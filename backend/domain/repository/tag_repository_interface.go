package repository

import (
	"context"

	"github.com/ariangn/todo-fullstack/backend/domain/entity"
)

type TagRepository interface {
	Create(ctx context.Context, t *entity.Tag) (*entity.Tag, error)
	FindByID(ctx context.Context, id string) (*entity.Tag, error)
	FindAllByUser(ctx context.Context, userID string) ([]*entity.Tag, error)
	FindByName(ctx context.Context, userID string, name string) (*entity.Tag, error)
	Update(ctx context.Context, t *entity.Tag) (*entity.Tag, error)
	Delete(ctx context.Context, id string) error
}
