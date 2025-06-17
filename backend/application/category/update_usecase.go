package category

import (
	"context"
	"errors"
	"time"

	"github.com/ariangn/todo-fullstack/backend/domain/entity"
	"github.com/ariangn/todo-fullstack/backend/domain/repository"
)

var (
	ErrCategoryNotFound  = errors.New("category not found")
	ErrCategoryForbidden = errors.New("cannot modify another user's category")
)

type UpdateUseCase interface {
	// userID: the caller's ID
	// id: the category to update
	// name, color, description: new values
	Execute(ctx context.Context, userID, id, name, color, description string) (*entity.Category, error)
}

type updateUseCase struct {
	categoryRepo repository.CategoryRepository
}

func NewUpdateUseCase(categoryRepo repository.CategoryRepository) UpdateUseCase {
	return &updateUseCase{categoryRepo}
}

func (uc *updateUseCase) Execute(
	ctx context.Context,
	userID, id, name, color, description string,
) (*entity.Category, error) {
	// 1) Fetch existing
	existing, err := uc.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrCategoryNotFound
	}

	// 2) Authorization
	if existing.UserID != userID {
		return nil, ErrCategoryForbidden
	}

	// 3) Apply updates
	existing.Name = name
	existing.Color = color
	existing.Description = &description
	existing.UpdatedAt = time.Now().UTC()

	// 4) Persist
	return uc.categoryRepo.Update(ctx, existing)
}
