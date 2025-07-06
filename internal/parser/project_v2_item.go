package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/taujago/go-github-telegram-bot/internal/telegram"
)

type ProjectV2ItemPayload struct {
	Action         string `json:"action"`
	ProjectsV2Item struct {
		ID            int    `json:"id"`
		NodeID        string `json:"node_id"`
		ProjectNodeID string `json:"project_node_id"`
		ContentNodeID string `json:"content_node_id"`
		ContentType   string `json:"content_type"`
		CreatedAt     string `json:"created_at"`
		UpdatedAt     string `json:"updated_at"`
		ArchivedAt    string `json:"archived_at"`
		Creator       struct {
			Login string `json:"login"`
		} `json:"creator"`
		Title string `json:"title"`
	} `json:"projects_v2_item"`
	Changes struct {
		FieldValue struct {
			FieldName string `json:"field_name"`
			From      struct {
				Name string `json:"name"`
			} `json:"from"`
			To struct {
				Name string `json:"name"`
			} `json:"to"`
		} `json:"field_value"`
	} `json:"changes"`
	Organization struct {
		Login string `json:"login"`
	} `json:"organization"`
}

// fallback title cache (in-memory)
var draftTitleCache = make(map[string]string)

func ParseProjectsV2Item(body []byte) error {
	if os.Getenv("DEBUG") == "true" {
		fmt.Println("DEBUG: Received projects_v2_item event")
	}

	var payload ProjectV2ItemPayload
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return err
	}

	creator := payload.ProjectsV2Item.Creator.Login
	contentNodeID := payload.ProjectsV2Item.ContentNodeID
	title := ""

	// Use title directly from payload if available
	if payload.ProjectsV2Item.Title != "" {
		title = payload.ProjectsV2Item.Title
		draftTitleCache[contentNodeID] = title
	}

	// fallback if no title available (still show something)
	if title == "" && payload.ProjectsV2Item.ContentType == "DraftIssue" {
		title = draftTitleCache[contentNodeID]
		if title == "" {
			title = fmt.Sprintf("Draft Task %s", contentNodeID)
		}
	}

	// final fallback
	if title == "" {
		title = fmt.Sprintf("%s (%s)", contentNodeID, payload.ProjectsV2Item.ContentType)
	}

	from := payload.Changes.FieldValue.From.Name
	to := payload.Changes.FieldValue.To.Name

	org := payload.Organization.Login
	projectNumber := fetchProjectNumber(payload.ProjectsV2Item.ProjectNodeID)

	cardURL := fmt.Sprintf("https://github.com/orgs/%s/projects/%s/views/1?pane=issue&itemId=%d", org, projectNumber, payload.ProjectsV2Item.ID)
	message := fmt.Sprintf(
		"🔁 %s moved task: *%s*\n➡️ %s → %s\n🔗 [Open Card](%s)",
		creator,
		escapeMarkdown(title),
		escapeMarkdown(from),
		escapeMarkdown(to),
		cardURL,
	)

	return telegram.SendMessage(message)
}

func fetchProjectNumber(projectNodeID string) string {
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		return "1" // fallback
	}

	graphql := `query { node(id: "` + projectNodeID + `") { ... on ProjectV2 { number } } }`
	reqBody := map[string]string{"query": graphql}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "1"
	}
	req.Header.Set("Authorization", "Bearer "+githubToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "1"
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Node struct {
				Number int `json:"number"`
			} `json:"node"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "1"
	}

	return fmt.Sprintf("%d", result.Data.Node.Number)
}

func escapeMarkdown(text string) string {
	replacer := strings.NewReplacer("_", "\\_", "*", "\\*", "[", "\\[", "]", "\\]", "(`", "\\(`", ")", "\\)")
	return replacer.Replace(text)
}
