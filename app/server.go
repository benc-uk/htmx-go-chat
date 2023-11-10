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

	echo := echo.New()

	// Configure with HTML template renderer and session middleware
	echo.HideBanner = true
	echo.Renderer = NewHTMLRenderer("templates")
	echo.Use(session.Middleware(sessions.NewCookieStore([]byte("very_secret_12345"))))

	// Add all the routes
	addRoutes(echo)

	// Start the server
	log.Println("🚀 Starting HTMX chat server on port: " + port)
	log.Println("✅ Version: " + os.Getenv("VERSION"))

	echo.Logger.Fatal(echo.Start(":" + port))
}
