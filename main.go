package main

import (
	"haikutie/data"
	"log"
	"net/http"

	"haikutie/handlers"
)

func main() {
	// Initialize database
	database := data.InitDB()
	defer database.Close()

	// Initialize handlers with database
	h := handlers.New(database)

	// Static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes
	http.HandleFunc("/", h.Home)
	http.HandleFunc("/login", h.Login)
	http.HandleFunc("/logout", h.Logout)
	http.HandleFunc("/register", h.Register)
	http.HandleFunc("/library", h.Library)
	http.HandleFunc("/compose", h.Compose)
	http.HandleFunc("/send", h.SendHaiku)
	http.HandleFunc("/haikus/received", h.ReceivedHaikus)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
