package todo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/ariangn/todo-fullstack/backend/domain/entity"
	"github.com/ariangn/todo-fullstack/backend/domain/repository"
	"github.com/ariangn/todo-fullstack/backend/domain/valueobject"
)

type CreateUseCase interface {
	Execute(ctx context.Context, userID, title string, body *string, status entity.Status, dueDate *valueobject.DueDateVO, categoryID *string, tagIDs []string) (*entity.Todo, error)
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
	status entity.Status,
	dueDate *valueobject.DueDateVO,
	categoryID *string,
	tagIDs []string,
) (*entity.Todo, error) {
	fmt.Println("UC START: Execute called")

	// validate TitleVO
	titleVO, err := valueobject.NewTitleVO(title)
	if err != nil {
		fmt.Println("UC ERROR: NewTitleVO failed:", err)
		return nil, err
	}
	fmt.Println("UC OK: TitleVO =", titleVO)

	// validate BodyVO if provided
	var bodyStr *string
	if body != nil {
		bodyVO, err := valueobject.NewBodyVO(*body)
		if err != nil {
			fmt.Println("UC ERROR: NewBodyVO failed:", err)
			return nil, err
		}
		s := bodyVO.String()
		bodyStr = &s
		fmt.Println("UC OK: BodyVO =", *bodyStr)
	}

	// validate DueDateVO if provided
	var dd *valueobject.DueDateVO
	if dueDate != nil {
		dd = dueDate
		fmt.Println("UC OK: DueDate =", dd)
	}

	// check category exists if categoryID != nil
	if categoryID != nil {
		fmt.Println("UC CHECK: Looking up category:", *categoryID)
		cat, err := uc.categoryRepo.FindByID(ctx, *categoryID)
		if err != nil {
			fmt.Println("UC ERROR: Category not found:", err)
			return nil, fmt.Errorf("invalid category ID: %w", err)
		}
		fmt.Printf("UC OK: Category exists: %+v\n", cat)
	}

	fmt.Println("UC OK: TagIDs passed =", tagIDs)

	// build domain entity
	todoEntity, err := entity.NewTodo(
		uuid.NewString(),
		titleVO.String(),
		bodyStr,
		status,
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
		fmt.Println("UC ERROR: entity.NewTodo failed:", err)
		return nil, err
	}
	fmt.Println("UC OK: Created todoEntity =", todoEntity)

	fmt.Println("UC FINAL: About to call todoRepo.Create with tagIDs =", tagIDs)

	created, err := uc.todoRepo.Create(ctx, todoEntity)
	if err != nil {
		fmt.Println("UC ERROR: todoRepo.Create failed:", err)
		return nil, err
	}
	fmt.Println("UC SUCCESS: Todo created =", created)

	return created, nil
}
