// ================================================================================
// Broker handles SSE connections and events, it is the core of the chat server
// ================================================================================

package main

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

const serverUsername = "💻 Server Message"
const maxStoredMessages = 1000
const maxMsgsReloaded = 50

// ChatMessage is the data structure used for chats & system messages
type ChatMessage struct {
	Username string // Username of the sender
	Message  string // Message body
	System   bool   // Is this a special system message?
	Store    bool   // Should this message be stored in the message store?
}

// Struct to hold the chat broker state
type ChatBroker struct {
	// Push messages here to broadcast them to all connected clients
	Broadcast chan ChatMessage

	// New client connections, channel holds the username
	newClients chan string

	// Closed client connections, channel holds the username
	closingClients chan string

	// Main connections registry, keyed on username
	// Each client has their own message channel
	clients map[string]chan ChatMessage

	// Simple in memory message store, could be replaced with a database
	msgStore []ChatMessage
}

// Dead simple struct to support SSE format
type SSE struct {
	Event string
	Data  string
}

// Write the SSE format message to a writer
func (sse *SSE) Write(w io.Writer) {
	fmt.Fprintf(w, "event: %s\n", sse.Event)
	fmt.Fprintf(w, "data: %s\n\n", sse.Data)
}

// Create a new chat broker
func NewChatBroker() (broker *ChatBroker) {
	broker = &ChatBroker{
		// Buffered channel so we don't block
		Broadcast:      make(chan ChatMessage, 100),
		newClients:     make(chan string),
		closingClients: make(chan string),
		clients:        make(map[string]chan ChatMessage),
	}

	// Set it running, listening and broadcasting events
	// Note: This runs in a goroutine so we don't block here
	go broker.listen()

	return
}

// HTTP handler for connecting clients to the chat stream and sending SSE events
func (broker *ChatBroker) handleStream(username string, c echo.Context, plain bool) error {
	w := c.Response().Writer
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Each connection registers its own message channel with the broker's connections registry
	messageChan := make(chan ChatMessage)
	broker.clients[username] = messageChan

	// Signal the broker that we have a new connection
	broker.newClients <- username

	// Remove this client from the map of connected clients, when this handler exits.
	defer func() {
		broker.closingClients <- username
	}()

	// Listen to connection closing and un-register client
	go func() {
		<-c.Request().Context().Done()
		broker.closingClients <- username
	}()

	// Main loop for sending messages to the client
	for {
		// Blocks here until there is a new message in this client's messageChan
		msg := <-messageChan

		sse := &SSE{
			Event: "chat",
			Data:  msg.Message,
		}

		// Render the message using HTML template
		msgHTML, _ := c.Echo().Renderer.(*HTMLRenderer).RenderToString("message", map[string]any{
			"username": msg.Username,
			"message":  msg.Message,
			"time":     time.Now().Format("15:04:05"),
			"isSelf":   username == msg.Username,
			"isServer": msg.System || msg.Username == serverUsername,
		})

		if plain {
			// Write a plain text response, really only used for debugging
			sse.Data = msg.Username + " says " + msg.Message
		} else {
			// Write the HTML response, but we need to strip out newlines for SSE
			sse.Data = strings.Replace(msgHTML, "\n", "", -1)
		}

		if msg.System {
			sse.Event = "system"
			sse.Data = msg.Message
		}

		// Write the SSE to the response writer
		sse.Write(w)

		// Flush the data immediately as we are streaming data
		c.Response().Flush()
	}
}

// Listen on different channels and act accordingly
func (broker *ChatBroker) listen() {
	for {
		select {
		// CASE: New client has connected
		case username := <-broker.newClients:
			log.Printf("User '%s' added: %d active clients", username, len(broker.clients))

			broker.Broadcast <- ChatMessage{
				Username: serverUsername,
				Message:  fmt.Sprintf("User '%s' has joined the chat 💬", username),
				System:   false,
			}

			broker.Broadcast <- ChatMessage{
				Username: "",
				Message:  fmt.Sprintf("There are %d users online", len(broker.clients)),
				System:   true,
			}

			// Send a bunch of existing stored messages to the new client so they get some history
			maxMsg := len(broker.msgStore) - maxMsgsReloaded
			if maxMsg < 0 {
				maxMsg = 0
			}

			for _, msg := range broker.msgStore[maxMsg:] {
				broker.clients[username] <- msg
			}

		// CASE: Client has detached and we want to stop sending them messages
		case username := <-broker.closingClients:
			delete(broker.clients, username)

			log.Printf("User '%s' disconnected: %d active clients", username, len(broker.clients))

			broker.Broadcast <- ChatMessage{
				Username: serverUsername,
				Message:  fmt.Sprintf("User '%s' has left the chat 👋", username),
				System:   false,
			}

			broker.Broadcast <- ChatMessage{
				Username: "",
				Message:  fmt.Sprintf("There are %d users online", len(broker.clients)),
				System:   true,
			}

		// CASE: Message incoming on the broadcast channel
		case message := <-broker.Broadcast:
			// Store the message in the message store
			if message.Store {
				broker.msgStore = append(broker.msgStore, message)
				if len(broker.msgStore) > maxStoredMessages {
					broker.msgStore = broker.msgStore[1:]
				}
			}

			// Loop through all connected clients and broadcast
			// the message to their individual message channel
			for username := range broker.clients {
				broker.clients[username] <- message
			}
		}
	}
}

// Check if a user exists in the broker
func (broker *ChatBroker) UserExists(username string) bool {
	_, ok := broker.clients[username]
	return ok
}

// Get all active users in the broker
func (broker *ChatBroker) GetUsers() []string {
	var users []string
	for username := range broker.clients {
		users = append(users, username)
	}

	return users
}
