package tag

import (
    "context"

    "github.com/google/uuid"

    "github.com/ariangn/todo-go/domain/entity"
    "github.com/ariangn/todo-go/domain/repository"
    "github.com/ariangn/todo-go/domain/valueobject"
)

type CreateUseCase interface {
    Execute(ctx context.Context, userID, name string) (*entity.Tag, error)
}

type createUseCase struct {
    tagRepo repository.TagRepository
}

func NewCreateUseCase(tagRepo repository.TagRepository) CreateUseCase {
    return &createUseCase{tagRepo}
}

func (uc *createUseCase) Execute(ctx context.Context, userID, name string) (*entity.Tag, error) {
    nameVO, err := valueobject.NewTitleVO(name) // reuse TitleVO for name non-empty
    if err != nil {
        return nil, err
    }
    userIDVO, err := valueobject.NewUserIDVO(userID)
    if err != nil {
        return nil, err
    }
    tagEntity, err := entity.NewTag(
        uuid.NewString(),
        nameVO.String(),
        userIDVO.String(),
    )
    if err != nil {
        return nil, err
    }
    return uc.tagRepo.Create(ctx, tagEntity)
}
