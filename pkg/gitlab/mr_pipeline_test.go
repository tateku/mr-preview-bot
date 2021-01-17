package gitlab

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/mock"
	"mr-preview-bot/pkg/preview_bot"
	"testing"
	"time"
)

type MockClient struct{
	mock.Mock
}

func (m *MockClient) ListProjectMergeRequests(projectID int, updatedAt time.Time) ([]*MergeRequestInfo, error) {
	args := m.Called(projectID, updatedAt)

	return args.Get(0).([]*MergeRequestInfo), args.Error(1)
}

func TestMRPipelineAggregate_should_return_no_events_if_client_returns_no_event(t *testing.T) {
	projectID := preview_bot.ProjectID(42)
	mockClient := &MockClient{}
    aggregate := createMRPipelineAggregate(projectID, mockClient, nil)

	command := &PollMRPipelineCommand{Now: time.Unix(45, 0)}

	mergeRequestInfos := []*MergeRequestInfo(nil)
    mockClient.On("ListProjectMergeRequests", 42, time.Unix(15, 0)).Return(mergeRequestInfos, nil)

	expectedEvents := []interface{}(nil)

	got := aggregate.HandlePollMRPipelineCommand(command)
	diff := cmp.Diff(expectedEvents, got)
	if diff != "" {
		t.Fatalf(diff)
	}

	mockClient.AssertExpectations(t)
}

func TestMRPipelineAggregate_should_return_no_events_if_client_fails(t *testing.T) {
	projectID := preview_bot.ProjectID(42)
	mockClient := &MockClient{}
	aggregate := createMRPipelineAggregate(projectID, mockClient, nil)

	command := &PollMRPipelineCommand{Now: time.Unix(45, 0)}

	listErr := errors.New("list projects merge requests error")
	mergeRequestInfos := []*MergeRequestInfo(nil)
	mockClient.On("ListProjectMergeRequests", 42, time.Unix(15, 0)).Return(mergeRequestInfos, listErr)

	expectedEvents := []interface{}(nil)

	got := aggregate.HandlePollMRPipelineCommand(command)
	diff := cmp.Diff(expectedEvents, got)
	if diff != "" {
		t.Fatalf(diff)
	}

	mockClient.AssertExpectations(t)
}

func TestMRPipelineAggregate_should_fire_mr_opened_event(t *testing.T) {
	projectID := preview_bot.ProjectID(42)
	mockClient := &MockClient{}
	aggregate := createMRPipelineAggregate(projectID, mockClient, nil)

	command := &PollMRPipelineCommand{Now: time.Unix(45, 0)}

	mergeRequestInfos := []*MergeRequestInfo{
		{
			MergeRequestID: 73,
			CreatedAt:      time.Unix(123, 0),
			UpdatedAt:      time.Unix(456, 0),
			State:          preview_bot.MergeRequestStatusOpened,
		},
	}
	mockClient.On("ListProjectMergeRequests", 42, time.Unix(15, 0)).Return(mergeRequestInfos, nil)

	expectedEvents := []interface{}{
		&preview_bot.MergeRequestOpenedEvent{
			MergeRequestID: 73,
			ProjectID:      projectID,
			DateTime:       time.Unix(123, 0),
		},
	}

	got := aggregate.HandlePollMRPipelineCommand(command)
	diff := cmp.Diff(expectedEvents, got)
	if diff != "" {
		t.Fatalf(diff)
	}

	mockClient.AssertExpectations(t)
}

func TestMRPipelineAggregate_should_return_no_event_if_client_returns_already_closed_mr(t *testing.T) {
	projectID := preview_bot.ProjectID(42)
	mockClient := &MockClient{}

	aggregate := createMRPipelineAggregate(projectID, mockClient, nil)

	command := &PollMRPipelineCommand{Now: time.Unix(45, 0)}

	mergeRequestInfos := []*MergeRequestInfo{
		{
			MergeRequestID: 73,
			CreatedAt:      time.Unix(123, 0),
			UpdatedAt:      time.Unix(456, 0),
			State:          preview_bot.MergeRequestStatusClosed,
		},
	}
	mockClient.On("ListProjectMergeRequests", 42, time.Unix(15, 0)).Return(mergeRequestInfos, nil)


	expectedEvents := []interface{}(nil)

	got := aggregate.HandlePollMRPipelineCommand(command)
	diff := cmp.Diff(expectedEvents, got)
	if diff != "" {
		t.Fatalf(diff)
	}

	mockClient.AssertExpectations(t)
}

func TestMRPipelineAggregate_should_return_no_event_for_already_open_mr(t *testing.T) {
	projectID := preview_bot.ProjectID(42)
	mockClient := &MockClient{}
	previousEvents := []interface{}{
		&preview_bot.MergeRequestOpenedEvent{
			MergeRequestID: 73,
			ProjectID:      projectID,
			DateTime:       time.Unix(123, 0),
		},
	}
	aggregate := createMRPipelineAggregate(projectID, mockClient, previousEvents)

	command := &PollMRPipelineCommand{Now: time.Unix(45, 0)}

	mergeRequestInfos := []*MergeRequestInfo{
		{
			MergeRequestID: 73,
			CreatedAt:      time.Unix(123, 0),
			UpdatedAt:      time.Unix(456, 0),
			State:          preview_bot.MergeRequestStatusOpened,
		},
		{
			MergeRequestID: 74,
			CreatedAt:      time.Unix(678, 0),
			UpdatedAt:      time.Unix(456, 0),
			State:          preview_bot.MergeRequestStatusOpened,
		},
	}
	mockClient.On("ListProjectMergeRequests", 42, time.Unix(15, 0)).Return(mergeRequestInfos, nil)

	expectedEvents := []interface{}{
		&preview_bot.MergeRequestOpenedEvent{
			MergeRequestID: 74,
			ProjectID:      projectID,
			DateTime:       time.Unix(678, 0),
		},
	}

	got := aggregate.HandlePollMRPipelineCommand(command)
	diff := cmp.Diff(expectedEvents, got)
	if diff != "" {
		t.Fatalf(diff)
	}

	mockClient.AssertExpectations(t)
}

func TestMRPipelineAggregate_should_fire_mr_finalized_event(t *testing.T) {
	projectID := preview_bot.ProjectID(42)
	mockClient := &MockClient{}
	previousEvents := []interface{}{
		&preview_bot.MergeRequestOpenedEvent{
			MergeRequestID: 73,
			ProjectID:      projectID,
			DateTime:       time.Unix(123, 0),
		},
	}
	aggregate := createMRPipelineAggregate(projectID, mockClient, previousEvents)

	command := &PollMRPipelineCommand{Now: time.Unix(45, 0)}

	mergeRequestInfos := []*MergeRequestInfo{
		{
			MergeRequestID: 73,
			CreatedAt:      time.Unix(123, 0),
			UpdatedAt:      time.Unix(456, 0),
			State:          preview_bot.MergeRequestStatusClosed,
		},
	}
	mockClient.On("ListProjectMergeRequests", 42, time.Unix(15, 0)).Return(mergeRequestInfos, nil)

	expectedEvents := []interface{}{
		&preview_bot.MergeRequestFinalizedEvent{
			MergeRequestID: 73,
			ProjectID:      projectID,
			DateTime:       time.Unix(456, 0),
			Status: preview_bot.MergeRequestStatusClosed,
		},
	}

	got := aggregate.HandlePollMRPipelineCommand(command)
	diff := cmp.Diff(expectedEvents, got)
	if diff != "" {
		t.Fatalf(diff)
	}

	mockClient.AssertExpectations(t)
}

func TestMRPipelineAggregate_should_return_no_event_for_already_finalized_event(t *testing.T) {
	projectID := preview_bot.ProjectID(42)
	mockClient := &MockClient{}
	previousEvents := []interface{}{
		&preview_bot.MergeRequestOpenedEvent{
			MergeRequestID: 73,
			ProjectID:      projectID,
			DateTime:       time.Unix(123, 0),
		},
		&preview_bot.MergeRequestFinalizedEvent{
			MergeRequestID: 73,
			ProjectID:      projectID,
			DateTime:       time.Unix(456, 0),
			Status: preview_bot.MergeRequestStatusClosed,
		},
	}
	aggregate := createMRPipelineAggregate(projectID, mockClient, previousEvents)

	command := &PollMRPipelineCommand{Now: time.Unix(45, 0)}

	mergeRequestInfos := []*MergeRequestInfo{
		{
			MergeRequestID: 73,
			CreatedAt:      time.Unix(123, 0),
			UpdatedAt:      time.Unix(789, 0),
			State:          preview_bot.MergeRequestStatusClosed,
		},
	}
	mockClient.On("ListProjectMergeRequests", 42, time.Unix(15, 0)).Return(mergeRequestInfos, nil)

	expectedEvents := []interface{}(nil)

	got := aggregate.HandlePollMRPipelineCommand(command)
	diff := cmp.Diff(expectedEvents, got)
	if diff != "" {
		t.Fatalf(diff)
	}

	mockClient.AssertExpectations(t)
}

func createMRPipelineAggregate(projectID preview_bot.ProjectID, client client, previousEvents []interface{}) *MRPipelineAggregate {
	aggregate := NewMRPipelineAggregate(projectID, client)

	for _, event := range previousEvents {
		aggregate.HandleEvent(event)
	}

	return aggregate
}
