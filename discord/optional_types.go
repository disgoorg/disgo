package discord

import (
	"time"
)

type NullTime *Time

//goland:noinspection GoUnusedExportedFunction
func NewTime(time time.Time) *NullTime {
	ot := NullTime(&Time{Time: time})
	return &ot
}

//goland:noinspection GoUnusedExportedFunction
func NewNullTime() *NullTime {
	ot := NullTime(nil)
	return &ot
}

type NullIcon *Icon

//goland:noinspection GoUnusedExportedFunction
func NewNIcon(icon Icon) *NullIcon {
	oi := NullIcon(&icon)
	return &oi
}

//goland:noinspection GoUnusedExportedFunction
func NewNullIcon() *NullIcon {
	oi := NullIcon(nil)
	return &oi
}
