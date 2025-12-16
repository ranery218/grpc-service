package repo

import (
	"context"
	"friend-service/internal/auth-service/domain/auth/entities"
)

type AuthRepo interface {
	Create(ctx context.Context, attrs entities.AuthCredentialsAttrs) (entities.AuthCredentials, error)
	Get(ctx context.Context, filter entities.AuthCredentialsFilter) (entities.AuthCredentials, error)
	Update(ctx context.Context, userID string, attrs entities.AuthCredentialsAttrs) (entities.AuthCredentials, error)
	Delete(ctx context.Context, userID string) error
}
