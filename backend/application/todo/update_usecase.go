package todo

import (
	"context"
	"errors"
	"time"

	"github.com/ariangn/todo-fullstack/backend/domain/entity"
	"github.com/ariangn/todo-fullstack/backend/domain/repository"
)

var ErrTodoNotFound = errors.New("todo not found")

type UpdateUseCase interface {
	Execute(ctx context.Context, t *entity.Todo) (*entity.Todo, error)
}

type updateUseCase struct {
	todoRepo repository.TodoRepository
}

func NewUpdateUseCase(todoRepo repository.TodoRepository) UpdateUseCase {
	return &updateUseCase{todoRepo}
}

func (uc *updateUseCase) Execute(ctx context.Context, t *entity.Todo) (*entity.Todo, error) {
	existing, err := uc.todoRepo.FindByID(ctx, t.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrTodoNotFound
	}
	t.UpdatedAt = time.Now().UTC()
	return uc.todoRepo.Update(ctx, t)
}
