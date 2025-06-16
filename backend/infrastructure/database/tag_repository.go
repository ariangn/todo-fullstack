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

type tagRepository struct {
	supabase *SupabaseClient
}

func NewTagRepository(supabase *SupabaseClient) repository.TagRepository {
	return &tagRepository{supabase}
}

func (r *tagRepository) FindByName(ctx context.Context, userID string, name string) (*entity.Tag, error) {
	builder := r.supabase.DB.
		From("tags").
		Select("*", "", false).
		Eq("user_id", userID).
		Eq("name", name).
		Single()

	raw, _, err := builder.Execute()
	if err != nil {
		// not found is not a fatal error
		return nil, nil
	}

	var m model.TagModel
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}

	return model.ToDomainTag(&m), nil
}

func (r *tagRepository) Create(ctx context.Context, t *entity.Tag) (*entity.Tag, error) {
	if t == nil {
		return nil, errors.New("tag entity is required")
	}
	id := t.ID
	if id == "" {
		id = uuid.NewString()
	}
	insert := map[string]interface{}{
		"id":      id,
		"user_id": t.UserID,
		"name":    t.Name,
	}

	_, _, err := r.supabase.DB.
		From("tags").
		Insert(insert, false, "", "", ""). // ← don't expect any data back
		Execute()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *tagRepository) FindByID(ctx context.Context, id string) (*entity.Tag, error) {
	builder := r.supabase.DB.
		From("tags").
		Select("*", "", false).
		Eq("id", id).
		Single()

	raw, _, err := builder.Execute()
	if err != nil {
		return nil, err
	}

	var m model.TagModel
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return model.ToDomainTag(&m), nil
}

func (r *tagRepository) FindAllByUser(ctx context.Context, userID string) ([]*entity.Tag, error) {
	builder := r.supabase.DB.
		From("tags").
		Select("*", "", false).
		Eq("user_id", userID)

	raw, _, err := builder.Execute()
	if err != nil {
		return nil, err
	}

	var models []model.TagModel
	if err := json.Unmarshal(raw, &models); err != nil {
		return nil, err
	}

	var tags []*entity.Tag
	for _, m := range models {
		tags = append(tags, model.ToDomainTag(&m))
	}
	return tags, nil
}

func (r *tagRepository) Update(ctx context.Context, t *entity.Tag) (*entity.Tag, error) {
	if t.ID == "" {
		return nil, errors.New("tag ID is required")
	}
	updates := map[string]interface{}{}
	if t.Name != "" {
		updates["name"] = t.Name
	}

	builder := r.supabase.DB.
		From("tags").
		Update(updates, "*", "").
		Eq("id", t.ID).
		Single()

	raw, _, err := builder.Execute()
	if err != nil {
		return nil, err
	}

	var m model.TagModel
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return model.ToDomainTag(&m), nil
}

func (r *tagRepository) Delete(ctx context.Context, id string) error {
	builder := r.supabase.DB.
		From("tags").
		Delete("*", "").
		Eq("id", id)

	_, _, err := builder.Execute()
	return err
}
