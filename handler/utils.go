package handler

import "strings"

type handlerHolder[T any] struct {
	path    string
	handler T
}

func findHandler[T any](handlers []*handlerHolder[T], path string, delimiter string) (T, map[string]string, bool) {
	parts := strings.Split(path, delimiter)
	values := map[string]string{}

	for _, h := range handlers {
		handlerParts := strings.Split(h.path, delimiter)
		if len(parts) != len(handlerParts) {
			continue
		}
		matches := true
		for i, part := range parts {
			handlerPart := handlerParts[i]
			if strings.HasPrefix(handlerPart, "{") && strings.HasSuffix(handlerParts[i], "}") {
				values[handlerPart[1:len(handlerPart)-1]] = part
				continue
			}
			if part != handlerPart {
				matches = false
				break
			}
		}
		if !matches {
			continue
		}
		return h.handler, values, true
	}
	var emptyHandler T
	return emptyHandler, nil, false
}
