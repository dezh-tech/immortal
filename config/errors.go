package config

import "fmt"

// Error represents an error in loading or validating config.
type Error struct {
	reason string
}

func (e Error) Error() string {
	return fmt.Sprintf("config error: %s\n", e.reason)
}
