package main

import (
	"database/sql"
	"friend-service/internal/auth-service/app"
	"friend-service/internal/auth-service/config"
	accessgenerator "friend-service/internal/auth-service/infra/access_generator"
	uuid_gen "friend-service/internal/auth-service/infra/id_generator"
	bcrypt_hasher "friend-service/internal/auth-service/infra/password_hasher/bcrypt"
	refreshgenerator "friend-service/internal/auth-service/infra/refresh_generator"
	"friend-service/internal/auth-service/infra/repository/postgres"
	usergrpcclient "friend-service/internal/auth-service/infra/user_client/user_grpc_client"
	grpcapi "friend-service/internal/auth-service/transport/grpc"
	"friend-service/internal/auth-service/usecases/auth"
	"friend-service/internal/auth-service/usecases/token"
	authv1 "friend-service/proto/gen/auth/v1"
	"log"
	"net"
	"time"

	_ "github.com/lib/pq"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if err := initAuthApp(); err != nil {
		log.Fatal(err)
	}
}

func initAuthApp() error {
	cfg, err := config.Load("config/auth.yaml")
	if err != nil {
		return err
	}
	db, err := sql.Open("postgres", cfg.DB.DSN)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	if err = db.Ping(); err != nil {
		return err
	}
	defer db.Close()

	authRepo := postgres.NewAuthRepository(db)
	tokenRepo := postgres.NewTokenRepository(db)

	refreshGenerator := refreshgenerator.NewRefreshTokenGenerator()

	accessGenerator := accessgenerator.NewAccessTokenGenerator([]byte(cfg.JWT.Secret), cfg.JWT.AccessTTL, cfg.JWT.Iss, cfg.JWT.Aud)

	hasher := bcrypt_hasher.NewPasswordHasher(0)

	idGenerator := uuid_gen.NewUUIDGenerator()

	conn, err := grpc.NewClient(
		cfg.Clients.UserServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	userClient := usergrpcclient.New(conn)

	loginService := auth.NewLoginService(authRepo, hasher)
	registerService := auth.NewRegisterService(authRepo, hasher, idGenerator)
	createRefreshTokenService := token.NewCreateRefreshService(hasher, refreshGenerator, accessGenerator, idGenerator, tokenRepo)
	revokeRefreshTokenService := token.NewRevokeRefreshService(hasher, tokenRepo)
	rotateRefreshTokenService := token.NewRotateRefreshService(hasher, refreshGenerator, accessGenerator, idGenerator, tokenRepo)

	authApp := app.NewAuthApp(registerService, loginService, createRefreshTokenService, revokeRefreshTokenService, rotateRefreshTokenService, userClient)

	authServer := grpcapi.NewAuthServer(authApp)

	lis, err := net.Listen("tcp", cfg.Server.GRPCAddr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	authv1.RegisterAuthServiceServer(grpcServer, authServer)
	return grpcServer.Serve(lis)
}
