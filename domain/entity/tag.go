package entity

import (
    "errors"
    "time"
)

type Tag struct {
    ID        string
    Name      string
    UserID    string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// NewTag enforces: Name non-empty, UserID non-empty
func NewTag(id, name, userID string) (*Tag, error) {
    if name == "" {
        return nil, errors.New("tag name cannot be empty")
    }
    if userID == "" {
        return nil, errors.New("userID cannot be empty")
    }
    return &Tag{
        ID:        id,
        Name:      name,
        UserID:    userID,
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
    }, nil
}
