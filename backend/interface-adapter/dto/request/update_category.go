package request

type UpdateCategoryDTO struct {
    Name        *string `json:"name,omitempty"`
    Color       *string `json:"color,omitempty"`
    Description *string `json:"description,omitempty"`
}
