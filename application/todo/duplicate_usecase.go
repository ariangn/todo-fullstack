package todo

import (
    "context"
    "time"

    "github.com/google/uuid"

    "github.com/ariangn/todo-go/domain/entity"
    "github.com/ariangn/todo-go/domain/repository"
)

type DuplicateUseCase interface {
    Execute(ctx context.Context, id string) (*entity.Todo, error)
}

type duplicateUseCase struct {
    todoRepo repository.TodoRepository
}

func NewDuplicateUseCase(todoRepo repository.TodoRepository) DuplicateUseCase {
    return &duplicateUseCase{todoRepo}
}

func (uc *duplicateUseCase) Execute(ctx context.Context, id string) (*entity.Todo, error) {
    original, err := uc.todoRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    if original == nil {
        return nil, ErrTodoNotFound
    }
    // create a new Todo entity that’s a copy (except new ID, CreatedAt, UpdatedAt, CompletedAt=nil)
    newTodo := &entity.Todo{
        ID:          uuid.NewString(),
        Title:       original.Title + " (Copy)",
        Body:        original.Body,
        Status:      entity.StatusTodo,
        DueDate:     original.DueDate,
        CompletedAt: nil,
        UserID:      original.UserID,
        CategoryID:  original.CategoryID,
        TagIDs:      original.TagIDs,
        CreatedAt:   time.Now().UTC(),
        UpdatedAt:   time.Now().UTC(),
    }
    return uc.todoRepo.Create(ctx, newTodo)
}
