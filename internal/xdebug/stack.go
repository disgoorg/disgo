package xdebug

import (
	"bytes"
	"runtime/debug"
)

// Stack returns a formatted stack trace of the goroutine that calls it.
// The skip parameter specifies the number of stack frames to skip
// before recording in the returned trace; 0 means to include the
// caller of Stack.
func Stack(skip int) []byte {
	return filterStack(debug.Stack(), skip)
}

func filterStack(stack []byte, skip int) []byte {
	if skip == 0 {
		return stack
	}

	lines := bytes.Split(stack, []byte{'\n'})
	var b bytes.Buffer

	// write the first line (goroutine info)
	b.Write(lines[0])
	b.WriteByte('\n')
	lines = lines[1:]

	var frameCount int
	for {
		if frameCount >= skip || len(lines) == 0 {
			break
		}

		frameCount++
		lines = lines[1:]
		if len(lines) > 1 {
			lines = lines[1:]
		}
	}
	b.Write(bytes.Join(lines, []byte{'\n'}))

	return b.Bytes()
}
