package json

import "time"

type NullBool *bool

//goland:noinspection GoUnusedExportedFunction
func NewBool(b bool) *NullBool {
	v := NullBool(&b)
	return &v
}

//goland:noinspection GoUnusedExportedFunction
func NewNullBool() *NullBool {
	v := NullBool(nil)
	return &v
}

type NullString *string

//goland:noinspection GoUnusedExportedFunction
func NewString(str string) *NullString {
	v := NullString(&str)
	return &v
}

//goland:noinspection GoUnusedExportedFunction
func NewNullString() *NullString {
	v := NullString(nil)
	return &v
}

type NullInt *int

//goland:noinspection GoUnusedExportedFunction
func NewInt(int int) *NullInt {
	v := NullInt(&int)
	return &v
}

//goland:noinspection GoUnusedExportedFunction
func NewNullInt() *NullInt {
	v := NullInt(nil)
	return &v
}

type NullFloat *float64

//goland:noinspection GoUnusedExportedFunction
func NewFloat(float float64) *NullFloat {
	v := NullFloat(&float)
	return &v
}

//goland:noinspection GoUnusedExportedFunction
func NewNullFloat() *NullFloat {
	v := NullFloat(nil)
	return &v
}

type NullTime *time.Time

//goland:noinspection GoUnusedExportedFunction
func NewTime(time time.Time) *NullTime {
	v := NullTime(&time)
	return &v
}

//goland:noinspection GoUnusedExportedFunction
func NewNullTime() *NullTime {
	v := NullTime(nil)
	return &v
}
