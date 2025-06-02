package repository

import (
	"context"

	"github.com/ariangn/todo-fullstack/backend/domain/entity"
)

type TodoRepository interface {
	Create(ctx context.Context, t *entity.Todo) (*entity.Todo, error)
	FindByID(ctx context.Context, id string) (*entity.Todo, error)
	FindAllByUser(ctx context.Context, userID string) ([]*entity.Todo, error)
	Update(ctx context.Context, t *entity.Todo) (*entity.Todo, error)
	Delete(ctx context.Context, id string) error
}
