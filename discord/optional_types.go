package discord

import "time"

type OptionalBool *bool

func NewOptionalBool(b bool) *OptionalBool {
	ob := OptionalBool(&b)
	return &ob
}

func NewNullBool() *OptionalBool {
	ob := OptionalBool(nil)
	return &ob
}

type OptionalString *string

func NewOptionalString(str string) *OptionalString {
	ostr := OptionalString(&str)
	return &ostr
}

func NewNullString() *OptionalString {
	ostr := OptionalString(nil)
	return &ostr
}

type OptionalInt *int

func NewOptionalInt(int int) *OptionalInt {
	oint := OptionalInt(&int)
	return &oint
}

func NewNullInt() *OptionalInt {
	oint := OptionalInt(nil)
	return &oint
}

type OptionalTime *Time

func NewOptionalTime(time time.Time) *OptionalTime {
	ot := OptionalTime(&Time{Time: time})
	return &ot
}

func NewNullTime() *OptionalTime {
	ot := OptionalTime(nil)
	return &ot
}

type OptionalIcon *Icon

func NewOptionalIcon(icon Icon) *OptionalIcon {
	oi := OptionalIcon(&icon)
	return &oi
}

func NewNullIcon() *OptionalIcon {
	oi := OptionalIcon(nil)
	return &oi
}

type OptionalFloat *float64

func NewOptionalFloat(float float64) *OptionalFloat {
	ofloat := OptionalFloat(&float)
	return &ofloat
}

func NewNullFloat() *OptionalFloat {
	ofloat := OptionalFloat(nil)
	return &ofloat
}
