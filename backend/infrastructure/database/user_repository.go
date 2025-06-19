package database

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"

	"github.com/ariangn/todo-fullstack/backend/domain/entity"
	"github.com/ariangn/todo-fullstack/backend/domain/repository"
	"github.com/ariangn/todo-fullstack/backend/infrastructure/database/model"
)

type userRepository struct {
	supabase *SupabaseClient
}

func NewUserRepository(db *SupabaseClient) repository.UserRepository {
	return &userRepository{supabase: db}
}

func (r *userRepository) Create(ctx context.Context, u *entity.User) (*entity.User, error) {
	// 1) assign a new ID
	u.ID = uuid.NewString()

	toInsert := map[string]interface{}{
		"id":         u.ID,
		"email":      u.Email,
		"password":   u.Password,
		"name":       u.Name,
		"avatar_url": u.AvatarURL,
		"timezone":   u.Timezone,
	}

	// 2) perform the insert, ignore the raw result
	if _, _, err := r.supabase.DB.
		From("users").
		Insert(toInsert, false, "", "*", "").
		Execute(); err != nil {
		// handle duplicate‐email more cleanly
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, errors.New("email is already taken")
		}
		return nil, err
	}

	// 3) now fetch the user we just created
	return r.FindByID(ctx, u.ID)
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	// SELECT * FROM users WHERE email = '<email>' LIMIT 1
	builder := r.supabase.DB.
		From("users").
		Select("*", "", false).
		Eq("email", email).
		Single()

	raw, _, err := builder.Execute()
	if err != nil {
		return nil, err
	}

	var m model.UserModel
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return model.ToDomainUser(&m), nil
}

// FindByID fetches a user by its ID. Returns an error if not found.
func (r *userRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	// SELECT * FROM users WHERE id = '<id>' LIMIT 1
	builder := r.supabase.DB.
		From("users").
		Select("*", "", false).
		Eq("id", id).
		Single()

	raw, _, err := builder.Execute()
	if err != nil {
		return nil, err
	}

	// If raw is empty or unmarshal fails, treat as “not found”
	var m model.UserModel
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}

	// Convert to domain and return
	userEntity := model.ToDomainUser(&m)
	if userEntity == nil || userEntity.ID == "" {
		return nil, errors.New("user not found")
	}
	return userEntity, nil
}
