package database

import (
    "context"
    "errors"

    "github.com/ariangn/todo-go/domain/entity"
    "github.com/ariangn/todo-go/domain/repository"
    "github.com/ariangn/todo-go/infrastructure/database/model"
    "github.com/google/uuid"
)

type todoRepository struct {
    supabase *SupabaseClient
}

func NewTodoRepository(supabase *SupabaseClient) repository.TodoRepository {
    return &todoRepository{supabase}
}

func (r *todoRepository) Create(ctx context.Context, t *entity.Todo) (*entity.Todo, error) {
    t.ID = uuid.NewString()
    row := map[string]interface{}{
        "id":            t.ID,
        "title":         t.Title,
        "body":          t.Body,
        "status":        string(t.Status),
        "due_date":      t.DueDate,
        "completed_at":  t.CompletedAt,
        "user_id":       t.UserID,
        "category_id":   t.CategoryID,
        // `tag_ids` cannot be inserted directly; must insert into todo_tags join table separately.
        // For simplicity, ignore tags on initial insert.
    }

    resp, err := r.supabase.DB.
        From("todos").
        Insert(row).
        Single().
        Execute(ctx)
    if err != nil {
        return nil, err
    }
    var m model.TodoModel
    if err := resp.JSON(&m); err != nil {
        return nil, err
    }
    return model.ToDomainTodo(&m), nil
}

func (r *todoRepository) FindByID(ctx context.Context, id string) (*entity.Todo, error) {
    resp, err := r.supabase.DB.
        From("todos_with_tag_ids"). // assume a view that includes tag_ids array
        Select("*").
        Eq("id", id).
        Single().
        Execute(ctx)
    if err != nil {
        return nil, err
    }
    var m model.TodoModel
    if err := resp.JSON(&m); err != nil {
        return nil, err
    }
    return model.ToDomainTodo(&m), nil
}

func (r *todoRepository) FindAllByUser(ctx context.Context, userID string) ([]*entity.Todo, error) {
    resp, err := r.supabase.DB.
        From("todos_with_tag_ids").
        Select("*").
        Eq("user_id", userID).
        Execute(ctx)
    if err != nil {
        return nil, err
    }
    var models []model.TodoModel
    if err := resp.JSON(&models); err != nil {
        return nil, err
    }
    var todos []*entity.Todo
    for _, m := range models {
        d := model.ToDomainTodo(&m)
        todos = append(todos, d)
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
    // Changing category if provided
    if t.CategoryID != nil {
        updates["category_id"] = t.CategoryID
    }

    resp, err := r.supabase.DB.
        From("todos").
        Update(updates).
        Eq("id", t.ID).
        Single().
        Execute(ctx)
    if err != nil {
        return nil, err
    }
    var m model.TodoModel
    if err := resp.JSON(&m); err != nil {
        return nil, err
    }
    return model.ToDomainTodo(&m), nil
}

func (r *todoRepository) Delete(ctx context.Context, id string) error {
    _, err := r.supabase.DB.
        From("todos").
        Delete().
        Eq("id", id).
        Execute(ctx)
    return err
}
