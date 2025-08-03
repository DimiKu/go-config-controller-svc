package custom_errors

import "errors"

var (
	ErrNotLockedConfigNotFound = errors.New("not locked config not found")
	ErrUserAlreadyExist        = errors.New("user already exist")
	ErrLoginError              = errors.New("login failed")
	ErrWrongTask               = errors.New("wrong task")
)
