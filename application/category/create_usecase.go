package category

import (
    "context"

    "github.com/google/uuid"

    "github.com/ariangn/todo-go/domain/entity"
    "github.com/ariangn/todo-go/domain/repository"
    "github.com/ariangn/todo-go/domain/value-object"
)

type CreateUseCase interface {
    Execute(ctx context.Context, userID, name, color string, description *string) (*entity.Category, error)
}

type createUseCase struct {
    categoryRepo repository.CategoryRepository
}

func NewCreateUseCase(categoryRepo repository.CategoryRepository) CreateUseCase {
    return &createUseCase{categoryRepo}
}

func (uc *createUseCase) Execute(ctx context.Context, userID, name, color string, description *string) (*entity.Category, error) {
    // Validate NameVO
    nameVO, err := value-object.NewTitleVO(name) // reuse TitleVO for non-empty check
    if err != nil {
        return nil, err
    }
    // Color is pre-validated by front-end

    userIDVO, err := value-object.NewUserIDVO(userID)
    if err != nil {
        return nil, err
    }

    catEntity, err := entity.NewCategory(
        uuid.NewString(),
        nameVO.String(),
        color,
        userIDVO.String(),
        description,
    )
    if err != nil {
        return nil, err
    }
    return uc.categoryRepo.Create(ctx, catEntity)
}
