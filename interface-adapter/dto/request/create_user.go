package request

type CreateUserDTO struct {
    Email     string  `json:"email"`
    Password  string  `json:"password"`
    Name      *string `json:"name,omitempty"`
    Timezone  string  `json:"timezone"`
    AvatarURL *string `json:"avatarUrl,omitempty"`
}
