package bcrypt_hasher

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher struct {
	cost int
}

func NewPasswordHasher(cost int) *PasswordHasher {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	return &PasswordHasher{cost: cost}
}

func (h *PasswordHasher) Hash(ctx context.Context, plaintext string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), h.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (h *PasswordHasher) Compare(ctx context.Context, hash string, plaintext string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext))
}
