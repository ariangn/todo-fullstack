package entity

import (
    "time"
)

type User struct {
    ID        string
    Email     string
    Password  string
    Name      *string
    AvatarURL *string
    Timezone  string
    CreatedAt time.Time
    UpdatedAt time.Time
}

func NewUser(
    id, email, password string,
    name, avatarURL *string,
    timezone string,
) (*User, error) {
    if email == "" {
        return nil, ErrEmailEmpty
    }
    if timezone == "" {
        return nil, ErrTimezoneEmpty
    }
    return &User{
        ID:        id,
        Email:     email,
        Password:  password,
        Name:      name,
        AvatarURL: avatarURL,
        Timezone:  timezone,
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
    }, nil
}

var (
    ErrEmailEmpty    = errors.New("email cannot be empty")
    ErrTimezoneEmpty = errors.New("timezone cannot be empty")
)
