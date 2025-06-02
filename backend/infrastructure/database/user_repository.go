package database

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/ariangn/todo-fullstack/backend/domain/entity"
	"github.com/ariangn/todo-fullstack/backend/domain/repository"
	"github.com/ariangn/todo-fullstack/backend/infrastructure/database/model"
)

type userRepository struct {
	supabase *SupabaseClient
}

func NewUserRepository(supabase *SupabaseClient) repository.UserRepository {
	return &userRepository{supabase}
}

func (r *userRepository) Create(ctx context.Context, u *entity.User) (*entity.User, error) {
	u.ID = uuid.NewString()

	// Build row to insert
	toInsert := map[string]interface{}{
		"id":         u.ID,
		"email":      u.Email,
		"password":   u.Password,
		"name":       u.Name,
		"avatar_url": u.AvatarURL,
		"timezone":   u.Timezone,
	}

	// Perform: INSERT INTO users VALUES(...) RETURNING *
	builder := r.supabase.DB.
		From("users").
		Insert(toInsert, false, "", "*", "").
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
