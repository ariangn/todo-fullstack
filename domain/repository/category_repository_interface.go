package repository

import (
    "context"
    "github.com/ariangn/todo-go/domain/entity"
)

type CategoryRepository interface {
    Create(ctx context.Context, c *entity.Category) (*entity.Category, error)
    FindByID(ctx context.Context, id string) (*entity.Category, error)
    FindAllByUser(ctx context.Context, userID string) ([]*entity.Category, error)
    Update(ctx context.Context, c *entity.Category) (*entity.Category, error)
    Delete(ctx context.Context, id string) error
}
