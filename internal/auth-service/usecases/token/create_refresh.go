package token

import (
	"context"
	"fmt"
	"friend-service/internal/auth-service/domain/auth/entities"
	"friend-service/internal/auth-service/domain/auth/repo"
	"friend-service/internal/auth-service/usecases/token/ports"
	"time"
)

type CreateRefreshRequest struct {
	UserID string
	TTL    time.Duration
}

type CreateRefreshResponse struct {
	RawToken         string
	AccessToken      string
	RefreshExpiresAt time.Time
	AccessExpiresAt  time.Time
}

type CreateRefreshService struct {
	hasher          ports.Hasher
	tokenGenerator  ports.TokenGenerator
	accessGenerator ports.AccessGenerator
	idGen           ports.IDGen
	repo            repo.RefreshTokenRepo
}

func NewCreateRefreshService(hasher ports.Hasher, tokenGenerator ports.TokenGenerator, accessGenerator ports.AccessGenerator, idGen ports.IDGen, repo repo.RefreshTokenRepo) *CreateRefreshService {
	return &CreateRefreshService{
		hasher:          hasher,
		tokenGenerator:  tokenGenerator,
		accessGenerator: accessGenerator,
		idGen:           idGen,
		repo:            repo,
	}
}

func (s *CreateRefreshService) Create(ctx context.Context, req CreateRefreshRequest) (CreateRefreshResponse, error) {
	rawToken, err := s.tokenGenerator.GenerateToken(ctx, 32)
	if err != nil {
		return CreateRefreshResponse{}, err
	}

	hashedToken, err := s.hasher.Hash(ctx, rawToken)
	if err != nil {
		return CreateRefreshResponse{}, err
	}

	refreshExpiresAt := time.Now().Add(req.TTL)

	sessionID, err := s.idGen.NewID()
	if err != nil {
		return CreateRefreshResponse{}, err
	}

	refreshTokenAttrs := entities.RefreshTokenAttrs{
		ID:                 sessionID,
		UserID:             req.UserID,
		HashedRefreshToken: hashedToken,
		ExpiresAt:          refreshExpiresAt,
	}

	savedToken, err := s.repo.Create(ctx, refreshTokenAttrs)
	if err != nil {
		return CreateRefreshResponse{}, err
	}

	accessToken, accessExp, err := s.accessGenerator.GenerateAccessToken(ctx, req.UserID)
	if err != nil {
		return CreateRefreshResponse{}, err
	}

	rawToken = fmt.Sprintf("%s.%s", savedToken.ID, rawToken)
	return CreateRefreshResponse{
		RawToken:    rawToken,
		AccessToken: accessToken,
		RefreshExpiresAt:   refreshExpiresAt,
		AccessExpiresAt:    accessExp,
	}, nil
}
