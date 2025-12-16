package refreshgenerator

import (
	"context"
	"crypto/rand"
	"encoding/base64"
)

type RefreshTokenGenerator struct{}

func NewRefreshTokenGenerator() *RefreshTokenGenerator {
	return &RefreshTokenGenerator{}
}

func (g *RefreshTokenGenerator) GenerateToken(ctx context.Context, nBytes int) (string, error) {
	b := make([]byte, nBytes)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
