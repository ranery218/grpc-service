package repo

import (
	"context"
	"friend-service/internal/user-service/domain/user/entities"
)

type FriendsRepo interface {
	Create(ctx context.Context, attrs entities.FriendsAttrs) (entities.Friends, error)
	Get(ctx context.Context, filter entities.FriendsFilter) (entities.Friends, error)
	Update(ctx context.Context, filter entities.FriendsFilter, attrs entities.FriendsAttrs) (entities.Friends, error)
}
