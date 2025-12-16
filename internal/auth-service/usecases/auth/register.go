package auth

import (
	"context"
	"errors"
	"friend-service/internal/auth-service/domain/auth/entities"
	"friend-service/internal/auth-service/domain/auth/repo"
	"friend-service/internal/auth-service/usecases/auth/ports"

	"github.com/samber/mo"
)

type RegisterRequest struct {
	Password string
	Email    string
}

type RegisterResponse struct {
	UserID string
}

type RegisterService struct {
	repo   repo.AuthRepo
	hasher ports.PasswordHasher
	idGen  ports.IDGen
}

func NewRegisterService(repo repo.AuthRepo, hasher ports.PasswordHasher, idGen ports.IDGen) *RegisterService {
	return &RegisterService{repo: repo, hasher: hasher, idGen: idGen}
}

func (s *RegisterService) Register(ctx context.Context, req RegisterRequest) (RegisterResponse, error) {
	if err := ValidateEmail(req.Email); err != nil {
		return RegisterResponse{}, err
	}

	if err := ValidatePassword(req.Password); err != nil {
		return RegisterResponse{}, err
	}

	email := req.Email

	if _, err := s.repo.Get(ctx, entities.AuthCredentialsFilter{Email: mo.Some(email)}); err == nil {
		return RegisterResponse{}, errors.New("email is taken")
	}

	hashedPassword, err := s.hasher.Hash(ctx, req.Password)
	if err != nil {
		return RegisterResponse{}, err
	}

	userID, err := s.idGen.NewID()
	if err != nil {
		return RegisterResponse{}, err
	}

	creds := entities.AuthCredentialsAttrs{
		UserID:         userID,
		Email:          email,
		HashedPassword: hashedPassword,
	}

	userEntities, err := s.repo.Create(ctx, creds)
	if err != nil {
		return RegisterResponse{}, err
	}

	return RegisterResponse{UserID: userEntities.UserID}, nil
}
