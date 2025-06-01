package valueobject

import "time"

type CompletedAtVO time.Time

func NewCompletedAtVO(v time.Time) CompletedAtVO {
    return CompletedAtVO(v)
}

func (c CompletedAtVO) Time() time.Time {
    return time.Time(c)
}

func (c CompletedAtVO) String() string {
    return time.Time(c).Format(time.RFC3339)
}
