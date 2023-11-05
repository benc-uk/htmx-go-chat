// ================================================================================
// Broker handles SSE connections and events, it is the core of the chat server
// ================================================================================

package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// ChatMessage is the data structure we send via the broker to connected clients
type ChatMessage struct {
	Username string // Username of the sender
	Message  string // Message body
	Server   bool   // Is this a special server message?
}

type ChatBroker struct {
	// Push messages here to broadcast them.
	Broadcast chan ChatMessage

	// New client connections, channel holds the username
	newClients chan string

	// Closed client connections, channel holds the username
	closingClients chan string

	// Client connections registry, key is the username
	clients map[string]chan ChatMessage
}

func NewChatBroker() (broker *ChatBroker) {
	broker = &ChatBroker{
		Broadcast:      make(chan ChatMessage, 10),
		newClients:     make(chan string),
		closingClients: make(chan string),
		clients:        make(map[string]chan ChatMessage),
	}

	// Set it running, listening and broadcasting events
	go broker.listen()

	return
}

// HTTP handler for connecting clients to the chat stream
func (broker *ChatBroker) handleStream(c echo.Context, username string) error {
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

	// Listen to connection close and un-register client
	go func() {
		<-c.Request().Context().Done()
		broker.closingClients <- username
	}()

	// Main loop for sending messages to the client
	for {
		msg := <-messageChan

		sess, _ := session.Get("session", c)

		// Render the message using HTML template
		msgData, _ := c.Echo().Renderer.(*HTMLRenderer).RenderToString("message", map[string]any{
			"username": msg.Username,
			"message":  msg.Message,
			"time":     time.Now().Format("15:04:05"),
			"isSelf":   sess.Values["username"] == msg.Username,
		})

		// Remove all newlines
		msgData = strings.Replace(msgData, "\n", "", -1)

		// Set the message type
		msgType := "chat"

		// If this is a server message just send it plain text
		if msg.Server {
			msgType = "server"
			msgData = msg.Message
		}

		// Write an SSE formatted response, yes the data is HTML!
		fmt.Fprintf(w, "event: %s\n", msgType)
		fmt.Fprintf(w, "data: %s\n\n", msgData)

		// Flush the data immediately as we are streaming data
		c.Response().Flush()
	}
}

// Listen on different channels and act accordingly
func (broker *ChatBroker) listen() {
	for {
		select {
		// New client has connected
		case username := <-broker.newClients:
			log.Printf("User '%s' added: %d active clients", username, len(broker.clients))

			broker.Broadcast <- ChatMessage{
				Username: "💻 Server Message",
				Message:  fmt.Sprintf("User '%s' has joined the chat 💬", username),
				Server:   false,
			}

			broker.Broadcast <- ChatMessage{
				Username: "",
				Message:  fmt.Sprintf("There are %d users online", len(broker.clients)),
				Server:   true,
			}

		// Client has detached and we want to stop sending them messages
		case username := <-broker.closingClients:
			delete(broker.clients, username)

			log.Printf("User '%s' disconnected: %d active clients", username, len(broker.clients))

			broker.Broadcast <- ChatMessage{
				Username: "💻 Server Message",
				Message:  fmt.Sprintf("User '%s' has left the chat 👋", username),
				Server:   false,
			}

			broker.Broadcast <- ChatMessage{
				Username: "Server",
				Message:  fmt.Sprintf("There are %d users online", len(broker.clients)),
				Server:   true,
			}

		// We got a new message from the outside
		case message := <-broker.Broadcast:
			// Loop through all connected clients and broadcast the message
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
