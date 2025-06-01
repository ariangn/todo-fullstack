package category

import (
    "context"
    "errors"

    "github.com/ariangn/todo-go/domain/entity"
    "github.com/ariangn/todo-go/domain/repository"
)

var ErrCategoryNotFound = errors.New("category not found")

type UpdateUseCase interface {
    Execute(ctx context.Context, c *entity.Category) (*entity.Category, error)
}

type updateUseCase struct {
    categoryRepo repository.CategoryRepository
}

func NewUpdateUseCase(categoryRepo repository.CategoryRepository) UpdateUseCase {
    return &updateUseCase{categoryRepo}
}

func (uc *updateUseCase) Execute(ctx context.Context, c *entity.Category) (*entity.Category, error) {
    existing, err := uc.categoryRepo.FindByID(ctx, c.ID)
    if err != nil {
        return nil, err
    }
    if existing == nil {
        return nil, ErrCategoryNotFound
    }
    c.UpdatedAt = time.Now().UTC()
    return uc.categoryRepo.Update(ctx, c)
}
