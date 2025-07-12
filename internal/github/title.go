package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchIssueOrPRTitleByNodeID(nodeID string, token string) (string, string, error) {
	query := fmt.Sprintf(`{
		node(id: "%s") {
			__typename
			... on Issue {
				title
				url
			}
			... on PullRequest {
				title
				url
			}
		}
	}`, nodeID)

	payload := map[string]string{"query": query}
	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Data struct {
			Node struct {
				Typename string `json:"__typename"`
				Title    string `json:"title"`
				URL      string `json:"url"`
			} `json:"node"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", err
	}

	return result.Data.Node.Title, result.Data.Node.URL, nil
}
