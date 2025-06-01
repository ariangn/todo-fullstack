package repository

import (
    "context"
    "github.com/ariangn/todo-go/domain/entity"
)

type UserRepository interface {
    Create(ctx context.Context, u *entity.User) (*entity.User, error)
    FindByEmail(ctx context.Context, email string) (*entity.User, error)
}
