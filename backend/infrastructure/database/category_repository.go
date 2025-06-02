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

type categoryRepository struct {
	supabase *SupabaseClient
}

func NewCategoryRepository(supabase *SupabaseClient) repository.CategoryRepository {
	return &categoryRepository{supabase}
}

func (r *categoryRepository) Create(ctx context.Context, c *entity.Category) (*entity.Category, error) {
	c.ID = uuid.NewString()
	toInsert := map[string]interface{}{
		"id":          c.ID,
		"name":        c.Name,
		"color":       c.Color,
		"description": c.Description,
		"user_id":     c.UserID,
	}

	builder := r.supabase.DB.
		From("categories").
		Insert(toInsert, false, "", "*", "").
		Single()

	raw, _, err := builder.Execute()
	if err != nil {
		return nil, err
	}

	var m model.CategoryModel
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return model.ToDomainCategory(&m), nil
}

func (r *categoryRepository) FindByID(ctx context.Context, id string) (*entity.Category, error) {
	builder := r.supabase.DB.
		From("categories").
		Select("*", "", false).
		Eq("id", id).
		Single()

	raw, _, err := builder.Execute()
	if err != nil {
		return nil, err
	}

	var m model.CategoryModel
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return model.ToDomainCategory(&m), nil
}

func (r *categoryRepository) FindAllByUser(ctx context.Context, userID string) ([]*entity.Category, error) {
	builder := r.supabase.DB.
		From("categories").
		Select("*", "", false).
		Eq("user_id", userID)

	raw, _, err := builder.Execute()
	if err != nil {
		return nil, err
	}

	var models []model.CategoryModel
	if err := json.Unmarshal(raw, &models); err != nil {
		return nil, err
	}

	var cats []*entity.Category
	for _, m := range models {
		cats = append(cats, model.ToDomainCategory(&m))
	}
	return cats, nil
}

func (r *categoryRepository) Update(ctx context.Context, c *entity.Category) (*entity.Category, error) {
	if c.ID == "" {
		return nil, errors.New("category ID is required")
	}
	updates := map[string]interface{}{}
	if c.Name != "" {
		updates["name"] = c.Name
	}
	if c.Color != "" {
		updates["color"] = c.Color
	}
	if c.Description != nil {
		updates["description"] = c.Description
	}

	builder := r.supabase.DB.
		From("categories").
		Update(updates, "*", "").
		Eq("id", c.ID).
		Single()

	raw, _, err := builder.Execute()
	if err != nil {
		return nil, err
	}

	var m model.CategoryModel
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return model.ToDomainCategory(&m), nil
}

func (r *categoryRepository) Delete(ctx context.Context, id string) error {
	builder := r.supabase.DB.
		From("categories").
		Delete("*", "").
		Eq("id", id)

	_, _, err := builder.Execute()
	return err
}
