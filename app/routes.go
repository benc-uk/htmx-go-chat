package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func addRoutes(e *echo.Echo, broker *ChatBroker) {
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
		if broker.UserExists(username) {
			return c.Render(http.StatusOK, "login", map[string]any{
				"error": "That name is already taken, please pick another name.",
			})
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
	e.GET("/connect_chat", func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		username := sess.Values["username"].(string)

		return broker.handleStream(username, c)
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

		// Push the new chat message to broker
		broker.Broadcast <- ChatMessage{
			Username: username,
			Message:  msgText,
			Store:    true,
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
	e.GET("/about", func(c echo.Context) error {
		ver := os.Getenv("VERSION")
		if ver == "" {
			ver = "Unknown!"
		}

		return c.Render(http.StatusOK, "about", map[string]any{
			"version": ver,
		})
	})

	//
	// Display the users list in a modal popup
	//
	e.GET("/users_modal", func(c echo.Context) error {
		users := broker.GetUsers()
		return c.Render(http.StatusOK, "users", map[string]any{
			"users": users,
		})
	})

	//
	// Display the users list
	//
	e.GET("/users", func(c echo.Context) error {
		users := broker.GetUsers()

		// Return users as a basic HTML list
		var html string
		for _, user := range users {
			html += "<li>" + user + "</li>"
		}

		return c.HTML(http.StatusOK, html)
	})
}
