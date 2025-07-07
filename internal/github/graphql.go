package github

// **Add these imports at the top**
import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func DoGraphQL(query string, variables interface{}, target interface{}) error {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("GITHUB_TOKEN not set")
	}

	payload := struct {
		Query     string      `json:"query"`
		Variables interface{} `json:"variables"`
	}{Query: query, Variables: variables}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("GraphQL query failed: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(&target)
}

func FetchIssueOrPRTitleByNodeID(nodeID string) (string, string, error) {
	const q = `
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
    }`
	variables := map[string]string{"id": nodeID}

	var resp struct {
		Data struct {
			Node struct {
				Title string `json:"title"`
				URL   string `json:"url"`
			} `json:"node"`
		} `json:"data"`
	}

	if err := DoGraphQL(q, variables, &resp); err != nil {
		return "", "", err
	}
	return resp.Data.Node.Title, resp.Data.Node.URL, nil
}
