// internal/parser/project_v2_item.go
package parser

import (
	"errors"
	"fmt"
)

type ProjectV2ItemEditedPayload struct {
	Action         string         `json:"action"`
	ProjectsV2Item ProjectsV2Item `json:"projects_v2_item"`
	Changes        *Changes       `json:"changes,omitempty"`
	Sender         User           `json:"sender"`
}

type ProjectsV2Item struct {
	ID            int64  `json:"id"`
	NodeID        string `json:"node_id"`
	ProjectNodeID string `json:"project_node_id"`
	ColumnName    string `json:"column_name,omitempty"` // hypothetical enrichment field
	CardTitle     string `json:"card_title,omitempty"`  // hypothetical enrichment field
	ProjectNumber string `json:"project_number,omitempty"`
	CardURL       string `json:"card_url,omitempty"`
}

type Changes struct {
	FieldValue *FieldValueChange `json:"field_value,omitempty"`
}

type FieldValueChange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type User struct {
	Login string `json:"login"`
}

type ParsedProjectV2ItemInfo struct {
	CardTitle     string
	ProjectNumber string
	ColumnFrom    string
	ColumnTo      string
	ColumnName    string
	CardURL       string
}

func ParseProjectV2ItemEdited(payload ProjectV2ItemEditedPayload) (*ParsedProjectV2ItemInfo, error) {
	if payload.Changes == nil || payload.Changes.FieldValue == nil {
		return nil, errors.New("missing field value change")
	}

	// Replace the following with actual GraphQL enrichment if needed
	return &ParsedProjectV2ItemInfo{
		CardTitle:     "Card Title Placeholder",
		ProjectNumber: "1",
		ColumnFrom:    payload.Changes.FieldValue.From,
		ColumnTo:      payload.Changes.FieldValue.To,
		CardURL:       fmt.Sprintf("https://github.com/orgs/PMJ-Project/projects/1/views/1?pane=issue&itemId=%d", payload.ProjectsV2Item.ID),
	}, nil
}

func ParseProjectV2ItemReordered(payload ProjectV2ItemEditedPayload) (*ParsedProjectV2ItemInfo, error) {
	// Replace with real GraphQL-based enrichment if needed
	return &ParsedProjectV2ItemInfo{
		CardTitle:     "Card Title Placeholder",
		ProjectNumber: "1",
		ColumnName:    "Doing", // this may be extracted via enrichment logic
		CardURL:       fmt.Sprintf("https://github.com/orgs/PMJ-Project/projects/1/views/1?pane=issue&itemId=%d", payload.ProjectsV2Item.ID),
	}, nil
}
