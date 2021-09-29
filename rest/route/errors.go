package route

import "fmt"

func ErrUnexpectedQueryParam(param string) error {
	return fmt.Errorf("unexpected query param '%s' received", param)
}

func ErrInvalidArgCount(argCount int, paramCount int) error {
	return fmt.Errorf("invalid amount of arguments received. expected: %d, received: %d", argCount, paramCount)
}

func ErrFileExtensionNotSupported(fileExtension string) error {
	return fmt.Errorf("provided file extension: %s is not supported by discord on this end", fileExtension)
}
