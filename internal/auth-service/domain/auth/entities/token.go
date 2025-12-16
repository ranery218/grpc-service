package entities

import (
	"time"

	"github.com/samber/mo"
)

type RefreshToken struct {
	ID                 string
	UserID             string
	HashedRefreshToken string
	ExpiresAt          time.Time
}

type RefreshTokenAttrs struct {
	ID                 string
	UserID             string
	HashedRefreshToken string
	ExpiresAt          time.Time
}

type RefreshTokenFilter struct {
	ID                 mo.Option[string]
	UserID             mo.Option[string]
	HashedRefreshToken mo.Option[string]
}
