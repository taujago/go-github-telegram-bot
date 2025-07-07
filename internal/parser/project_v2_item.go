package parser

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/taujago/go-github-telegram-bot/internal/github"
)

// ProjectV2ItemPayload models the incoming webhook for a project_v2_item event
type ProjectV2ItemPayload struct {
	Action         string `json:"action"`
	ProjectsV2Item struct {
		ID            int64   `json:"id"`
		NodeID        string  `json:"node_id"`
		ProjectNodeID string  `json:"project_node_id"`
		ContentNodeID string  `json:"content_node_id"`
		ContentType   string  `json:"content_type"`
		CreatedAt     string  `json:"created_at"`
		UpdatedAt     string  `json:"updated_at"`
		ArchivedAt    *string `json:"archived_at"`
	} `json:"projects_v2_item"`
	Sender struct {
		Login string `json:"login"`
	} `json:"sender"`
	Organization struct {
		Login string `json:"login"`
	} `json:"organization"`
}

// ParseProjectV2Item parses the project_v2_item event and enriches it with the issue/PR title via GraphQL
func ParseProjectV2Item(body []byte) (string, error) {
	var payload ProjectV2ItemPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}

	// Default values
	title := "_unknown item_"
	url := ""

	if payload.ProjectsV2Item.ContentNodeID != "" {
		t, u, err := github.FetchIssueOrPRTitleByNodeID(payload.ProjectsV2Item.ContentNodeID)
		if err != nil {
			log.Printf("❌ Failed to fetch GitHub title from content_node_id: %v", err)
		} else {
			title = t
			url = u
		}
	}

	return fmt.Sprintf(
		"🧩 **%s** %s a project card in org *%s*\n📌 *%s*\n🔗 %s",
		payload.Sender.Login,
		payload.Action,
		payload.Organization.Login,
		title,
		url,
	), nil
}
