package request

type CreateCategoryDTO struct {
    Name        string  `json:"name"`
    Color       string  `json:"color"`
    Description *string `json:"description,omitempty"`
}
