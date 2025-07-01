package url

import "errors"

// Predefined errors used in the URL repository to represent various error conditions
// when performing operations on URL resources.
var (
	// ErrNotFound indicates that the requested resource was not found
	ErrNotFound = errors.New("not found")
	// ErrIDIsBusy indicates that the specified ID is currently in use or unavailable
	ErrIDIsBusy = errors.New("id is busy")
	// ErrDeleted indicates that the resource has been previously deleted
	ErrDeleted = errors.New("deleted")
	// ErrConflict indicates a conflict occurred, typically due to a concurrent modification or constraint violation
	ErrConflict = errors.New("conflict")
)
