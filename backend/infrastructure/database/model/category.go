package model

import (
	"time"

	"github.com/ariangn/todo-fullstack/backend/domain/entity"
)

type CategoryModel struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Color       string    `json:"color"`
	Description *string   `json:"description"`
	UserID      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToDomainCategory(m *CategoryModel) *entity.Category {
	return &entity.Category{
		ID:          m.ID,
		Name:        m.Name,
		Color:       m.Color,
		Description: m.Description,
		UserID:      m.UserID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func FromDomainCategory(c *entity.Category) *CategoryModel {
	return &CategoryModel{
		ID:          c.ID,
		Name:        c.Name,
		Color:       c.Color,
		Description: c.Description,
		UserID:      c.UserID,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
