package custom_errors

import "errors"

var (
	ErrUnauthorized    = errors.New("unauthorized")
	ErrOrderNotFound   = errors.New("order not found")
	ErrOrderCompleted  = errors.New("order already completed")
	ErrInvalidObjectID = errors.New("invalid object ID")
)
