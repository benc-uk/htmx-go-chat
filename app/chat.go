// ================================================================================
// Customised SSE broker to handle our chat messages
// ================================================================================

package main

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/benc-uk/go-rest-api/pkg/sse"
)

const serverUsername = "💻 Server Message"
const maxMsgsReloaded = 50

// ChatMessage is the data structure used for chats & system messages
type ChatMessage struct {
	Username string // Username of the sender
	Message  string // Message body
	System   bool   // Is this a special system message?
}

func initChat(msgStore *[]ChatMessage, renderer HTMLRenderer) sse.Broker[ChatMessage] {
	// The chat broker, it will handle chat messages and SSE events
	broker := sse.NewBroker[ChatMessage]()

	// Handle users joining the chat
	broker.ClientConnectedHandler = func(clientID string) {
		broker.Broadcast <- ChatMessage{
			Username: serverUsername,
			Message:  fmt.Sprintf("User '%s' has joined the chat 💬", clientID),
			System:   false,
		}

		broker.Broadcast <- ChatMessage{
			Username: "",
			Message:  fmt.Sprintf("There are %d users online", broker.GetClientCount()),
			System:   true,
		}

		// Send some stored messages to the new client
		maxMsg := len(*msgStore) - maxMsgsReloaded
		if maxMsg < 0 {
			maxMsg = 0
		}

		for _, msg := range (*msgStore)[maxMsg:] {
			broker.SendToClient(clientID, msg)
		}
	}

	// Handle users leaving the chat
	broker.ClientDisconnectedHandler = func(clientID string) {
		broker.Broadcast <- ChatMessage{
			Username: serverUsername,
			Message:  fmt.Sprintf("User '%s' has left the chat 👋", clientID),
			System:   false,
		}

		broker.Broadcast <- ChatMessage{
			Username: "",
			Message:  fmt.Sprintf("There are %d users online", broker.GetClientCount()),
			System:   true,
		}
	}

	// Handle chat & system messages and format them for SSE
	broker.MessageAdapter = func(msg ChatMessage, clientID string) sse.SSE {
		sse := sse.SSE{
			Event: "chat",
			Data:  "",
		}

		// Render the message using HTML template
		msgHTML, _ := renderer.RenderToString("message", map[string]any{
			"username": msg.Username,
			"message":  template.HTML(msg.Message), // nolint:gosec
			"time":     time.Now().Format("15:04:05"),
			"isSelf":   clientID == msg.Username,
			"isServer": msg.System || msg.Username == serverUsername,
		})

		// Write the HTML response, but we need to strip out newlines for SSE
		sse.Data = strings.Replace(msgHTML, "\n", "", -1)

		if msg.System {
			sse.Event = "system"
			sse.Data = msg.Message
		}

		return sse
	}

	return *broker
}