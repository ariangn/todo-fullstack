package tag

import "context"

import "github.com/ariangn/todo-go/domain/repository"

type DeleteUseCase interface {
    Execute(ctx context.Context, id string) error
}

type deleteUseCase struct {
    tagRepo repository.TagRepository
}

func NewDeleteUseCase(tagRepo repository.TagRepository) DeleteUseCase {
    return &deleteUseCase{tagRepo}
}

func (uc *deleteUseCase) Execute(ctx context.Context, id string) error {
    return uc.tagRepo.Delete(ctx, id)
}
