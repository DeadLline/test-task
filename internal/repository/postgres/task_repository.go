package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	const query = `
	INSERT INTO tasks (title, description, status, scheduled_at, created_at, updated_at, recurrence, is_template)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	RETURNING id, title, description, status, scheduled_at, created_at, updated_at, recurrence, is_template
	`

	recurrence, _ := json.Marshal(task.Recurrence)

	row := r.pool.QueryRow(ctx, query,
		task.Title,
		task.Description,
		task.Status,
		task.ScheduledAt,
		task.CreatedAt,
		task.UpdatedAt,
		recurrence,
		task.IsTemplate,
	)

	return scanTask(row)
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	const query = `
	SELECT id, title, description, status, scheduled_at, created_at, updated_at, recurrence, is_template
	FROM tasks WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	task, err := scanTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrNotFound
		}
		return nil, err
	}
	return task, nil
}

func (r *Repository) Update(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	const query = `
	UPDATE tasks
	SET title=$1, description=$2, status=$3, scheduled_at=$4, updated_at=$5, recurrence=$6, is_template=$7
	WHERE id=$8
	RETURNING id, title, description, status, scheduled_at, created_at, updated_at, recurrence, is_template
	`

	recurrence, _ := json.Marshal(task.Recurrence)

	row := r.pool.QueryRow(ctx, query,
		task.Title,
		task.Description,
		task.Status,
		task.ScheduledAt,
		task.UpdatedAt,
		recurrence,
		task.IsTemplate,
		task.ID,
	)

	return scanTask(row)
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM tasks WHERE id = $1`
	res, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return taskdomain.ErrNotFound
	}
	return nil
}

func (r *Repository) List(ctx context.Context) ([]taskdomain.Task, error) {
	const query = `
	SELECT id, title, description, status, scheduled_at, created_at, updated_at, recurrence, is_template
	FROM tasks ORDER BY id DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []taskdomain.Task
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *t)
	}
	return tasks, rows.Err()
}

type taskScanner interface {
	Scan(dest ...any) error
}

func scanTask(scanner taskScanner) (*taskdomain.Task, error) {
	var task taskdomain.Task
	var status string
	var recurrence []byte

	err := scanner.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&status,
		&task.ScheduledAt,
		&task.CreatedAt,
		&task.UpdatedAt,
		&recurrence,
		&task.IsTemplate,
	)
	if err != nil {
		return nil, err
	}

	task.Status = taskdomain.Status(status)

	if len(recurrence) > 0 {
		var r taskdomain.RecurrenceRule
		_ = json.Unmarshal(recurrence, &r)
		task.Recurrence = &r
	}

	return &task, nil
}