package user

import (
	"context"

	"github.com/ariangn/todo-fullstack/backend/domain/entity"
	"github.com/ariangn/todo-fullstack/backend/domain/repository"
)

// FindByIDUseCase defines the interface for fetching a user by their ID.
type FindByIDUseCase interface {
	Execute(ctx context.Context, id string) (*entity.User, error)
}

// findByIDUseCase is the concrete implementation.
type findByIDUseCase struct {
	repo repository.UserRepository
}

// NewFindByIDUseCase constructs a FindByIDUseCase given a UserRepository.
func NewFindByIDUseCase(repo repository.UserRepository) FindByIDUseCase {
	return &findByIDUseCase{repo: repo}
}

// Execute calls the repository’s FindByID method.
func (uc *findByIDUseCase) Execute(ctx context.Context, id string) (*entity.User, error) {
	return uc.repo.FindByID(ctx, id)
}
