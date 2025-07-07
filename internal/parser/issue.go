package parser

import (
	"encoding/json"
	"fmt"
)

type issuePayload struct {
	Action string `json:"action"`
	Issue  struct {
		HTMLURL string `json:"html_url"`
		Title   string `json:"title"`
		State   string `json:"state"`
		User    struct {
			Login string `json:"login"`
		} `json:"user"`
	} `json:"issue"`
	Comment *struct {
		Body    string `json:"body"`
		HTMLURL string `json:"html_url"`
		User    struct {
			Login string `json:"login"`
		} `json:"user"`
	} `json:"comment,omitempty"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

func ParseIssue(body []byte) (string, error) {
	var payload issuePayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}

	var msg string

	switch payload.Action {
	case "opened":
		msg = fmt.Sprintf(
			"📌 *New Issue Opened*\n📁 Repo: `%s`\n👤 Author: `%s`\n📝 Title: *%s*\n🔗 %s",
			payload.Repository.FullName,
			payload.Issue.User.Login,
			payload.Issue.Title,
			payload.Issue.HTMLURL,
		)
	case "closed":
		msg = fmt.Sprintf(
			"✅ *Issue Closed*\n📁 Repo: `%s`\n📝 Title: *%s*\n🔗 %s",
			payload.Repository.FullName,
			payload.Issue.Title,
			payload.Issue.HTMLURL,
		)
	case "edited":
		msg = fmt.Sprintf(
			"✏️ *Issue Updated*\n📁 Repo: `%s`\n📝 Title: *%s*\n🔗 %s",
			payload.Repository.FullName,
			payload.Issue.Title,
			payload.Issue.HTMLURL,
		)
	case "reopened":
		msg = fmt.Sprintf(
			"🔄 *Issue Reopened*\n📁 Repo: `%s`\n📝 Title: *%s*\n🔗 %s",
			payload.Repository.FullName,
			payload.Issue.Title,
			payload.Issue.HTMLURL,
		)
	case "created":
		if payload.Comment != nil {
			msg = fmt.Sprintf(
				"💬 *New Comment on Issue*\n📁 Repo: `%s`\n👤 Comment by: `%s`\n📝 Comment: _%s_\n🔗 %s",
				payload.Repository.FullName,
				payload.Comment.User.Login,
				payload.Comment.Body,
				payload.Comment.HTMLURL,
			)
		}
	default:
		msg = fmt.Sprintf("⚠️ Unhandled issue action: %s", payload.Action)
	}

	return msg, nil
}
