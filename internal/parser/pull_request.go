package parser

import (
	"encoding/json"
	"fmt"
)

type PullRequestPayload struct {
	Action      string `json:"action"`
	PullRequest struct {
		Title   string `json:"title"`
		HTMLURL string `json:"html_url"`
		Merged  bool   `json:"merged"`
		User    struct {
			Login string `json:"login"`
		} `json:"user"`
	} `json:"pull_request"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
	Sender struct {
		Login string `json:"login"`
	} `json:"sender"`
}

func ParsePullRequest(body []byte) (string, error) {
	var payload PullRequestPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}

	repo := payload.Repository.FullName
	title := payload.PullRequest.Title
	url := payload.PullRequest.HTMLURL
	sender := payload.Sender.Login

	switch payload.Action {
	case "opened":
		return fmt.Sprintf("📣 **%s** opened a pull request in `%s`\n\n📄 *%s*\n🔗 %s", sender, repo, title, url), nil
	case "closed":
		if payload.PullRequest.Merged {
			return fmt.Sprintf("✅ **%s** merged a pull request in `%s`\n\n🔀 *%s*\n🔗 %s", sender, repo, title, url), nil
		}
		return fmt.Sprintf("🚫 **%s** closed a pull request in `%s` without merging\n\n📄 *%s*\n🔗 %s", sender, repo, title, url), nil
	case "reopened":
		return fmt.Sprintf("🔁 **%s** reopened a pull request in `%s`\n\n📄 *%s*\n🔗 %s", sender, repo, title, url), nil
	default:
		return "", nil // ignore other actions
	}
}
