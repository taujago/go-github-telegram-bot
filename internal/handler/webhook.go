// internal/handler/webhook.go
package handler

import (
	"encoding/json"
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
	case "issues":
		message, err = parser.ParseIssue(body)
	case "issue_comment":
		message, err = parser.ParseIssueComment(body)
	case "project_v2_item":
		if isDebugEnabled() {
			log.Println("Handling project_v2_item...")
		}

		var payload parser.ProjectV2ItemEditedPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			http.Error(w, "Failed to parse project_v2_item", http.StatusBadRequest)
			return
		}

		switch payload.Action {
		case "edited":
			if payload.Changes != nil && payload.Changes.FieldValue != nil {
				if isDebugEnabled() {
					log.Println("Detected FieldValue change:", payload.Changes.FieldValue)
				}

				info, err := parser.ParseProjectV2ItemEdited(payload)
				if err != nil {
					http.Error(w, "Failed to enrich project_v2_item: "+err.Error(), http.StatusInternalServerError)
					return
				}

				message = fmt.Sprintf(
					"\U0001F4E6 *%s*\n\U0001F3C1 Project #%s\n\U0001F464 Moved by: _%s_\n\U0001F501 `%s` → `%s`\n🔗 [View Card](%s)",
					info.CardTitle,
					info.ProjectNumber,
					payload.Sender.Login,
					info.ColumnFrom,
					info.ColumnTo,
					info.CardURL,
				)
			} else {
				if isDebugEnabled() {
					log.Println("Ignored project_v2_item: not a field_value change")
				}
				fmt.Fprint(w, "Ignored project_v2_item action or no column change")
				return
			}

		case "reordered":
			if isDebugEnabled() {
				log.Println("Detected reordered action on project_v2_item")
			}

			info, err := parser.ParseProjectV2ItemReordered(payload)
			if err != nil {
				http.Error(w, "Failed to parse reordered payload: "+err.Error(), http.StatusInternalServerError)
				return
			}

			message = fmt.Sprintf(
				"\U0001F500 *%s*\n\U0001F3C1 Project #%s\n\U0001F464 Reordered by: _%s_\n\U0001F4CB Column: `%s`\n🔗 [View Card](%s)",
				info.CardTitle,
				info.ProjectNumber,
				payload.Sender.Login,
				info.ColumnName,
				info.CardURL,
			)

		default:
			if isDebugEnabled() {
				log.Println("Unhandled project_v2_item action:", payload.Action)
			}
			fmt.Fprint(w, "Unhandled project_v2_item action")
			return
		}

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
