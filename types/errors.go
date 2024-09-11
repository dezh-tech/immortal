package types

import "fmt"

// EncodeError represents an encoding error.
type EncodeError struct {
	Reason string
}

func (e EncodeError) Error() string {
	return fmt.Sprintf("encoding error: %s", e.Reason)
}

// DecodeError represents an decoding error.
type DecodeError struct {
	Reason string
}

func (e DecodeError) Error() string {
	return fmt.Sprintf("decoding error: %s", e.Reason)
}
