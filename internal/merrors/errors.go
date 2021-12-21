package merrors

import "sync"

var _ error = (*Error)(nil)

type Error struct {
	errs []error
	mu   *sync.Mutex
}

func (e Error) Error() string {
	e.mu.Lock()
	defer e.mu.Unlock()
	str := "["
	for i := range e.errs {
		str += e.errs[i].Error()
		if i < len(e.errs) {
			str += ", "
		}
	}
	return str + "]"
}

func (e *Error) Add(err error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.errs = append(e.errs, err)
}
