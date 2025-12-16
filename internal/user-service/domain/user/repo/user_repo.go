package repo

import (
	"context"
	"friend-service/internal/user-service/domain/user/entities"
)

type UserRepo interface {
	Create(ctx context.Context, attrs entities.ProfileAttrs) (entities.Profile, error)
	Get(ctx context.Context, filter entities.ProfileFilter) (entities.Profile, error)
	GetAll(ctx context.Context, filter entities.ProfileFilter) ([]entities.Profile, error)
	Update(ctx context.Context, filter entities.ProfileFilter, attrs entities.ProfileAttrs) (entities.Profile, error)
	Delete(ctx context.Context, filter entities.ProfileFilter) error
}
