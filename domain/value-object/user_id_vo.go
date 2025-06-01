package value-object

import (
    "errors"
    "github.com/google/uuid"
)

type UserIDVO string

func NewUserIDVO(v string) (UserIDVO, error) {
    _, err := uuid.Parse(v)
    if err != nil {
        return "", errors.New("invalid UUID format for UserID")
    }
    return UserIDVO(v), nil
}

func (u UserIDVO) String() string {
    return string(u)
}
