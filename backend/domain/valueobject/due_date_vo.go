package valueobject

import (
	"time"
)

type DueDateVO time.Time

func NewDueDateVO(v time.Time) (DueDateVO, error) {
	return DueDateVO(v), nil
}

func (d DueDateVO) Time() time.Time {
	return time.Time(d)
}

func (d DueDateVO) String() string {
	return time.Time(d).Format(time.RFC3339)
}
