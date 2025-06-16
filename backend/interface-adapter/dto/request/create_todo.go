package request

import "time"

type CreateTodoDTO struct {
	Title      string     `json:"title"`
	Body       *string    `json:"body"`
	DueDate    *time.Time `json:"dueDate"`
	Status     string     `json:"status"`
	CategoryID *string    `json:"categoryId"`
	TagIDs     []string   `json:"tagIds"`
}
