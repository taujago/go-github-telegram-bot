package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/taujago/go-github-telegram-bot/internal/parser"
	"github.com/taujago/go-github-telegram-bot/internal/telegram"
)

func isDebugEnabled() bool {
	return os.Getenv("DEBUG") == "true"
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-GitHub-Event") != "push" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Ignoring non-push event")
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}

	if isDebugEnabled() {
		log.Println("====== DEBUG: RAW GitHub Payload ======")
		log.Println(string(body))
		log.Println("========================================")
	}

	message, err := parser.ParsePush(body)
	if err != nil {
		http.Error(w, "Failed to parse payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := telegram.SendMessage(message); err != nil {
		http.Error(w, "Failed to send to Telegram: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "OK")
}
