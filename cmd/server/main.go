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

	botToken := os.Getenv("BOT_TOKEN")
	chatID := os.Getenv("CHAT_ID")

	log.Println("📨 Telegram Config:")
	log.Println("  BOT_TOKEN:", maskToken(botToken))
	log.Println("  CHAT_ID  :", chatID)

	http.HandleFunc("/webhook", handler.WebhookHandler)

	log.Printf("🚀 Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func maskToken(token string) string {
	if len(token) <= 10 {
		return "********"
	}
	return token[:10] + "..."
}
