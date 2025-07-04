package parser

import (
	"encoding/json"
	"fmt"

	"github.com/taujago/go-github-telegram-bot/internal/github"
)

func ParseProjectsV2Item(body []byte) (string, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return "", err
	}

	action := raw["action"].(string)
	sender := raw["sender"].(map[string]interface{})["login"].(string)
	contentNodeID := raw["projects_v2_item"].(map[string]interface{})["content_node_id"].(string)

	title, url, err := github.GetTitleAndURLFromNodeID(contentNodeID)
	if err != nil {
		title = "_(unknown title)_"
		url = ""
	}

	if action == "edited" {
		if changes, ok := raw["changes"].(map[string]interface{}); ok {
			if fieldVal, ok := changes["field_value"].(map[string]interface{}); ok {
				if fieldVal["field_name"] == "Status" {
					from := fieldVal["from"].(map[string]interface{})["name"].(string)
					to := fieldVal["to"].(map[string]interface{})["name"].(string)
					return fmt.Sprintf("\U0001F501 **%s** moved task: [%s](%s)\n\u27a1\ufe0f **%s → %s**", sender, title, url, from, to), nil
				}
			}
		}
	}

	return fmt.Sprintf("\U0001F501 **%s** %s a task: [%s](%s)", sender, action, title, url), nil
}
