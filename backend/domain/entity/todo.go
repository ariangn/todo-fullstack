package entity

import (
	"errors"
	"time"
)

type Status string

const (
	StatusTodo       Status = "TODO"
	StatusInProgress Status = "IN_PROGRESS"
	StatusCompleted  Status = "COMPLETED"
)

type Todo struct {
	ID          string
	Title       string
	Body        *string
	Status      Status
	DueDate     *time.Time
	CompletedAt *time.Time
	UserID      string
	CategoryID  *string
	TagIDs      []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewTodo enforces: Title non-empty & dueDate (if set) not in past.
func NewTodo(
	id, title string,
	body *string,
	status Status,
	dueDate *time.Time,
	userID string,
	categoryID *string,
	tagIDs []string,
) (*Todo, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}
	return &Todo{
		ID:          id,
		Title:       title,
		Body:        body,
		Status:      status,
		DueDate:     dueDate,
		CompletedAt: nil,
		UserID:      userID,
		CategoryID:  categoryID,
		TagIDs:      tagIDs,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}, nil
}

var ErrDueDateInPast = errors.New("due date cannot be in the past")
