// ================================================================================
// Application entry point and HTTP server setup & start
// ================================================================================

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
	e.Renderer = NewHTMLRenderer("templates")

	// We need server side sessions to store the state of the user
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("very_secret_12345"))))

	// Add all the routes
	addRoutes(e)

	// Start the server
	log.Println("Starting chat server on port: " + port)
	log.Println("Version: " + os.Getenv("VERSION"))

	e.Logger.Fatal(e.Start(":" + port))
}
