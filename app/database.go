// ================================================================================
// Very simple database store using SQLite to hold a history of chat messages
// ================================================================================

package main

import (
	"database/sql"
	"log"
	"os"
	"path"

	_ "modernc.org/sqlite"
)

// Open the database and create the table if not exists
func openDB() *sql.DB {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "."
	}

	db, err := sql.Open("sqlite", path.Join(dbPath, "chat.db"))
	if err != nil {
		log.Fatal(err)
	}

	// Create table if not exists
	_, err = db.Exec(`
	  CREATE TABLE IF NOT EXISTS 
	  messages (id INTEGER PRIMARY KEY, username TEXT, message TEXT, timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)
	`)
	if err != nil {
		log.Printf("Error creating database, check path '%s' exists and is writeable", dbPath)
		log.Fatal(err)
	}

	log.Printf("🧮 Database '%s' opened successfully", path.Join(dbPath, "chat.db"))

	return db
}

// Store a message in the database
func storeMessage(db *sql.DB, msg ChatMessage) {
	// TODO: Add some kind of cleanup to remove old messages
	_, err := db.Exec("INSERT INTO messages (username, message) VALUES (?, ?)", msg.Username, msg.Message)
	if err != nil {
		log.Println("DB error: ", err)
	}
}

// Fetch the last n messages from the database
func fetchMessages(db *sql.DB, count int) []ChatMessage {
	rows, err := db.Query("SELECT username, message FROM messages ORDER BY timestamp ASC LIMIT ?", count)
	if err != nil {
		log.Println(err)
	}

	messages := []ChatMessage{}

	for rows.Next() {
		var username, message string

		err = rows.Scan(&username, &message)
		if err != nil {
			log.Println(err)
		}

		messages = append(messages, ChatMessage{
			Username: username,
			Message:  message,
		})
	}

	return messages
}
