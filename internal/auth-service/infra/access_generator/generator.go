package accessgenerator

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenGenerator struct {
	secret []byte
	ttl    time.Duration
	iss    string
	aud    string
}

func NewAccessTokenGenerator(secret []byte, ttl time.Duration, iss string, aud string) *AccessTokenGenerator {
	return &AccessTokenGenerator{
		secret: secret,
		ttl:    ttl,
		iss:    iss,
		aud:    aud,
	}
}

func (g *AccessTokenGenerator) GenerateAccessToken(ctx context.Context, userID string) (token string, exp time.Time, err error) {
	now := time.Now()
	exp = now.Add(g.ttl)

	claims := jwt.MapClaims{
		"sub": userID,
		"iss": g.iss,
		"aud": g.aud,
		"iat": now.Unix(),
		"exp": exp.Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := jwtToken.SignedString(g.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, exp, nil
}
