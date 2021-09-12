package discord

import "time"

type OptionalBool *bool

func NewBool(b bool) *OptionalBool {
	ob := OptionalBool(&b)
	return &ob
}

func NewNullBool() *OptionalBool {
	ob := OptionalBool(nil)
	return &ob
}

type OptionalString *string

func NewString(str string) *OptionalString {
	ostr := OptionalString(&str)
	return &ostr
}

func NewNullString() *OptionalString {
	ostr := OptionalString(nil)
	return &ostr
}

type OptionalTime *Time

func NewTime(time time.Time) *OptionalTime {
	ot := OptionalTime(&Time{Time: time})
	return &ot
}

func NewNullTime() *OptionalTime {
	ot := OptionalTime(nil)
	return &ot
}
