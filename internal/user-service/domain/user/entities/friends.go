package entities

import "github.com/samber/mo"

type Friends struct {
	UserID  string
	Friends []string
}

type FriendsAttrs struct {
	UserID  string
	Friends []string
}

type FriendsFilter struct {
	UserID   mo.Option[string]
	FriendID mo.Option[string]
}
