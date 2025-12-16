package auth

import (
	"context"
	"friend-service/internal/auth-service/domain/auth/entities"
	"friend-service/internal/auth-service/domain/auth/repo"
	"friend-service/internal/auth-service/usecases/auth/ports"

	"github.com/samber/mo"
)

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	UserID string
}

type LoginService struct {
	repo   repo.AuthRepo
	hasher ports.PasswordHasher
}

func NewLoginService(repo repo.AuthRepo, hasher ports.PasswordHasher) *LoginService {
	return &LoginService{repo: repo, hasher: hasher}
}

func (s *LoginService) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	if err := ValidateEmail(req.Email); err != nil {
		return LoginResponse{}, err
	}
	if err := ValidatePassword(req.Password); err != nil {
		return LoginResponse{}, err
	}

	email := req.Email
	password := req.Password

	creds, err := s.repo.Get(ctx, entities.AuthCredentialsFilter{Email: mo.Some(email)})
	if err != nil {
		return LoginResponse{}, err
	}

	if err := s.hasher.Compare(ctx, creds.HashedPassword, password); err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{UserID: creds.UserID}, nil
}
