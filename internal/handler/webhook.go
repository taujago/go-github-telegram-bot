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
	eventType := r.Header.Get("X-GitHub-Event")

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

	var message string
	switch eventType {
	case "push":
		message, err = parser.ParsePush(body)
	case "pull_request":
		message, err = parser.ParsePullRequest(body)
	case "create":
		message, err = parser.ParseCreate(body)
	case "delete":
		message, err = parser.ParseDelete(body)
	case "workflow_run":
		message, err = parser.ParseWorkflowRun(body)
	case "project":

		message, err = parser.ParseProject(body)
	case "project_card":
		message, err = parser.ParseProjectCard(body)
	case "project_column":
		message, err = parser.ParseProjectColumn(body)
	case "projects_v2_item":
		message, err = parser.ParseProjectsV2Item(body)

	default:
		fmt.Fprint(w, "Ignored event type: "+eventType)
		return
	}

	if err != nil {
		http.Error(w, "Error parsing payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := telegram.SendMessage(message); err != nil {
		http.Error(w, "Failed to send to Telegram: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "OK")
}
