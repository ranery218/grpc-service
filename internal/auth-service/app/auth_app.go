package app

import (
	"context"
	"friend-service/internal/auth-service/usecases/auth"
	"friend-service/internal/auth-service/usecases/token"
	"time"
)

type UserClient interface {
	CreateProfile(ctx context.Context, userID string, username string) error
}

type RegisterDTO struct {
	Email    string
	Username string
	Password string
}

type LoginDTO struct {
	Email    string
	Password string
}

type LogoutDTO struct {
	RefreshToken string
}

type RefreshDTO struct {
	RefreshToken string
}

type AuthTokens struct {
	AccessToken  string
	AccessExp    time.Time
	RefreshToken string
	RefreshExp   time.Time
}

type AuthApp struct {
	registerService           *auth.RegisterService
	loginService              *auth.LoginService
	createRefreshTokenService *token.CreateRefreshService
	revokeRefreshTokenService *token.RevokeRefreshService
	rotateRefreshTokenService *token.RotateRefreshService
	userClient                UserClient
}

func NewAuthApp(
	registerService *auth.RegisterService,
	loginService *auth.LoginService,
	createRefreshTokenService *token.CreateRefreshService,
	revokeRefreshTokenService *token.RevokeRefreshService,
	rotateRefreshTokenService *token.RotateRefreshService,
	userClient                UserClient,
) *AuthApp {
	return &AuthApp{
		registerService:           registerService,
		loginService:              loginService,
		createRefreshTokenService: createRefreshTokenService,
		revokeRefreshTokenService: revokeRefreshTokenService,
		rotateRefreshTokenService: rotateRefreshTokenService,
		userClient:                userClient,
	}
}

func (app *AuthApp) Register(ctx context.Context, dto RegisterDTO) (AuthTokens, error) {
	if err := auth.ValidateEmail(dto.Email); err != nil {
		return AuthTokens{}, err
	}
	if err := auth.ValidatePassword(dto.Password); err != nil {
		return AuthTokens{}, err
	}
	if err := auth.ValidateUsername(dto.Username); err != nil {
		return AuthTokens{}, err
	}

	registerReq := auth.RegisterRequest{
		Password: dto.Password,
		Email:    dto.Email,
	}
	resp, err := app.registerService.Register(ctx, registerReq)
	if err != nil {
		return AuthTokens{}, err
	}

	userID := resp.UserID
	createTokenReq := token.CreateRefreshRequest{
		UserID: userID,
		TTL:    time.Hour * 24 * 7,
	}
	tokens, err := app.createRefreshTokenService.Create(ctx, createTokenReq)
	if err != nil {
		return AuthTokens{}, err
	}

	err = app.userClient.CreateProfile(ctx, userID, dto.Username)
	if err != nil {
		return AuthTokens{}, err
	}

	return AuthTokens{
		AccessToken:  tokens.AccessToken,
		AccessExp:    tokens.AccessExpiresAt,
		RefreshToken: tokens.RawToken,
		RefreshExp:   tokens.RefreshExpiresAt,
	}, nil
}

func (app *AuthApp) Login(ctx context.Context, dto LoginDTO) (AuthTokens, error) {
	if err := auth.ValidateEmail(dto.Email); err != nil {
		return AuthTokens{}, err
	}
	if err := auth.ValidatePassword(dto.Password); err != nil {
		return AuthTokens{}, err
	}

	loginReq := auth.LoginRequest{
		Email:    dto.Email,
		Password: dto.Password,
	}

	resp, err := app.loginService.Login(ctx, loginReq)
	if err != nil {
		return AuthTokens{}, err
	}

	userID := resp.UserID
	createTokenReq := token.CreateRefreshRequest{
		UserID: userID,
		TTL:    time.Hour * 24 * 7,
	}
	tokens, err := app.createRefreshTokenService.Create(ctx, createTokenReq)
	if err != nil {
		return AuthTokens{}, err
	}

	return AuthTokens{
		AccessToken:  tokens.AccessToken,
		AccessExp:    tokens.AccessExpiresAt,
		RefreshToken: tokens.RawToken,
		RefreshExp:   tokens.RefreshExpiresAt,
	}, nil
}

func (app *AuthApp) Logout(ctx context.Context, dto LogoutDTO) error {
	revokeReq := token.RevokeRefreshRequest{
		RawToken: dto.RefreshToken,
	}
	return app.revokeRefreshTokenService.Revoke(ctx, revokeReq)
}

func (app *AuthApp) Refresh(ctx context.Context, dto RefreshDTO) (AuthTokens, error) {
	rotateReq := token.RotateRefreshRequest{
		RawToken: dto.RefreshToken,
		TTL:      time.Hour * 24 * 7,
	}

	tokens, err := app.rotateRefreshTokenService.Rotate(ctx, rotateReq)
	if err != nil {
		return AuthTokens{}, err
	}

	return AuthTokens{
		AccessToken:  tokens.AccessToken,
		AccessExp:    tokens.AccessExpiresAt,
		RefreshToken: tokens.RawToken,
		RefreshExp:   tokens.RefreshExpiresAt,
	}, nil
}
