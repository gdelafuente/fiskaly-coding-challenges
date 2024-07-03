package commands

import (
	"errors"
)

var (
	ErrValidation            = errors.New("invalid command or query")
	ErrAlgorithmNotSupported = errors.New("algorithm not supported")
	ErrSavingDevice          = errors.New("failed to save device")
)
