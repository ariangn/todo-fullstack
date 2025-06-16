package tag

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/ariangn/todo-fullstack/backend/domain/entity"
	"github.com/ariangn/todo-fullstack/backend/domain/repository"
	"github.com/ariangn/todo-fullstack/backend/domain/valueobject"
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
	fmt.Println("UC START: Creating tag with name =", name)

	nameVO, err := valueobject.NewTitleVO(name)
	if err != nil {
		fmt.Println("UC ERROR: TitleVO failed:", err)
		return nil, err
	}
	fmt.Println("UC OK: TitleVO =", nameVO)

	userIDVO, err := valueobject.NewUserIDVO(userID)
	if err != nil {
		fmt.Println("UC ERROR: UserIDVO failed:", err)
		return nil, err
	}
	fmt.Println("UC OK: UserIDVO =", userIDVO)

	tagEntity, err := entity.NewTag(
		uuid.NewString(),
		nameVO.String(),
		userIDVO.String(),
	)
	if err != nil {
		fmt.Println("UC ERROR: NewTag failed:", err)
		return nil, err
	}
	fmt.Println("UC OK: Created tagEntity =", tagEntity)

	created, err := uc.tagRepo.Create(ctx, tagEntity)
	if err != nil {
		fmt.Println("UC ERROR: tagRepo.Create failed:", err)
		return nil, err
	}
	fmt.Println("UC SUCCESS: Tag created =", created)

	return created, nil
}
