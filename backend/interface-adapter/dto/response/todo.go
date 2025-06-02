package response

import "time"

type TodoResponseDTO struct {
    ID          string     `json:"id"`
    Title       string     `json:"title"`
    Body        *string    `json:"body,omitempty"`
    Status      string     `json:"status"`
    DueDate     *time.Time `json:"dueDate,omitempty"`
    CompletedAt *time.Time `json:"completedAt,omitempty"`
    UserID      string     `json:"userId"`
    CategoryID  *string    `json:"categoryId,omitempty"`
    TagIDs      []string   `json:"tagIds,omitempty"`
    CreatedAt   time.Time  `json:"createdAt"`
    UpdatedAt   time.Time  `json:"updatedAt"`
}
