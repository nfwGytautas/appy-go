package appy_errors

import "errors"

// Error: No Logger provider specified
var ErrNoLogger = errors.New("no Logger provider specified")
