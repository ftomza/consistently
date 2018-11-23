package consistently

import (
	"errors"
)

var (

	// ErrVersionNotValid Update value not valid version current value
	ErrVersionNotValid = errors.New("version not valid")
)

// Consistently is a struct, embed it in your struct to enable check version for the struct
type Consistently struct {
	Version string `consistently:"version"`
}
