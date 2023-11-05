package main

import (
	"log"
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
	e.Renderer = NewHTMLRenderer()
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("very_secret_12345"))))

	// This is our chat broker, it will handle our clients and SSE messages
	broker := NewChatBroker()

	addRoutes(e, broker)

	// Start the server
	log.Println("Starting chat server on port: " + port)
	log.Println("Open http://localhost:" + port + " in the browser to access the chat")
	e.Logger.Fatal(e.Start(":" + port))
}
