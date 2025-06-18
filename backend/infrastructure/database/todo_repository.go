package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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
	}

	// Step 1: Insert todo using return=minimal (no unmarshalling needed)
	_, _, err := r.supabase.DB.
		From("todos").
		Insert(toInsert, false, "", "minimal", ""). // ← use return=minimal
		Single().
		Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to insert todo: %w", err)
	}

	// Step 2: Insert into todo_tags
	for _, tagID := range t.TagIDs {
		join := map[string]interface{}{
			"todo_id": t.ID,
			"tag_id":  tagID,
		}
		_, _, err := r.supabase.DB.
			From("todo_tags").
			Insert(join, false, "", "minimal", ""). // also minimal here
			Execute()
		if err != nil {
			return nil, fmt.Errorf("failed to insert todo_tag: %w", err)
		}
	}

	// Step 3: Return the original entity
	return t, nil
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

	// 1) Update core todo fields
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
	if len(updates) > 0 {
		if _, _, err := r.supabase.DB.
			From("todos").
			Update(updates, "", "").
			Eq("id", t.ID).
			Execute(); err != nil {
			return nil, err
		}
	}

	// 2) Only touch tags if caller provided TagIDs
	if t.TagIDs != nil {
		// a) sanitize new list (drop empty)
		newIDs := make(map[string]struct{}, len(t.TagIDs))
		for _, tid := range t.TagIDs {
			if s := strings.TrimSpace(tid); s != "" {
				newIDs[s] = struct{}{}
			}
		}

		// b) fetch existing tag_ids for this todo
		raw, _, err := r.supabase.DB.
			From("todo_tags").
			Select("tag_id", "", false).
			Eq("todo_id", t.ID).
			Execute()
		if err != nil {
			return nil, err
		}
		var rows []struct {
			TagID string `json:"tag_id"`
		}
		if err := json.Unmarshal(raw, &rows); err != nil {
			return nil, err
		}
		oldIDs := make(map[string]struct{}, len(rows))
		for _, r := range rows {
			oldIDs[r.TagID] = struct{}{}
		}

		// c) compute diffs
		toAdd := []map[string]interface{}{}
		for id := range newIDs {
			if _, seen := oldIDs[id]; !seen {
				toAdd = append(toAdd, map[string]interface{}{
					"todo_id": t.ID,
					"tag_id":  id,
				})
			}
		}
		toRemove := []string{}
		for id := range oldIDs {
			if _, keep := newIDs[id]; !keep {
				toRemove = append(toRemove, id)
			}
		}

		// d) batch delete removed tags
		if len(toRemove) > 0 {
			if _, _, err := r.supabase.DB.
				From("todo_tags").
				Delete("", "").
				Eq("todo_id", t.ID).
				In("tag_id", toRemove).
				Execute(); err != nil {
				return nil, err
			}
		}

		// e) batch insert added tags
		if len(toAdd) > 0 {
			if _, _, err := r.supabase.DB.
				From("todo_tags").
				Insert(toAdd, false, "", "", "").
				Execute(); err != nil {
				return nil, err
			}
		}
	}

	// 3) Return the fresh todo
	return r.FindByID(ctx, t.ID)
}

func (r *todoRepository) Delete(ctx context.Context, id string) error {
	builder := r.supabase.DB.
		From("todos").
		Delete("*", "").
		Eq("id", id)

	_, _, err := builder.Execute()
	return err
}
