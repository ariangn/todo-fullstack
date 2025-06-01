package todo

import (
    "context"

    "github.com/ariangn/todo-go/domain/repository"
)

type DeleteUseCase interface {
    Execute(ctx context.Context, id string) error
}

type deleteUseCase struct {
    todoRepo repository.TodoRepository
}

func NewDeleteUseCase(todoRepo repository.TodoRepository) DeleteUseCase {
    return &deleteUseCase{todoRepo}
}

func (uc *deleteUseCase) Execute(ctx context.Context, id string) error {
    return uc.todoRepo.Delete(ctx, id)
}
