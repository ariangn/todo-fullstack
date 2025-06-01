package model

import (
    "time"

    "github.com/ariangn/todo-go/domain/entity"
)

// mirrors the JSON structure returned by PostgREST for the "users" table
type UserModel struct {
    ID        string     `json:"id"`
    Email     string     `json:"email"`
    Password  string     `json:"password"`
    Name      *string    `json:"name"`
    AvatarURL *string    `json:"avatar_url"`
    Timezone  string     `json:"timezone"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
}

func ToDomainUser(m *UserModel) *entity.User {
    return &entity.User{
        ID:        m.ID,
        Email:     m.Email,
        Password:  m.Password,
        Name:      m.Name,
        AvatarURL: m.AvatarURL,
        Timezone:  m.Timezone,
        CreatedAt: m.CreatedAt,
        UpdatedAt: m.UpdatedAt,
    }
}

func FromDomainUser(u *entity.User) *UserModel {
    return &UserModel{
        ID:        u.ID,
        Email:     u.Email,
        Password:  u.Password,
        Name:      u.Name,
        AvatarURL: u.AvatarURL,
        Timezone:  u.Timezone,
        CreatedAt: u.CreatedAt,
        UpdatedAt: u.UpdatedAt,
    }
}
