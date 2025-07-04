package parser

import (
	"encoding/json"
	"fmt"
)

type WorkflowPayload struct {
	Action      string `json:"action"`
	WorkflowRun struct {
		Name       string `json:"name"`
		Status     string `json:"status"`
		Conclusion string `json:"conclusion"`
		HTMLURL    string `json:"html_url"`
		HeadBranch string `json:"head_branch"`
	} `json:"workflow_run"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

func ParseWorkflowRun(body []byte) (string, error) {
	var payload WorkflowPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}

	if payload.Action != "completed" {
		return "", nil
	}

	statusIcon := "✅"
	if payload.WorkflowRun.Conclusion != "success" {
		statusIcon = "❌"
	}

	return fmt.Sprintf("%s GitHub Action *%s* on branch `%s` for `%s`\nStatus: *%s*\n🔗 %s",
		statusIcon,
		payload.WorkflowRun.Name,
		payload.WorkflowRun.HeadBranch,
		payload.Repository.FullName,
		payload.WorkflowRun.Conclusion,
		payload.WorkflowRun.HTMLURL,
	), nil
}
