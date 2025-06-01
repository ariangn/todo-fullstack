package request

// add fields here
// does not implement filtering logic
type TodoFiltersDTO struct {
    Status     *string `form:"status,omitempty"`
    CategoryID *string `form:"categoryId,omitempty"`
    TagID      *string `form:"tagId,omitempty"`
    DueDateFrom *string `form:"dueDateFrom,omitempty"`
    DueDateTo   *string `form:"dueDateTo,omitempty"`
}
