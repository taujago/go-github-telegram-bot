package parser

import (
	"encoding/json"
	"fmt"
)

type ProjectCardPayload struct {
	Action string `json:"action"`
	Card   struct {
		Note string `json:"note"`
		URL  string `json:"url"`
	} `json:"project_card"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
	Sender struct {
		Login string `json:"login"`
	} `json:"sender"`
}

func ParseProjectCard(body []byte) (string, error) {
	var payload ProjectCardPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}

	note := payload.Card.Note
	if note == "" {
		note = "_linked issue or PR_"
	}

	return fmt.Sprintf("🗂️ **%s** %s a card in project `%s`\n📝 %s",
		payload.Sender.Login,
		payload.Action,
		payload.Repository.FullName,
		note,
	), nil
}
