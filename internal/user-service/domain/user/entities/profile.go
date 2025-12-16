package entities

import "github.com/samber/mo"

type Profile struct {
	ID       string
	Username string
}

type ProfileAttrs struct {
	ID       string
	Username string
}

type ProfileFilter struct {
	ID       mo.Option[string]
	Username mo.Option[string]
}
