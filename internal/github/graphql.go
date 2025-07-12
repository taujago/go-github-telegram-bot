package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const githubGraphQLEndpoint = "https://api.github.com/graphql"

type GraphQLResponse struct {
	Data struct {
		Node struct {
			FieldValues struct {
				Nodes []struct {
					ProjectField struct {
						Name string `json:"name"`
					} `json:"projectField"`
					Value string `json:"value"`
				} `json:"nodes"`
			} `json:"fieldValues"`
		} `json:"node"`
	} `json:"data"`
}

func GetCurrentColumn(projectItemID, token string) (string, error) {
	query := fmt.Sprintf(`{
		node(id: "%s") {
			... on ProjectV2Item {
				fieldValues(first: 10) {
					nodes {
						projectField {
							name
						}
						value
					}
				}
			}
		}
	}`, projectItemID)

	payload := map[string]string{"query": query}
	data, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", githubGraphQLEndpoint, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result GraphQLResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	for _, node := range result.Data.Node.FieldValues.Nodes {
		if node.ProjectField.Name == "Status" {
			return node.Value, nil
		}
	}

	return "(unknown)", nil
}
