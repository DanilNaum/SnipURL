package url

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrIDIsBusy = errors.New("id is busy")

	ErrConflict = errors.New("conflict")
)
