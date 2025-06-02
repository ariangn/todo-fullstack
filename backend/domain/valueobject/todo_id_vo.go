package valueobject

import (
    "errors"
    "github.com/google/uuid"
)

type TodoIDVO string

func NewTodoIDVO(v string) (TodoIDVO, error) {
    _, err := uuid.Parse(v)
    if err != nil {
        return "", errors.New("invalid UUID format for TodoID")
    }
    return TodoIDVO(v), nil
}

func (t TodoIDVO) String() string {
    return string(t)
}
