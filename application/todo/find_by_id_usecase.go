// application/todo/find_by_id_use_case.go
package todo

import (
	"context"

	"github.com/ariangn/todo-go/domain/entity"
	"github.com/ariangn/todo-go/domain/repository"
)

type FindByIDUseCase interface {
	Execute(ctx context.Context, id string) (*entity.Todo, error)
}

type findByIDUseCase struct {
	todoRepo repository.TodoRepository
}

func NewFindByIDUseCase(todoRepo repository.TodoRepository) FindByIDUseCase {
	return &findByIDUseCase{todoRepo}
}

func (uc *findByIDUseCase) Execute(ctx context.Context, id string) (*entity.Todo, error) {
	return uc.todoRepo.FindByID(ctx, id)
}
