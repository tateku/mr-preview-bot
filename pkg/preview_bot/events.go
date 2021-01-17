package preview_bot

import (
	"time"
)

type ProjectID int
type MergeRequestID int
type PipelineID int
type PreviewID int

type PreviewEventBaseFields struct {
	PreviewID PreviewID
	PreviewProjectID ProjectID
	MergeRequestID MergeRequestID
	DateTime time.Time
}

type PreviewStartedEvent struct {
	PreviewEventBaseFields
}

type PreviewReadyEvent struct {
	PreviewEventBaseFields
	Info struct {
		URLs map[string]string
	}
}

type PreviewFailedEvent struct {
	PreviewEventBaseFields
}

type PreviewDeletedEvent struct {
	PreviewEventBaseFields
}

type PreviewCancelledEvent struct {
	PreviewEventBaseFields
}

type PipelineFinalizedStatus string

const (
	PipelineFinalizedStatusFailed PipelineFinalizedStatus = "failed"
	PipelineFinalizedStatusSuccess PipelineFinalizedStatus = "success"
)

type PipelineStartedEvent struct {
	PipelineID PipelineID
	ProjectID ProjectID
	MergeRequestID MergeRequestID
	DateTime time.Time
}

type PipelineFinalizedEvent struct {
	PipelineID PipelineID
	ProjectID ProjectID
	MergeRequestID MergeRequestID
	DateTime time.Time
	Status PipelineFinalizedStatus
}

type MergeRequestStatus string

const (
	MergeRequestStatusOpened MergeRequestStatus = "opened"
	MergeRequestStatusClosed MergeRequestStatus = "closed"
	MergeRequestStatusLocked MergeRequestStatus = "locked"
	MergeRequestStatusMerged MergeRequestStatus = "merged"
)

type MergeRequestOpenedEvent struct {
	MergeRequestID MergeRequestID
	ProjectID ProjectID
	DateTime time.Time
}

type MergeRequestFinalizedEvent struct {
	MergeRequestID MergeRequestID
	ProjectID      ProjectID
	DateTime       time.Time
	Status         MergeRequestStatus
}
