package category

import "context"

import "github.com/ariangn/todo-go/domain/repository"

type DeleteUseCase interface {
    Execute(ctx context.Context, id string) error
}

type deleteUseCase struct {
    categoryRepo repository.CategoryRepository
}

func NewDeleteUseCase(categoryRepo repository.CategoryRepository) DeleteUseCase {
    return &deleteUseCase{categoryRepo}
}

func (uc *deleteUseCase) Execute(ctx context.Context, id string) error {
    return uc.categoryRepo.Delete(ctx, id)
}
