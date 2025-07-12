package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v50/github"
	"github.com/taujago/go-github-telegram-bot/github"
	"github.com/taujago/go-github-telegram-bot/telegram"
	"go.uber.org/zap"
)

type Handler struct {
	github   *github.Client
	telegram *telegram.Client
	logger   *zap.Logger
}

func NewHandler(github *github.Client, telegram *telegram.Client, logger *zap.Logger) *Handler {
	return &Handler{
		github:   github,
		telegram: telegram,
		logger:   logger,
	}
}

func (h *Handler) HandleGitHubWebhook(c *gin.Context) {
	eventType := c.GetHeader("X-GitHub-Event")
	h.logger.Info("Received GitHub event", zap.String("type", eventType))

	switch eventType {
	case "project_card":
		var event github.ProjectCardEvent
		if err := c.ShouldBindJSON(&event); err != nil {
			h.logger.Error("Failed to parse project_card event", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}

		h.logger.Debug("Project card event",
			zap.String("action", event.GetAction()),
			zap.Int64("card_id", event.GetProjectCard().GetID()),
		)

		switch event.GetAction() {
		case "moved":
			if err := h.handleCardMovement(event); err != nil {
				h.logger.Error("Failed to handle card movement", zap.Error(err))
			}
		case "created":
			h.logger.Debug("Card created event received")
		case "deleted":
			h.logger.Debug("Card deleted event received")
		default:
			h.logger.Debug("Unhandled project card action", zap.String("action", event.GetAction()))
		}
	default:
		h.logger.Debug("Unhandled event type", zap.String("type", eventType))
	}

	c.Status(http.StatusOK)
}

func (h *Handler) handleCardMovement(event github.ProjectCardEvent) error {
	card := event.GetProjectCard()
	changes := event.GetChanges()

	if changes.ColumnID == nil || changes.ColumnID.From == nil {
		return fmt.Errorf("missing column change information in payload")
	}

	previousColumnID := *changes.ColumnID.From
	newColumnID := card.GetColumnID()

	// Get column details
	previousColumn, err := h.github.GetProjectColumn(previousColumnID)
	if err != nil {
		return fmt.Errorf("failed to get previous column: %w", err)
	}

	newColumn, err := h.github.GetProjectColumn(newColumnID)
	if err != nil {
		return fmt.Errorf("failed to get new column: %w", err)
	}

	// Get card content if available
	cardContent := card.GetNote()
	if contentURL := card.GetContentURL(); contentURL != "" {
		if issue, err := h.github.GetIssueFromURL(contentURL); err == nil {
			cardContent = fmt.Sprintf("#%d: %s", issue.GetNumber(), issue.GetTitle())
		}
	}

	// Prepare and send message
	message := fmt.Sprintf(
		"🃏 *Card Moved*\n\n"+
			"• *Card*: %s\n"+
			"• *From*: %s\n"+
			"• *To*: %s\n"+
			"• *Project*: %s",
		cardContent,
		previousColumn.GetName(),
		newColumn.GetName(),
		previousColumn.GetProjectURL(), // You might want to get actual project name
	)

	return h.telegram.SendMessage(message)
}
