package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	e := echo.New()

	// Configure with HTML template renderer and session middleware
	e.HideBanner = true
	e.Renderer = HTMLRenderer()
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("very_secret_12345"))))

	broker := NewChatBroker()

	// Root route renders the main page
	e.GET("/", func(c echo.Context) error {
		sess, _ := session.Get("session", c)

		return c.Render(http.StatusOK, "index", map[string]any{
			// Username might be empty or nil, the template will handle it
			"username": sess.Values["username"],
		})

	})

	// Login POST will set the username in the session and render the chat
	e.POST("/login", func(c echo.Context) error {
		username := c.FormValue("username")
		if username == "" {
			return c.Render(http.StatusOK, "login", map[string]any{
				"error": "Username can not be empty.",
			})
		}

		// Check user exists
		if _, ok := broker.Usernames[username]; ok {
			return c.Render(http.StatusOK, "login", map[string]any{
				"error": "That name is already taken, please pick another name.",
			})
		}

		sess, _ := session.Get("session", c)
		sess.Values["username"] = username
		err := sess.Save(c.Request(), c.Response())
		if err != nil {
			log.Println("Session error: ", err)
			return c.Redirect(http.StatusFound, "/")
		}

		// Render the chat template
		return c.Render(http.StatusOK, "chat", map[string]any{
			"username": username,
		})
	})

	// Connect clients to the chat stream
	e.GET("/connect_chat", func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		return broker.handleStream(c, sess.Values["username"].(string))
	})

	// Chat route for sending messages
	e.POST("/chat", func(c echo.Context) error {
		msgText := c.FormValue("message")
		username := c.FormValue("username")

		// Push the new chat message to broker
		broker.ChatMessages <- ChatMessage{
			Username: username,
			Message:  msgText,
		}

		return c.HTML(http.StatusOK, "")
	})

	// Used to logout
	e.POST("/logout", func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		delete(broker.Usernames, sess.Values["username"].(string))
		sess.Values["username"] = ""
		err := sess.Save(c.Request(), c.Response())
		if err != nil {
			log.Println("Session error: ", err)
		}

		return c.Render(http.StatusOK, "login", nil)
	})

	// Start the server
	log.Println("Starting server on port: " + port)
	e.Logger.Fatal(e.Start(":" + port))
}
