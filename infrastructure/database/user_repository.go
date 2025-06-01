package database

import (
    "context"
    "errors"

    "github.com/ariangn/todo-go/domain/entity"
    "github.com/ariangn/todo-go/domain/repository"
    "github.com/ariangn/todo-go/infrastructure/database/model"
    "github.com/google/uuid"
)

type userRepository struct {
    supabase *SupabaseClient
}

func NewUserRepository(supabase *SupabaseClient) repository.UserRepository {
    return &userRepository{supabase}
}

func (r *userRepository) Create(ctx context.Context, u *entity.User) (*entity.User, error) {
    u.ID = uuid.NewString()
    row := map[string]interface{}{
        "id":         u.ID,
        "email":      u.Email,
        "password":   u.Password,
        "name":       u.Name,
        "avatar_url": u.AvatarURL,
        "timezone":   u.Timezone,
    }

    resp, err := r.supabase.DB.
        From("users").
        Insert(row).
        Single().
        Execute(ctx)
    if err != nil {
        return nil, err
    }
    var m model.UserModel
    if err := resp.JSON(&m); err != nil {
        return nil, err
    }
    return model.ToDomainUser(&m), nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
    resp, err := r.supabase.DB.
        From("users").
        Select("*").
        Eq("email", email).
        Single().
        Execute(ctx)
    if err != nil {
        return nil, err
    }
    var m model.UserModel
    if err := resp.JSON(&m); err != nil {
        return nil, err
    }
    return model.ToDomainUser(&m), nil
}
