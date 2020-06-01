package passwd

import "fmt"

// ErrNotFound is used when an entry is not found.
type ErrNotFound struct {
	err string
}

// Error takes an error and returns a string. Satisfies the interface.
func (e *ErrNotFound) Error() string {
	return fmt.Sprintf(e.err)
}
