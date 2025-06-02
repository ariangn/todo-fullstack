package user

import (
	"context"
	"errors"

	"github.com/ariangn/todo-fullstack/backend/domain/entity"
	"github.com/ariangn/todo-fullstack/backend/domain/repository"
	"github.com/ariangn/todo-fullstack/backend/domain/valueobject"
	"github.com/google/uuid"
)

type RegisterUseCase interface {
	Execute(ctx context.Context, email, password string, name *string, timezone string, avatarURL *string) (*entity.User, error)
}

type registerUseCase struct {
	userRepo repository.UserRepository
}

func NewRegisterUseCase(userRepo repository.UserRepository) RegisterUseCase {
	return &registerUseCase{userRepo}
}

func (uc *registerUseCase) Execute(ctx context.Context, email, password string, name *string, timezone string, avatarURL *string) (*entity.User, error) {
	// validate email via EmailVO
	emailVO, err := valueobject.NewEmailVO(email)
	if err != nil {
		return nil, err
	}
	// validate & hash password via PasswordVO
	pwdVO, err := valueobject.NewPasswordVO(password)
	if err != nil {
		return nil, err
	}
	hashedPwd := pwdVO.Hash()

	// use ValueObject for timezone (just non-empty check)
	if timezone == "" {
		return nil, ErrTimezoneMissing
	}

	// construct domain entity
	userEntity, err := entity.NewUser(
		uuid.NewString(),
		emailVO.String(),
		hashedPwd,
		name,
		avatarURL,
		timezone,
	)
	if err != nil {
		return nil, err
	}
	return uc.userRepo.Create(ctx, userEntity)
}

var ErrTimezoneMissing = errors.New("timezone is required")
