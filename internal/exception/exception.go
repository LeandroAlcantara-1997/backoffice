package exception

import "errors"

var (
	ErrProcessingTimeout = errors.New("processing timeout error")
)
