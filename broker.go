package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type ChatMessage struct {
	Username string
	Message  string
	Server   bool
}

type ChatBroker struct {
	// Push messages here to broadcast them.
	ChatMessages chan ChatMessage

	// New client connections, channel of channels!
	newClients chan chan ChatMessage

	// Closed client connections, channel of channels!
	closingClients chan chan ChatMessage

	// Client connections registry
	clients map[chan ChatMessage]bool

	// List of Usernames
	Usernames map[string]bool
}

func NewChatBroker() (broker *ChatBroker) {
	broker = &ChatBroker{
		ChatMessages:   make(chan ChatMessage, 1),
		newClients:     make(chan chan ChatMessage),
		closingClients: make(chan chan ChatMessage),
		clients:        make(map[chan ChatMessage]bool),
		Usernames:      make(map[string]bool),
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

	log.Printf("Client connected: %s", username)

	broker.Usernames[username] = true

	// Each connection registers its own message channel with the broker's connections registry
	messageChan := make(chan ChatMessage)

	// Signal the broker that we have a new connection
	broker.newClients <- messageChan

	// Remove this client from the map of connected clients, when this handler exits.
	defer func() {
		broker.closingClients <- messageChan
	}()

	// Listen to connection close and un-register client
	go func() {
		<-c.Request().Context().Done()
		broker.closingClients <- messageChan
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
		// New client has connected, register their message channel
		case s := <-broker.newClients:
			broker.clients[s] = true

			log.Printf("Client added: %d active clients", len(broker.clients))
			broker.ChatMessages <- ChatMessage{
				Username: "Server",
				Message:  fmt.Sprintf("There are %d users online", len(broker.clients)),
				Server:   true,
			}

		// Client has detached and we want to stop sending them messages
		case s := <-broker.closingClients:
			delete(broker.clients, s)

			log.Printf("Removed client: %d active clients", len(broker.clients))
			broker.ChatMessages <- ChatMessage{
				Username: "Server",
				Message:  fmt.Sprintf("There are %d users online", len(broker.clients)),
				Server:   true,
			}

		// We got a new message from the outside
		case message := <-broker.ChatMessages:
			// Loop through all connected clients and broadcast the message
			for clientMessageChan := range broker.clients {
				clientMessageChan <- message
			}
		}
	}
}
