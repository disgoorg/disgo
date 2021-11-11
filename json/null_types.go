package json

type NullBool *bool

func NewBool(b bool) *NullBool {
	v := NullBool(&b)
	return &v
}

func NewNullBool() *NullBool {
	v := NullBool(nil)
	return &v
}

type NullString *string

func NewString(str string) *NullString {
	v := NullString(&str)
	return &v
}

func NewNullString() *NullString {
	v := NullString(nil)
	return &v
}

type NullInt *int

func NewInt(int int) *NullInt {
	v := NullInt(&int)
	return &v
}

func NewNullInt() *NullInt {
	v := NullInt(nil)
	return &v
}

type NullFloat *float64

func NewFloat(float float64) *NullFloat {
	v := NullFloat(&float)
	return &v
}

func NewNullFloat() *NullFloat {
	v := NullFloat(nil)
	return &v
}
