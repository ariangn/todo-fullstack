package response

import "time"

type TagResponseDTO struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    UserID    string    `json:"userId"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}
