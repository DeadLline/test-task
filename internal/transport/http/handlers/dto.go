package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type RecurrenceDTO struct {
	Type       string      `json:"type"`
	EveryNDays int         `json:"every_n_days,omitempty"`
	DayOfMonth int         `json:"day_of_month,omitempty"`
	Dates      []time.Time `json:"dates,omitempty"`
	Even       *bool       `json:"even,omitempty"`
}

type taskMutationDTO struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`

	ScheduledAt time.Time      `json:"scheduled_at"`
	Recurrence  *RecurrenceDTO `json:"recurrence"`
	IsTemplate  bool           `json:"is_template"`
}

type taskDTO struct {
	ID          int64             `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`

	ScheduledAt time.Time `json:"scheduled_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Recurrence *RecurrenceDTO `json:"recurrence,omitempty"`
	IsTemplate bool           `json:"is_template"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	return taskDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		ScheduledAt: task.ScheduledAt,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		Recurrence:  mapRecurrenceToDTO(task.Recurrence),
		IsTemplate:  task.IsTemplate,
	}
}

func mapRecurrenceToDTO(r *taskdomain.RecurrenceRule) *RecurrenceDTO {
	if r == nil {
		return nil
	}
	return &RecurrenceDTO{
		Type:       string(r.Type),
		EveryNDays: r.EveryNDays,
		DayOfMonth: r.DayOfMonth,
		Dates:      r.Dates,
		Even:       r.Even,
	}
}
