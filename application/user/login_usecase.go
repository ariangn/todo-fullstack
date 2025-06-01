package user

import (
    "context"
    "errors"
    "time"

    "github.com/ariangn/todo-go/domain/repository"
    "github.com/ariangn/todo-go/domain/valueobject"
    "github.com/ariangn/todo-go/infrastructure/auth"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type LoginUseCase interface {
    Execute(ctx context.Context, email, password string) (string /* JWT token */, error)
}

type loginUseCase struct {
    userRepo   repository.UserRepository
    authClient auth.AuthClientInterface
}

func NewLoginUseCase(userRepo repository.UserRepository, authClient auth.AuthClientInterface) LoginUseCase {
    return &loginUseCase{userRepo, authClient}
}

func (uc *loginUseCase) Execute(ctx context.Context, email, password string) (string, error) {
    // lookup user by email
    existing, err := uc.userRepo.FindByEmail(ctx, email)
    if err != nil {
        return "", err
    }
    if existing == nil {
        return "", ErrInvalidCredentials
    }
    // verify password
    pwdVO := valueobject.NewPasswordVOWithHash(existing.Password)
    if !pwdVO.Verify(password) {
        return "", ErrInvalidCredentials
    }
    // gnerate JWT (24h TTL)
    token, err := uc.authClient.GenerateToken(existing.ID, 24*time.Hour)
    if err != nil {
        return "", err
    }
    return token, nil
}
