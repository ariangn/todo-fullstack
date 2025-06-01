package value-object

import (
    "errors"
    "regexp"
)

var emailRegex = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

type EmailVO string

func NewEmailVO(v string) (EmailVO, error) {
    if !emailRegex.MatchString(v) {
        return "", errors.New("invalid email format")
    }
    return EmailVO(v), nil
}

func (e EmailVO) String() string {
    return string(e)
}
