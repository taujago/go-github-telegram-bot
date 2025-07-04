package parser

import (
	"encoding/json"
	"fmt"
)

type CreatePayload struct {
	Ref        string `json:"ref"`
	RefType    string `json:"ref_type"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
	Sender struct {
		Login string `json:"login"`
	} `json:"sender"`
}

func ParseCreate(body []byte) (string, error) {
	var payload CreatePayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}

	if payload.RefType != "branch" {
		return "", nil
	}

	return fmt.Sprintf("🌿 **%s** created new branch `%s` in `%s`",
		payload.Sender.Login, payload.Ref, payload.Repository.FullName), nil
}
