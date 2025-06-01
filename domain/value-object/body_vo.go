package value-object

import "errors"

type BodyVO string

func NewBodyVO(v string) (BodyVO, error) {
    if len(v) > 1000 {
        return "", errors.New("body cannot exceed 1000 characters")
    }
    return BodyVO(v), nil
}

func (b BodyVO) String() string {
    return string(b)
}
