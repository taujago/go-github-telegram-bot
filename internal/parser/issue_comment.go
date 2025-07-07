package parser

import (
	"encoding/json"
	"fmt"
)

type issueCommentPayload struct {
	Action string `json:"action"`
	Issue  struct {
		HTMLURL string `json:"html_url"`
		Title   string `json:"title"`
		Number  int    `json:"number"`
		User    struct {
			Login string `json:"login"`
		} `json:"user"`
	} `json:"issue"`
	Comment struct {
		Body    string `json:"body"`
		HTMLURL string `json:"html_url"`
		User    struct {
			Login string `json:"login"`
		} `json:"user"`
	} `json:"comment"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

func ParseIssueComment(body []byte) (string, error) {
	var payload issueCommentPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}

	if payload.Action != "created" {
		return "", fmt.Errorf("ignored issue_comment action: %s", payload.Action)
	}

	msg := fmt.Sprintf(
		"💬 *New Comment on Issue*\n📁 Repo: `%s`\n🔖 Issue: *%s* (#%d)\n👤 Comment by: `%s`\n📝 _%s_\n🔗 %s",
		payload.Repository.FullName,
		payload.Issue.Title,
		payload.Issue.Number,
		payload.Comment.User.Login,
		payload.Comment.Body,
		payload.Comment.HTMLURL,
	)

	return msg, nil
}
