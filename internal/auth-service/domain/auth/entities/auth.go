package entities

import "github.com/samber/mo"

type AuthCredentials struct {
	UserID         string
	Email          string
	HashedPassword string
}

type AuthCredentialsAttrs struct {
	UserID         string
	Email          string
	HashedPassword string
}

type AuthCredentialsFilter struct {
	UserID   mo.Option[string]
	Email    mo.Option[string]
}
