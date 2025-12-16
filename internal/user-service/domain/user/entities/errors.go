package entities
import "errors"

var (
	ErrUsersFriendsNotFound = errors.New("user's friends not found")
	ErrEmptyFilter		  = errors.New("at least one filter must be provided")
)
