package database

import (
    "context"
    "errors"

    "github.com/ariangn/todo-go/domain/entity"
    "github.com/ariangn/todo-go/domain/repository"
    "github.com/ariangn/todo-go/infrastructure/database/model"
    "github.com/google/uuid"
)

type categoryRepository struct {
    supabase *SupabaseClient
}

func NewCategoryRepository(supabase *SupabaseClient) repository.CategoryRepository {
    return &categoryRepository{supabase}
}

func (r *categoryRepository) Create(ctx context.Context, c *entity.Category) (*entity.Category, error) {
    c.ID = uuid.NewString()
    row := map[string]interface{}{
        "id":          c.ID,
        "name":        c.Name,
        "color":       c.Color,
        "description": c.Description,
        "user_id":     c.UserID,
    }
    resp, err := r.supabase.DB.
        From("categories").
        Insert(row).
        Single().
        Execute(ctx)
    if err != nil {
        return nil, err
    }
    var m model.CategoryModel
    if err := resp.JSON(&m); err != nil {
        return nil, err
    }
    return model.ToDomainCategory(&m), nil
}

func (r *categoryRepository) FindByID(ctx context.Context, id string) (*entity.Category, error) {
    resp, err := r.supabase.DB.
        From("categories").
        Select("*").
        Eq("id", id).
        Single().
        Execute(ctx)
    if err != nil {
        return nil, err
    }
    var m model.CategoryModel
    if err := resp.JSON(&m); err != nil {
        return nil, err
    }
    return model.ToDomainCategory(&m), nil
}

func (r *categoryRepository) FindAllByUser(ctx context.Context, userID string) ([]*entity.Category, error) {
    resp, err := r.supabase.DB.
        From("categories").
        Select("*").
        Eq("user_id", userID).
        Execute(ctx)
    if err != nil {
        return nil, err
    }
    var models []model.CategoryModel
    if err := resp.JSON(&models); err != nil {
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
    resp, err := r.supabase.DB.
        From("categories").
        Update(updates).
        Eq("id", c.ID).
        Single().
        Execute(ctx)
    if err != nil {
        return nil, err
    }
    var m model.CategoryModel
    if err := resp.JSON(&m); err != nil {
        return nil, err
    }
    return model.ToDomainCategory(&m), nil
}

func (r *categoryRepository) Delete(ctx context.Context, id string) error {
    _, err := r.supabase.DB.
        From("categories").
        Delete().
        Eq("id", id).
        Execute(ctx)
    return err
}
