package sweep

import "errors"

// ErrClosed is the error returned when sweep is closed already.
var ErrClosed = errors.New("sweep closed")

// ErrEntryNotFound is the error returned when an
// entry doesn't exist in sweep.
var ErrEntryNotFound = errors.New("entry not found")

// ErrEntryTooLarge is the error returned when an entry is too large
// going to be put in sweep.
var ErrEntryTooLarge = errors.New("entry is too large in size")
