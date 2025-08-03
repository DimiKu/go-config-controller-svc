package custom_errors

import "errors"

var (
	ErrFieldsContainsBadChars = errors.New("field contains invalid character ';'")
)
