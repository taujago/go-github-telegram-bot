package parser

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Commit struct {
	ID       string   `json:"id"`
	Message  string   `json:"message"`
	Added    []string `json:"added"`
	Removed  []string `json:"removed"`
	Modified []string `json:"modified"`
	URL      string   `json:"url"`
}

type PushPayload struct {
	Ref        string `json:"ref"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
	Pusher struct {
		Name string `json:"name"`
	} `json:"pusher"`
	Commits []Commit `json:"commits"`
}

func ParsePush(body []byte) (string, error) {
	var payload PushPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}

	repo := payload.Repository.FullName
	branch := strings.TrimPrefix(payload.Ref, "refs/heads/")
	pusher := payload.Pusher.Name

	if len(payload.Commits) == 0 {
		return fmt.Sprintf("📦 **%s** pushed to `%s@%s` with no commits.", pusher, repo, branch), nil
	}

	commit := payload.Commits[0]
	action := "updated"
	files := append(append(commit.Added, commit.Modified...), commit.Removed...)

	switch {
	case len(commit.Added) > 0:
		action = "added"
		files = commit.Added
	case len(commit.Removed) > 0:
		action = "removed"
		files = commit.Removed
	}

	fileList := "`" + strings.Join(files, "`, `") + "`"
	return fmt.Sprintf(
		"📦 **%s** %s %s in `%s@%s`\n\n📝 Commit Message: *%s*\n\n🔗 %s",
		pusher,
		action,
		fileList,
		repo,
		branch,
		commit.Message,
		commit.URL,
	), nil
}
