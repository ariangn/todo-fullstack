package database

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"

	"github.com/ariangn/todo-fullstack/backend/domain/entity"
	"github.com/ariangn/todo-fullstack/backend/domain/repository"
	"github.com/ariangn/todo-fullstack/backend/infrastructure/database/model"
)

type todoRepository struct {
	supabase *SupabaseClient
}

func NewTodoRepository(supabase *SupabaseClient) repository.TodoRepository {
	return &todoRepository{supabase}
}

func (r *todoRepository) Create(ctx context.Context, t *entity.Todo) (*entity.Todo, error) {
	t.ID = uuid.NewString()

	toInsert := map[string]interface{}{
		"id":           t.ID,
		"title":        t.Title,
		"body":         t.Body,
		"status":       string(t.Status),
		"due_date":     t.DueDate,
		"completed_at": t.CompletedAt,
		"user_id":      t.UserID,
		"category_id":  t.CategoryID,
		// Note: we do not insert TagIDs here; handle tags via a separate join if needed
	}

	builder := r.supabase.DB.
		From("todos").
		Insert(toInsert, false, "", "*", "").
		Single()

	raw, _, err := builder.Execute()
	if err != nil {
		return nil, err
	}

	var m model.TodoModel
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return model.ToDomainTodo(&m), nil
}

func (r *todoRepository) FindByID(ctx context.Context, id string) (*entity.Todo, error) {
	// Assumes a view “todos_with_tag_ids” exists that aggregates tag_ids.
	builder := r.supabase.DB.
		From("todos_with_tag_ids").
		Select("*", "", false).
		Eq("id", id).
		Single()

	raw, _, err := builder.Execute()
	if err != nil {
		return nil, err
	}

	var m model.TodoModel
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return model.ToDomainTodo(&m), nil
}

func (r *todoRepository) FindAllByUser(ctx context.Context, userID string) ([]*entity.Todo, error) {
	builder := r.supabase.DB.
		From("todos_with_tag_ids").
		Select("*", "", false).
		Eq("user_id", userID)

	raw, _, err := builder.Execute()
	if err != nil {
		return nil, err
	}

	var models []model.TodoModel
	if err := json.Unmarshal(raw, &models); err != nil {
		return nil, err
	}

	var todos []*entity.Todo
	for _, m := range models {
		todos = append(todos, model.ToDomainTodo(&m))
	}
	return todos, nil
}

func (r *todoRepository) Update(ctx context.Context, t *entity.Todo) (*entity.Todo, error) {
	if t.ID == "" {
		return nil, errors.New("todo ID is required")
	}
	updates := map[string]interface{}{}
	if t.Title != "" {
		updates["title"] = t.Title
	}
	if t.Body != nil {
		updates["body"] = t.Body
	}
	if t.Status != "" {
		updates["status"] = string(t.Status)
	}
	if t.DueDate != nil {
		updates["due_date"] = t.DueDate
	}
	if t.CompletedAt != nil {
		updates["completed_at"] = t.CompletedAt
	}
	if t.CategoryID != nil {
		updates["category_id"] = t.CategoryID
	}

	builder := r.supabase.DB.
		From("todos").
		Update(updates, "*", "").
		Eq("id", t.ID).
		Single()

	raw, _, err := builder.Execute()
	if err != nil {
		return nil, err
	}

	var m model.TodoModel
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return model.ToDomainTodo(&m), nil
}

func (r *todoRepository) Delete(ctx context.Context, id string) error {
	builder := r.supabase.DB.
		From("todos").
		Delete("*", "").
		Eq("id", id)

	_, _, err := builder.Execute()
	return err
}
