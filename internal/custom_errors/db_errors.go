package custom_errors

import "errors"

var (
	ErrNotLockedConfigNotFound = errors.New("not locked config not found")
)
