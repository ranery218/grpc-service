package grpcapi

import (
	"context"
	"errors"
	"friend-service/internal/auth-service/app"
	"friend-service/internal/auth-service/domain/auth/entities"
	"friend-service/internal/auth-service/usecases/token"
	authv1 "friend-service/proto/gen/auth/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	authv1.UnimplementedAuthServiceServer
	authApp *app.AuthApp
}

func NewAuthServer(authApp *app.AuthApp) *AuthServer {
	return &AuthServer{authApp: authApp}
}

func (s *AuthServer) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	username := req.GetUsername()
	password := req.GetPassword()
	email := req.GetEmail()

	if username == "" || password == "" || email == "" {
		return nil, status.Error(codes.InvalidArgument, "username, password and email are required")
	}

	registerDTO := app.RegisterDTO{
		Email:    email,
		Username: username,
		Password: password,
	}

	tokens, err := s.authApp.Register(ctx, registerDTO)
	if err != nil {
		return nil, mapErr(err)
	}

	return &authv1.RegisterResponse{
		AccessToken:      tokens.AccessToken,
		AccessExpiresAt:  tokens.AccessExp.Unix(),
		RefreshToken:     tokens.RefreshToken,
		RefreshExpiresAt: tokens.RefreshExp.Unix(),
	}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	email := req.GetEmail()
	password := req.GetPassword()
	if email == "" || password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	loginDTO := app.LoginDTO{
		Email:    email,
		Password: password,
	}

	tokens, err := s.authApp.Login(ctx, loginDTO)
	if err != nil {
		return nil, mapErr(err)
	}

	return &authv1.LoginResponse{
		AccessToken:      tokens.AccessToken,
		AccessExpiresAt:  tokens.AccessExp.Unix(),
		RefreshToken:     tokens.RefreshToken,
		RefreshExpiresAt: tokens.RefreshExp.Unix(),
	}, nil
}

func (s *AuthServer) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	refreshToken := req.GetRefreshToken()
	if refreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	logoutDTO := app.LogoutDTO{
		RefreshToken: refreshToken,
	}
	err := s.authApp.Logout(ctx, logoutDTO)
	if err != nil {
		return nil, mapErr(err)
	}

	return &authv1.LogoutResponse{}, nil
}

func (s *AuthServer) Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.RefreshResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	refreshToken := req.GetRefreshToken()
	if refreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	refreshDTO := app.RefreshDTO{
		RefreshToken: refreshToken,
	}

	tokens, err := s.authApp.Refresh(ctx, refreshDTO)
	if err != nil {
		return nil, mapErr(err)
	}

	return &authv1.RefreshResponse{
		AccessToken:      tokens.AccessToken,
		AccessExpiresAt:  tokens.AccessExp.Unix(),
		RefreshToken:     tokens.RefreshToken,
		RefreshExpiresAt: tokens.RefreshExp.Unix(),
	}, nil
}

func mapErr(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, token.ErrInvalidRefreshToken):
		return status.Error(codes.Unauthenticated, "invalid refresh token")
	case errors.Is(err, entities.ErrRefreshTokenNotFound):
		return status.Error(codes.Unauthenticated, "invalid refresh token")
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
