// ================================================================================
// All HTTP routes are defined here, purely for code organisation purposes.
// ================================================================================

package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/benc-uk/go-rest-api/pkg/sse"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const maxStoredMessages = 1000

func addRoutes(e *echo.Echo, broker sse.Broker[ChatMessage], msgStore *[]ChatMessage) {
	//
	// Root route renders the main index.html template
	//
	e.GET("/", func(c echo.Context) error {
		sess, _ := session.Get("session", c)

		return c.Render(http.StatusOK, "index", map[string]any{
			// Username might be empty or nil, the template will handle it
			"username": sess.Values["username"],
		})
	})

	//
	// Login POST will set the username in the session and render the chat view
	//
	e.POST("/login", func(c echo.Context) error {
		username := c.FormValue("username")
		if username == "" {
			return c.Render(http.StatusOK, "login", map[string]any{
				"error": "Username can not be empty.",
			})
		}

		// Check if name exists
		activeUsers := broker.GetClients()
		for _, user := range activeUsers {
			if user == username {
				return c.Render(http.StatusOK, "login", map[string]any{
					"error": "That name is already taken, please pick another name.",
				})
			}
		}

		sess, _ := session.Get("session", c)
		sess.Values["username"] = username
		err := sess.Save(c.Request(), c.Response())
		if err != nil {
			log.Println("Session error: ", err)
			return c.Render(http.StatusOK, "login", map[string]any{
				"error": err.Error(),
			})
		}

		// Got this far, we can render the chat template
		return c.Render(http.StatusOK, "chat", map[string]any{
			"username":       username,
			"addLoginButton": true,
		})
	})

	//
	// Connect clients to the chat stream using the broker
	//
	e.GET("/chat-stream", func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		username := sess.Values["username"].(string)

		return broker.Stream(username, c.Response().Writer, *c.Request())
	})

	//
	// Post messages to the chat for broadcast
	//
	e.POST("/chat", func(c echo.Context) error {
		msgText := c.FormValue("message")
		username := c.FormValue("username")

		// Trim the message
		msgText = strings.TrimSpace(msgText)

		if msgText == "" {
			return c.HTML(http.StatusBadRequest, "")
		}

		msg := ChatMessage{
			Username: username,
			Message:  msgText,
		}

		// Push the new chat message to broker
		broker.Broadcast <- msg

		// Add the message to the message store
		*msgStore = append(*msgStore, msg)

		// Trim the message store if it gets too big
		if len(*msgStore) > maxStoredMessages {
			*msgStore = (*msgStore)[1:]
		}

		return c.HTML(http.StatusOK, "")
	})

	//
	// Used to logout
	//
	e.POST("/logout", func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		sess.Values["username"] = ""

		err := sess.Save(c.Request(), c.Response())
		if err != nil {
			log.Println("Session error: ", err)
		}

		return c.Render(http.StatusOK, "login", nil)
	})

	//
	// Display the 'about' modal popup
	//
	e.GET("/modal-about", func(c echo.Context) error {
		ver := os.Getenv("VERSION")
		if ver == "" {
			ver = "Unknown!"
		}

		return c.Render(http.StatusOK, "modal-about", map[string]any{
			"version": ver,
		})
	})

	//
	// Display the users list in a modal popup
	//
	e.GET("/modal-users", func(c echo.Context) error {
		users := broker.GetClients()

		return c.Render(http.StatusOK, "modal-users", map[string]any{
			"users": users,
		})
	})

	//
	// Display the users list
	//
	e.GET("/users", func(c echo.Context) error {
		users := broker.GetClients()

		// Return users as a basic HTML list
		var html string
		for _, user := range users {
			html += "<li>" + user + "</li>"
		}

		return c.HTML(http.StatusOK, html)
	})
}
