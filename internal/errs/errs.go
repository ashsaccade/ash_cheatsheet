package errs

import (
	"github.com/pkg/errors"
)

var (
	ErrEmptyCardName = errors.New("empty card name")
	ErrEmptyCardDesc = errors.New("empty description name")
)
