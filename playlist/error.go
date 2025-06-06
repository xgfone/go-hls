package playlist

import "fmt"

// ParseError represents an error when parsing a master/media playlist.
type ParseError struct {
	Line int
	Err  error
}

// Error implements the error interface.
func (e ParseError) Error() string {
	return fmt.Sprintf("line %d: %v", e.Line, e.Err)
}

// Unwrap is used to unwrap the inner error.
func (e ParseError) Unwrap() error {
	return e.Err
}
