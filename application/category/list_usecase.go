package category

import (
    "context"

    "github.com/ariangn/todo-go/domain/entity"
    "github.com/ariangn/todo-go/domain/repository"
)

type ListUseCase interface {
    Execute(ctx context.Context, userID string) ([]*entity.Category, error)
}

type listUseCase struct {
    categoryRepo repository.CategoryRepository
}

func NewListUseCase(categoryRepo repository.CategoryRepository) ListUseCase {
    return &listUseCase{categoryRepo}
}

func (uc *listUseCase) Execute(ctx context.Context, userID string) ([]*entity.Category, error) {
    return uc.categoryRepo.FindAllByUser(ctx, userID)
}
