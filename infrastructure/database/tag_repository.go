package database

import (
    "context"
    "errors"

    "github.com/ariangn/todo-go/domain/entity"
    "github.com/ariangn/todo-go/domain/repository"
    "github.com/ariangn/todo-go/infrastructure/database/model"
    "github.com/google/uuid"
)

type tagRepository struct {
    supabase *SupabaseClient
}

func NewTagRepository(supabase *SupabaseClient) repository.TagRepository {
    return &tagRepository{supabase}
}

func (r *tagRepository) Create(ctx context.Context, t *entity.Tag) (*entity.Tag, error) {
    t.ID = uuid.NewString()
    row := map[string]interface{}{
        "id":      t.ID,
        "name":    t.Name,
        "user_id": t.UserID,
    }
    resp, err := r.supabase.DB.
        From("tags").
        Insert(row).
        Single().
        Execute(ctx)
    if err != nil {
        return nil, err
    }
    var m model.TagModel
    if err := resp.JSON(&m); err != nil {
        return nil, err
    }
    return model.ToDomainTag(&m), nil
}

func (r *tagRepository) FindByID(ctx context.Context, id string) (*entity.Tag, error) {
    resp, err := r.supabase.DB.
        From("tags").
        Select("*").
        Eq("id", id).
        Single().
        Execute(ctx)
    if err != nil {
        return nil, err
    }
    var m model.TagModel
    if err := resp.JSON(&m); err != nil {
        return nil, err
    }
    return model.ToDomainTag(&m), nil
}

func (r *tagRepository) FindAllByUser(ctx context.Context, userID string) ([]*entity.Tag, error) {
    resp, err := r.supabase.DB.
        From("tags").
        Select("*").
        Eq("user_id", userID).
        Execute(ctx)
    if err != nil {
        return nil, err
    }
    var models []model.TagModel
    if err := resp.JSON(&models); err != nil {
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
    resp, err := r.supabase.DB.
        From("tags").
        Update(updates).
        Eq("id", t.ID).
        Single().
        Execute(ctx)
    if err != nil {
        return nil, err
    }
    var m model.TagModel
    if err := resp.JSON(&m); err != nil {
        return nil, err
    }
    return model.ToDomainTag(&m), nil
}

func (r *tagRepository) Delete(ctx context.Context, id string) error {
    _, err := r.supabase.DB.
        From("tags").
        Delete().
        Eq("id", id).
        Execute(ctx)
    return err
}
