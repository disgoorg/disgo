package route

import "fmt"

// ErrUnexpectedQueryParam returns a new error for unexpected query parameters
func ErrUnexpectedQueryParam(param string) error {
	return fmt.Errorf("unexpected query param '%s' received", param)
}

// ErrInvalidArgCount returns a new error for invalid argument count
func ErrInvalidArgCount(argCount int, paramCount int) error {
	return fmt.Errorf("invalid amount of arguments received. expected: %d, received: %d", argCount, paramCount)
}

// ErrImageFormatNotSupported returns a new error for if you provide an image format that is not supported for this resource
func ErrImageFormatNotSupported(imageFormat ImageFormat) error {
	return fmt.Errorf("provided image format: %s is not supported on this route", imageFormat)
}
