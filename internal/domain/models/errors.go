package models

import "fmt"

var (
	ErrForbidden = fmt.Errorf("forbidden")
	ErrNotFound = fmt.Errorf("not found")
)
