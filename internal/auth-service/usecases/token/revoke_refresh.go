package token

import (
	"context"
	"errors"
	"friend-service/internal/auth-service/domain/auth/entities"
	"friend-service/internal/auth-service/domain/auth/repo"
	"friend-service/internal/auth-service/usecases/token/ports"

	"github.com/samber/mo"
)

type RevokeRefreshRequest struct {
	RawToken string
}

type RevokeRefreshService struct {
	hasher ports.Hasher
	repo   repo.RefreshTokenRepo
}

func NewRevokeRefreshService(hasher ports.Hasher, repo repo.RefreshTokenRepo) *RevokeRefreshService {
	return &RevokeRefreshService{
		hasher: hasher,
		repo:   repo,
	}
}

func (s *RevokeRefreshService) Revoke(ctx context.Context, req RevokeRefreshRequest) error {
	id, secret, err := parseRefresh(req.RawToken)
	if err != nil {
		return err
	}

	savedToken, err := s.repo.Get(ctx, entities.RefreshTokenFilter{ID: mo.Some(id)})
	if err != nil {
		if errors.Is(err, entities.ErrRefreshTokenNotFound) {
			return nil
		}
		return err
	}

	hashedToken := savedToken.HashedRefreshToken

	if err := s.hasher.Compare(ctx, hashedToken, secret); err != nil {
		return errors.New("invalid refresh token")
	}

	return s.repo.Delete(ctx, entities.RefreshTokenFilter{ID: mo.Some(id)})
}
