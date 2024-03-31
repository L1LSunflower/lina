package tools

import "errors"

// Custom errors
var (
	ErrGetDepends = errors.New("failed to get depends from context")
)
