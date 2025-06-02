package todo

import (
	"context"

	"github.com/ariangn/todo-fullstack/backend/domain/entity"
	"github.com/ariangn/todo-fullstack/backend/domain/repository"
)

type ListUseCase interface {
	Execute(ctx context.Context, userID string) ([]*entity.Todo, error)
}

type listUseCase struct {
	todoRepo repository.TodoRepository
}

func NewListUseCase(todoRepo repository.TodoRepository) ListUseCase {
	return &listUseCase{todoRepo}
}

func (uc *listUseCase) Execute(ctx context.Context, userID string) ([]*entity.Todo, error) {
	return uc.todoRepo.FindAllByUser(ctx, userID)
}
