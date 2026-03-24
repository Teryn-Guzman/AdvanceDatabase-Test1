package data

import (
	"errors"
)

var ErrRecordNotFound = errors.New("record not found")

var ErrDuplicateEmail = errors.New("a customer with this email already exists")
