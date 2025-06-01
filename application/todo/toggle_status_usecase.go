package todo

import (
	"context"
	"time"

	"github.com/ariangn/todo-go/domain/entity"
	"github.com/ariangn/todo-go/domain/repository"
)

type ToggleStatusUseCase interface {
	Execute(ctx context.Context, id string, newStatus entity.Status) (*entity.Todo, error)
}

type toggleStatusUseCase struct {
	todoRepo repository.TodoRepository
}

func NewToggleStatusUseCase(todoRepo repository.TodoRepository) ToggleStatusUseCase {
	return &toggleStatusUseCase{todoRepo}
}

func (uc *toggleStatusUseCase) Execute(ctx context.Context, id string, newStatus entity.Status) (*entity.Todo, error) {
	t, err := uc.todoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, ErrTodoNotFound
	}
	t.Status = newStatus
	t.UpdatedAt = time.Now().UTC()
	if newStatus == entity.StatusCompleted {
		now := time.Now().UTC()
		t.CompletedAt = &now
	} else {
		t.CompletedAt = nil
	}
	return uc.todoRepo.Update(ctx, t)
}
