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

	cookieKey := os.Getenv("COOKIE_KEY")
	if cookieKey == "" {
		cookieKey = "cookie-secret-1234567890"
	}

	echo := echo.New()

	// Configure with HTML template renderer and session middleware
	echo.HideBanner = true
	htmlRenderer := NewHTMLRenderer("templates")
	echo.Renderer = htmlRenderer
	echo.Use(session.Middleware(sessions.NewCookieStore([]byte(cookieKey))))

	// Open and/or create the database
	db := openDB()
	defer db.Close()

	// Initialise the chat broker and add routes
	broker := initChat(db, *htmlRenderer)
	addRoutes(echo, broker, db)

	// Start the server
	log.Println("🚀 Starting HTMX chat server on port: " + port)
	log.Println("✅ Version: " + os.Getenv("VERSION"))

	echo.Logger.Fatal(echo.Start(":" + port))
}
