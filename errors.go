package filesystem

import "errors"

var (
	// config errors
	ERR_CONFIG_KEY_NOT_FOUND = errors.New("key not found")
	// filesystem errors
	ERR_SOURCE_SAME = errors.New("source and destination are the same")
)
