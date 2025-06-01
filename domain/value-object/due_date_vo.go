package valueobject

import (
    "errors"
    "time"
)

type DueDateVO time.Time

func NewDueDateVO(v time.Time) (DueDateVO, error) {
    if v.Before(time.Now().UTC()) {
        return DueDateVO(time.Time{}), errors.New("due date cannot be in the past")
    }
    return DueDateVO(v), nil
}

func (d DueDateVO) Time() time.Time {
    return time.Time(d)
}

func (d DueDateVO) String() string {
    return time.Time(d).Format(time.RFC3339)
}
