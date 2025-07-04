package github

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func GetTitleAndURLFromNodeID(nodeID string) (string, string, error) {
	query := `
		query($id: ID!) {
			node(id: $id) {
				... on Issue {
					title
					url
				}
				... on PullRequest {
					title
					url
				}
			}
		}
	`

	vars := map[string]string{"id": nodeID}
	body, _ := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": vars,
	})

	req, _ := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("GITHUB_TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	respBody, _ := io.ReadAll(res.Body)
	var resp struct {
		Data struct {
			Node struct {
				Title string `json:"title"`
				URL   string `json:"url"`
			} `json:"node"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return "", "", err
	}

	return resp.Data.Node.Title, resp.Data.Node.URL, nil
}
