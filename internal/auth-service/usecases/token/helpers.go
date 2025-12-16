package token

import (
	"errors"
	"regexp"
	"strings"
)

func parseRefresh(raw string) (id, secret string, err error) {
	var tokenRe = regexp.MustCompile(`^[A-Za-z0-9_-]{8,128}$`)
	id, secret, ok := strings.Cut(raw, ".")
	switch {
	case !ok:
		return "", "", errors.New("invalid refresh format")
	case !tokenRe.MatchString(id):
		return "", "", errors.New("invalid token id")
	case !tokenRe.MatchString(secret):
		return "", "", errors.New("invalid token secret")
	}
	return id, secret, nil
}
