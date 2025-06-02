package model

import (
	"time"

	"github.com/ariangn/todo-fullstack/backend/domain/entity"
)

// mirrors the JSON for "todos" table
type TodoModel struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Body        *string    `json:"body"`
	Status      string     `json:"status"`
	DueDate     *time.Time `json:"due_date"`
	CompletedAt *time.Time `json:"completed_at"`
	UserID      string     `json:"user_id"`
	CategoryID  *string    `json:"category_id"`
	TagIDs      []string   `json:"tag_ids"` // assuming a computed JSON array of tag_ids
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func ToDomainTodo(m *TodoModel) *entity.Todo {
	return &entity.Todo{
		ID:          m.ID,
		Title:       m.Title,
		Body:        m.Body,
		Status:      entity.Status(m.Status),
		DueDate:     m.DueDate,
		CompletedAt: m.CompletedAt,
		UserID:      m.UserID,
		CategoryID:  m.CategoryID,
		TagIDs:      m.TagIDs,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func FromDomainTodo(t *entity.Todo) *TodoModel {
	return &TodoModel{
		ID:          t.ID,
		Title:       t.Title,
		Body:        t.Body,
		Status:      string(t.Status),
		DueDate:     t.DueDate,
		CompletedAt: t.CompletedAt,
		UserID:      t.UserID,
		CategoryID:  t.CategoryID,
		TagIDs:      t.TagIDs,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
