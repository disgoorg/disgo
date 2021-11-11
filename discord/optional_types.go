package discord

import (
	"time"
)

type NullTime *Time

func NewTime(time time.Time) *NullTime {
	ot := NullTime(&Time{Time: time})
	return &ot
}

func NewNullTime() *NullTime {
	ot := NullTime(nil)
	return &ot
}

type NullIcon *Icon

func NewNIcon(icon Icon) *NullIcon {
	oi := NullIcon(&icon)
	return &oi
}

func NewNullIcon() *NullIcon {
	oi := NullIcon(nil)
	return &oi
}
