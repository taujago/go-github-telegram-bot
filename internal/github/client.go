package github

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type Client struct {
	client *github.Client
	ctx    context.Context
}

func NewClient(token string) *Client {
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))

	return &Client{
		client: github.NewClient(tc),
		ctx:    ctx,
	}
}

func (g *Client) GetProjectColumn(columnID int64) (*github.ProjectColumn, error) {
	column, resp, err := g.client.Projects.GetProjectColumn(g.ctx, columnID)
	if err != nil {
		return nil, fmt.Errorf("failed to get column %d: %w", columnID, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return column, nil
}

func (g *Client) GetIssueFromURL(url string) (*github.Issue, error) {
	parts := strings.Split(url, "/")
	if len(parts) < 5 {
		return nil, fmt.Errorf("invalid content URL format")
	}

	owner := parts[len(parts)-4]
	repo := parts[len(parts)-3]
	number, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return nil, fmt.Errorf("invalid issue number: %w", err)
	}

	issue, resp, err := g.client.Issues.Get(g.ctx, owner, repo, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return issue, nil
}

func (g *Client) GetProjectCard(cardID int64) (*github.ProjectCard, error) {
	card, resp, err := g.client.Projects.GetProjectCard(g.ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get card %d: %w", cardID, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return card, nil
}
