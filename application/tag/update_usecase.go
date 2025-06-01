package tag

import (
    "context"
    "errors"

    "github.com/ariangn/todo-go/domain/entity"
    "github.com/ariangn/todo-go/domain/repository"
)

var ErrTagNotFound = errors.New("tag not found")

type UpdateUseCase interface {
    Execute(ctx context.Context, t *entity.Tag) (*entity.Tag, error)
}

type updateUseCase struct {
    tagRepo repository.TagRepository
}

func NewUpdateUseCase(tagRepo repository.TagRepository) UpdateUseCase {
    return &updateUseCase{tagRepo}
}

func (uc *updateUseCase) Execute(ctx context.Context, t *entity.Tag) (*entity.Tag, error) {
    existing, err := uc.tagRepo.FindByID(ctx, t.ID)
    if err != nil {
        return nil, err
    }
    if existing == nil {
        return nil, ErrTagNotFound
    }
    t.UpdatedAt = time.Now().UTC()
    return uc.tagRepo.Update(ctx, t)
}
