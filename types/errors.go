package types

import "fmt"

// ErrEncode represents an encoding error.
type ErrEncode struct {
	Reason string
}

func (e ErrEncode) Error() string {
	return fmt.Sprintf("encoding error: %s", e.Reason)
}

// ErrDecode represents an decoding error.
type ErrDecode struct {
	Reason string
}

func (e ErrDecode) Error() string {
	return fmt.Sprintf("decoding error: %s", e.Reason)
}
