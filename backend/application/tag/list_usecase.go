package tag

import (
	"context"

	"github.com/ariangn/todo-fullstack/backend/domain/entity"
	"github.com/ariangn/todo-fullstack/backend/domain/repository"
)

type ListUseCase interface {
	Execute(ctx context.Context, userID string) ([]*entity.Tag, error)
}

type listUseCase struct {
	tagRepo repository.TagRepository
}

func NewListUseCase(tagRepo repository.TagRepository) ListUseCase {
	return &listUseCase{tagRepo}
}

func (uc *listUseCase) Execute(ctx context.Context, userID string) ([]*entity.Tag, error) {
	return uc.tagRepo.FindAllByUser(ctx, userID)
}
