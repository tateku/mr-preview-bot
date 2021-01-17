package gitlab

import (
	"mr-preview-bot/pkg/preview_bot"
	"time"
)

// mr statuses: opened, closed, locked, merged.
// pipeline statuses: created, waiting_for_resource, preparing, pending, running, success, failed, canceled, skipped, manual, scheduled

type PollMRPipelineCommand struct {
	Now time.Time
}

type MRPipelineAggregate struct {
	ProjectID preview_bot.ProjectID
	client    client
	// internal state
	mrStates map[preview_bot.MergeRequestID]preview_bot.MergeRequestStatus
}

func (aggregate *MRPipelineAggregate) HandleEvent(event interface{}) {
	switch ev := event.(type) {
	case *preview_bot.MergeRequestOpenedEvent:
		{
			aggregate.handleMergeRequestOpenedEvent(ev)
		}
	case *preview_bot.MergeRequestFinalizedEvent:
		{
			aggregate.handleMergeRequestFinalizedEvent(ev)
		}
	}
}

func (aggregate *MRPipelineAggregate) handleMergeRequestOpenedEvent(event *preview_bot.MergeRequestOpenedEvent) {
	_, found := aggregate.mrStates[event.MergeRequestID]
	if found {
		// TODO: handle if opened event comes for already opened mr?
	}

	aggregate.mrStates[event.MergeRequestID] = preview_bot.MergeRequestStatusOpened
}

func (aggregate *MRPipelineAggregate) handleMergeRequestFinalizedEvent(event *preview_bot.MergeRequestFinalizedEvent) {
	_, found := aggregate.mrStates[event.MergeRequestID]
	if !found {
		// TODO: handle if closed event comes for non-opened
	}

	aggregate.mrStates[event.MergeRequestID] = event.Status
}

func (aggregate *MRPipelineAggregate) HandlePollMRPipelineCommand(command *PollMRPipelineCommand) []interface{} {
	thirtySecondsAgo := command.Now.Add(-time.Second * 30)
	// TODO: make sure we use last updated or thirty seconds ago
	mergeRequests, err := aggregate.client.ListProjectMergeRequests(int(aggregate.ProjectID), thirtySecondsAgo)
	if err != nil {
		return nil
	}

	var mergeRequestEvents []interface{}
	for _, mr := range mergeRequests {
		lastMRStatus, found := aggregate.mrStates[mr.MergeRequestID]

		if mr.State == preview_bot.MergeRequestStatusOpened {
			if found {
				continue
			}
			mergeRequestEvents = append(mergeRequestEvents, &preview_bot.MergeRequestOpenedEvent{
				MergeRequestID: mr.MergeRequestID,
				ProjectID:      aggregate.ProjectID,
				DateTime:       mr.CreatedAt,
			})
		} else if found && lastMRStatus == preview_bot.MergeRequestStatusOpened {
			mergeRequestEvents = append(mergeRequestEvents, &preview_bot.MergeRequestFinalizedEvent{
				MergeRequestID: mr.MergeRequestID,
				ProjectID:      aggregate.ProjectID,
				DateTime:       mr.UpdatedAt,
				Status: mr.State,
			})
		}
	}

	return mergeRequestEvents
}

func NewMRPipelineAggregate(projectID preview_bot.ProjectID, client client) *MRPipelineAggregate {
	return &MRPipelineAggregate{
		ProjectID: projectID,
		client:    client,
		mrStates:  make(map[preview_bot.MergeRequestID]preview_bot.MergeRequestStatus),
	}
}
