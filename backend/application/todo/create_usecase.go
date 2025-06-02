package todo

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ariangn/todo-fullstack/backend/domain/entity"
	"github.com/ariangn/todo-fullstack/backend/domain/repository"
	"github.com/ariangn/todo-fullstack/backend/domain/valueobject"
)

type CreateUseCase interface {
	Execute(ctx context.Context, userID, title string, body *string, dueDate *valueobject.DueDateVO, categoryID *string, tagIDs []string) (*entity.Todo, error)
}

type createUseCase struct {
	todoRepo     repository.TodoRepository
	categoryRepo repository.CategoryRepository
	tagRepo      repository.TagRepository
}

func NewCreateUseCase(
	todoRepo repository.TodoRepository,
	categoryRepo repository.CategoryRepository,
	tagRepo repository.TagRepository,
) CreateUseCase {
	return &createUseCase{todoRepo, categoryRepo, tagRepo}
}

func (uc *createUseCase) Execute(
	ctx context.Context,
	userID, title string,
	body *string,
	dueDate *valueobject.DueDateVO,
	categoryID *string,
	tagIDs []string,
) (*entity.Todo, error) {
	// validate TitleVO
	titleVO, err := valueobject.NewTitleVO(title)
	if err != nil {
		return nil, err
	}

	// validate BodyVO if provided
	var bodyStr *string
	if body != nil {
		bodyVO, err := valueobject.NewBodyVO(*body)
		if err != nil {
			return nil, err
		}
		s := bodyVO.String()
		bodyStr = &s
	}

	// validate DueDateVO if provided
	var dd *valueobject.DueDateVO
	if dueDate != nil {
		dd = dueDate
	}

	// check category exists if categoryID != nil
	if categoryID != nil {
		_, err := uc.categoryRepo.FindByID(ctx, *categoryID)
		if err != nil {
			return nil, err
		}
	}

	// check each tag exists if tagIDs provided
	for _, tID := range tagIDs {
		_, err := uc.tagRepo.FindByID(ctx, tID)
		if err != nil {
			return nil, err
		}
	}

	// build domain entity
	todoEntity, err := entity.NewTodo(
		uuid.NewString(),
		titleVO.String(),
		bodyStr,
		entity.StatusTodo,
		func() *time.Time {
			if dd != nil {
				t := dd.Time()
				return &t
			}
			return nil
		}(),
		userID,
		categoryID,
		tagIDs,
	)
	if err != nil {
		return nil, err
	}

	return uc.todoRepo.Create(ctx, todoEntity)
}
