package user

import "errors"

var (
	ErrUserIDRequired = errors.New("user ID is required")
	ErrFriendIDsEmpty = errors.New("friend IDs list cannot be empty")
	ErrNoneFriends    = errors.New("no friends found")
)
