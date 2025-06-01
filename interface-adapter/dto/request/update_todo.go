package request

import "time"

type UpdateTodoDTO struct {
    Title      *string    `json:"title,omitempty"`
    Body       *string    `json:"body,omitempty"`
    DueDate    *time.Time `json:"dueDate,omitempty"`
    Status     *string    `json:"status,omitempty"`
    CategoryID *string    `json:"categoryId,omitempty"`
    TagIDs     *[]string  `json:"tagIds,omitempty"`
}
