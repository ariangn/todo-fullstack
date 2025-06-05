package user

import (
	"context"
	"errors"

	"fmt"

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
		fmt.Println("Email validation failed:", err)
		return nil, err
	}
	// validate & hash password via PasswordVO
	pwdVO, err := valueobject.NewPasswordVO(password)
	if err != nil {
		fmt.Println("Password validation failed:", err)
		return nil, err
	}
	hashedPwd := pwdVO.Hash()

	// use ValueObject for timezone (just non-empty check)
	if timezone == "" {
		fmt.Println("Missing timezone")
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
		fmt.Println("User creation failed:", err)
		return nil, err
	}
	res, err := uc.userRepo.Create(ctx, userEntity)
	if err != nil {
		fmt.Println("User repo creation failed:", err)
		return nil, err
	}
	return res, nil
}

var ErrTimezoneMissing = errors.New("timezone is required")
