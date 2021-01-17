package gitlab

import (
	"mr-preview-bot/pkg/preview_bot"
	"time"
)

type MergeRequestInfo struct {
	MergeRequestID preview_bot.MergeRequestID
	CreatedAt time.Time
	UpdatedAt time.Time
	State preview_bot.MergeRequestStatus
}

type client interface {
	ListProjectMergeRequests(projectID int, updatedAt time.Time) ([]*MergeRequestInfo, error)
}
