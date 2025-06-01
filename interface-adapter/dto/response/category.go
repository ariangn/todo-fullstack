package response

import "time"

type CategoryResponseDTO struct {
    ID          string     `json:"id"`
    Name        string     `json:"name"`
    Color       string     `json:"color"`
    Description *string    `json:"description,omitempty"`
    UserID      string     `json:"userId"`
    CreatedAt   time.Time  `json:"createdAt"`
    UpdatedAt   time.Time  `json:"updatedAt"`
}
