package parser

import (
	"encoding/json"
	"fmt"
)

type ProjectsV2ItemPayload struct {
	Action         string `json:"action"`
	ProjectsV2Item struct {
		ContentType string  `json:"content_type"`
		CreatedAt   string  `json:"created_at"`
		UpdatedAt   string  `json:"updated_at"`
		ArchivedAt  *string `json:"archived_at"`
		Creator     struct {
			Login string `json:"login"`
		} `json:"creator"`
	} `json:"projects_v2_item"`
	Sender struct {
		Login string `json:"login"`
	} `json:"sender"`
}

func ParseProjectsV2Item(body []byte) (string, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return "", err
	}

	action := raw["action"].(string)
	sender := raw["sender"].(map[string]interface{})["login"].(string)

	// Check if this is a status move
	if action == "edited" {
		if changes, ok := raw["changes"].(map[string]interface{}); ok {
			if fieldVal, ok := changes["field_value"].(map[string]interface{}); ok {
				if fieldVal["field_name"] == "Status" {
					from := fieldVal["from"].(map[string]interface{})["name"].(string)
					to := fieldVal["to"].(map[string]interface{})["name"].(string)
					return fmt.Sprintf("📌 **%s** moved task from **%s** to **%s**", sender, from, to), nil
				}
			}
		}
	}

	// Fallback for created or edited without status change
	payload := struct {
		ProjectsV2Item struct {
			ContentType string `json:"content_type"`
			Creator     struct {
				Login string `json:"login"`
			} `json:"creator"`
		} `json:"projects_v2_item"`
	}{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}

	return fmt.Sprintf("🗂️ **%s** %s a `%s` in GitHub Projects v2 (created by **%s**)", sender, action, payload.ProjectsV2Item.ContentType, payload.ProjectsV2Item.Creator.Login), nil
}
