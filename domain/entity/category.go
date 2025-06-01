package entity

import (
    "errors"
    "time"
)

type Category struct {
    ID          string
    Name        string
    Color       string
    Description *string
    UserID      string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// NewCategory enforces: Name non-empty, UserID non-empty
func NewCategory(
    id, name, color, userID string,
    description *string,
) (*Category, error) {
    if name == "" {
        return nil, errors.New("category name cannot be empty")
    }
    if userID == "" {
        return nil, errors.New("userID cannot be empty")
    }
    return &Category{
        ID:          id,
        Name:        name,
        Color:       color,
        Description: description,
        UserID:      userID,
        CreatedAt:   time.Now().UTC(),
        UpdatedAt:   time.Now().UTC(),
    }, nil
}
