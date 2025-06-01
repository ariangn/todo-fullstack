package value-object

import "errors"

type TitleVO string

func NewTitleVO(v string) (TitleVO, error) {
    if len(v) == 0 {
        return "", errors.New("title cannot be empty")
    }
    if len(v) > 200 {
        return "", errors.New("title cannot exceed 200 characters")
    }
    return TitleVO(v), nil
}

func (t TitleVO) String() string {
    return string(t)
}
