package response

import "time"

type UserResponseDTO struct {
    ID        string     `json:"id"`
    Email     string     `json:"email"`
    Name      *string    `json:"name,omitempty"`
    AvatarURL *string    `json:"avatarUrl,omitempty"`
    Timezone  string     `json:"timezone"`
    Token     *string    `json:"token,omitempty"`
    CreatedAt time.Time  `json:"createdAt"`
    UpdatedAt time.Time  `json:"updatedAt"`
}
