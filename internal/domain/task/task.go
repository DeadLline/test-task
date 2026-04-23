package task

import "time"

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type RecurrenceType string

const (
	RecurrenceDaily   RecurrenceType = "daily"
	RecurrenceMonthly RecurrenceType = "monthly"
	RecurrenceDates   RecurrenceType = "dates"
	RecurrenceEvenOdd RecurrenceType = "even_odd"
)

type RecurrenceRule struct {
	Type       RecurrenceType `json:"type"`
	EveryNDays int            `json:"every_n_days,omitempty"`
	DayOfMonth int            `json:"day_of_month,omitempty"`
	Dates      []time.Time    `json:"dates,omitempty"`
	Even       *bool          `json:"even,omitempty"`
}

type Task struct {
	ID          int64           `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Status      Status          `json:"status"`

	ScheduledAt time.Time       `json:"scheduled_at"`

	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`

	Recurrence  *RecurrenceRule `json:"recurrence,omitempty"`
	IsTemplate  bool            `json:"is_template"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}