package request

import "time"

type CreateTodoDTO struct {
    Title      string     `json:"title"`
    Body       *string    `json:"body,omitempty"`
    DueDate    *time.Time `json:"dueDate,omitempty"`
    CategoryID *string    `json:"categoryId,omitempty"`
    TagIDs     []string   `json:"tagIds,omitempty"`
}
