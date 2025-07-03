package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/taujago/go-github-telegram-bot/internal/handler"
)

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/webhook", handler.WebhookHandler)

	log.Printf("🚀 Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
