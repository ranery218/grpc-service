package token

import (
	"context"
	"fmt"
	"friend-service/internal/auth-service/domain/auth/entities"
	"friend-service/internal/auth-service/domain/auth/repo"
	"friend-service/internal/auth-service/usecases/token/ports"
	"time"

	"github.com/samber/mo"
)

type RotateRefreshRequest struct {
	RawToken string
	TTL      time.Duration
}

type RotateRefreshResponse struct {
	RawToken         string
	AccessToken      string
	RefreshExpiresAt time.Time
	AccessExpiresAt  time.Time
}

type RotateRefreshService struct {
	hasher          ports.Hasher
	tokenGenerator  ports.TokenGenerator
	accessGenerator ports.AccessGenerator
	idGen           ports.IDGen
	repo            repo.RefreshTokenRepo
}

func NewRotateRefreshService(hasher ports.Hasher, tokenGenerator ports.TokenGenerator, accessGenerator ports.AccessGenerator, idGen ports.IDGen, repo repo.RefreshTokenRepo) *RotateRefreshService {
	return &RotateRefreshService{
		hasher:          hasher,
		tokenGenerator:  tokenGenerator,
		accessGenerator: accessGenerator,
		idGen:           idGen,
		repo:            repo,
	}
}

func (s *RotateRefreshService) Rotate(ctx context.Context, req RotateRefreshRequest) (RotateRefreshResponse, error) {
	id, secret, err := parseRefresh(req.RawToken)
	if err != nil {
		return RotateRefreshResponse{}, ErrInvalidRefreshToken
	}

	storedToken, err := s.repo.Get(ctx, entities.RefreshTokenFilter{ID: mo.Some(id)})
	if err != nil {
		return RotateRefreshResponse{}, ErrInvalidRefreshToken
	}

	if err = s.hasher.Compare(ctx, storedToken.HashedRefreshToken, secret); err != nil {
		return RotateRefreshResponse{}, ErrInvalidRefreshToken
	}

	if storedToken.ExpiresAt.Before(time.Now()) {
		return RotateRefreshResponse{}, ErrInvalidRefreshToken
	}

	newRawToken, err := s.tokenGenerator.GenerateToken(ctx, 32)
	if err != nil {
		return RotateRefreshResponse{}, err
	}
	newHashedToken, err := s.hasher.Hash(ctx, newRawToken)
	if err != nil {
		return RotateRefreshResponse{}, err
	}

	newExpiresAt := time.Now().Add(req.TTL)

	newID, err := s.idGen.NewID()
	if err != nil {
		return RotateRefreshResponse{}, err
	}

	newRefreshTokenAttrs := entities.RefreshTokenAttrs{
		ID:                 newID,
		UserID:             storedToken.UserID,
		HashedRefreshToken: newHashedToken,
		ExpiresAt:          newExpiresAt,
	}
	_, err = s.repo.Create(ctx, newRefreshTokenAttrs)
	if err != nil {
		return RotateRefreshResponse{}, err
	}

	err = s.repo.Delete(ctx, entities.RefreshTokenFilter{ID: mo.Some(storedToken.ID)})
	if err != nil {
		return RotateRefreshResponse{}, err
	}

	accessToken, accessExp, err := s.accessGenerator.GenerateAccessToken(ctx, storedToken.UserID)
	if err != nil {
		return RotateRefreshResponse{}, err
	}

	newRawToken = fmt.Sprintf("%s.%s", newID, newRawToken)

	return RotateRefreshResponse{
		RawToken:    newRawToken,
		AccessToken: accessToken,
		RefreshExpiresAt:   newExpiresAt,
		AccessExpiresAt:    accessExp,
	}, nil
}
