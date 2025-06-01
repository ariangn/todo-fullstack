package model

import (
    "time"

    "github.com/ariangn/todo-go/domain/entity"
)

type TagModel struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    UserID    string    `json:"user_id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func ToDomainTag(m *TagModel) *entity.Tag {
    return &entity.Tag{
        ID:        m.ID,
        Name:      m.Name,
        UserID:    m.UserID,
        CreatedAt: m.CreatedAt,
        UpdatedAt: m.UpdatedAt,
    }
}

func FromDomainTag(t *entity.Tag) *TagModel {
    return &TagModel{
        ID:        t.ID,
        Name:      t.Name,
        UserID:    t.UserID,
        CreatedAt: t.CreatedAt,
        UpdatedAt: t.UpdatedAt,
    }
}
