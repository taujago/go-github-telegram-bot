package parser

import (
	"encoding/json"
	"fmt"
)

type ProjectColumnPayload struct {
	Action string `json:"action"`
	Column struct {
		Name string `json:"name"`
	} `json:"project_column"`
	Sender struct {
		Login string `json:"login"`
	} `json:"sender"`
}

func ParseProjectColumn(body []byte) (string, error) {
	var payload ProjectColumnPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}

	return fmt.Sprintf("📁 **%s** %s project column *%s*",
		payload.Sender.Login,
		payload.Action,
		payload.Column.Name,
	), nil
}
