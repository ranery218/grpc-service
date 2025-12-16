package repo

import (
	"context"
	"friend-service/internal/auth-service/domain/auth/entities"
)

type RefreshTokenRepo interface {
	Create(ctx context.Context, attrs entities.RefreshTokenAttrs) (entities.RefreshToken, error)
	Get(ctx context.Context, filter entities.RefreshTokenFilter) (entities.RefreshToken, error)
	Delete(ctx context.Context, filter entities.RefreshTokenFilter) error
}
