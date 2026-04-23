package task

import (
	"context"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, in CreateInput) (*taskdomain.Task, error) {
	in.Title = strings.TrimSpace(in.Title)
	in.Description = strings.TrimSpace(in.Description)

	if in.Title == "" {
		return nil, ErrInvalidInput
	}

	if in.Status == "" {
		in.Status = taskdomain.StatusNew
	}

	if !in.Status.Valid() {
		return nil, ErrInvalidInput
	}

	if err := validateRecurrence(in.Recurrence); err != nil {
		return nil, err
	}

	now := s.now()

	t := &taskdomain.Task{
		Title:       in.Title,
		Description: in.Description,
		Status:      in.Status,
		ScheduledAt: in.ScheduledAt,
		Recurrence:  in.Recurrence,
		IsTemplate:  in.IsTemplate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return s.repo.Create(ctx, t)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, ErrInvalidInput
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, in UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, ErrInvalidInput
	}

	in.Title = strings.TrimSpace(in.Title)
	in.Description = strings.TrimSpace(in.Description)

	if in.Title == "" {
		return nil, ErrInvalidInput
	}

	if !in.Status.Valid() {
		return nil, ErrInvalidInput
	}

	if err := validateRecurrence(in.Recurrence); err != nil {
		return nil, err
	}

	t := &taskdomain.Task{
		ID:          id,
		Title:       in.Title,
		Description: in.Description,
		Status:      in.Status,
		ScheduledAt: in.ScheduledAt,
		Recurrence:  in.Recurrence,
		IsTemplate:  in.IsTemplate,
		UpdatedAt:   s.now(),
	}

	return s.repo.Update(ctx, t)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidInput
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	return s.repo.List(ctx)
}

func (s *Service) GenerateTasks(ctx context.Context) error {
	tasks, err := s.repo.List(ctx)
	if err != nil {
		return err
	}

	now := s.now()

	for _, t := range tasks {
		if !t.IsTemplate || t.Recurrence == nil {
			continue
		}

		if shouldCreate(t, now) {
			newTask := t
			newTask.ID = 0
			newTask.IsTemplate = false
			newTask.CreatedAt = now
			newTask.UpdatedAt = now

			_, err := s.repo.Create(ctx, &newTask)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func shouldCreate(t taskdomain.Task, now time.Time) bool {
	r := t.Recurrence

	switch r.Type {

	case taskdomain.RecurrenceDaily:
		days := int(now.Sub(t.CreatedAt).Hours() / 24)
		if r.EveryNDays <= 0 {
			return false
		}
		return days%r.EveryNDays == 0

	case taskdomain.RecurrenceMonthly:
		return now.Day() == r.DayOfMonth

	case taskdomain.RecurrenceDates:
		for _, d := range r.Dates {
			if sameDay(d, now) {
				return true
			}
		}

	case taskdomain.RecurrenceEvenOdd:
		if r.Even == nil {
			return false
		}
		if *r.Even {
			return now.Day()%2 == 0
		}
		return now.Day()%2 == 1
	}

	return false
}

func sameDay(a, b time.Time) bool {
	return a.Year() == b.Year() &&
		a.Month() == b.Month() &&
		a.Day() == b.Day()
}

func validateRecurrence(r *taskdomain.RecurrenceRule) error {
	if r == nil {
		return nil
	}

	switch r.Type {

	case taskdomain.RecurrenceDaily:
		if r.EveryNDays <= 0 {
			return ErrInvalidInput
		}

	case taskdomain.RecurrenceMonthly:
		if r.DayOfMonth < 1 || r.DayOfMonth > 30 {
			return ErrInvalidInput
		}

	case taskdomain.RecurrenceDates:
		if len(r.Dates) == 0 {
			return ErrInvalidInput
		}

	case taskdomain.RecurrenceEvenOdd:
		if r.Even == nil {
			return ErrInvalidInput
		}

	default:
		return ErrInvalidInput
	}

	return nil
}