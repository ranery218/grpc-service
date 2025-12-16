package ports

import (
	"context"
	"time"
)

type AccessGenerator interface {
	GenerateAccessToken(ctx context.Context, userID string) (token string, exp time.Time, err error)
}
