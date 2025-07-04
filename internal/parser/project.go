package parser

import (
	"encoding/json"
	"fmt"
)

type ProjectPayload struct {
	Action  string `json:"action"`
	Project struct {
		Name string `json:"name"`
		URL  string `json:"html_url"`
	} `json:"project"`
	Sender struct {
		Login string `json:"login"`
	} `json:"sender"`
}

func ParseProject(body []byte) (string, error) {
	var payload ProjectPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}

	return fmt.Sprintf("📋 **%s** %s project *%s*\n🔗 %s",
		payload.Sender.Login,
		payload.Action,
		payload.Project.Name,
		payload.Project.URL,
	), nil
}
