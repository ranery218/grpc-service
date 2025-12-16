package main

import (
	"context"
	"friend-service/internal/gateway/config"
	"friend-service/internal/gateway/middleware"
	authv1 "friend-service/proto/gen/auth/v1"
	userv1 "friend-service/proto/gen/user/v1"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if err := RunServer(); err !=nil {
		log.Fatal(err)
	}
}

func RunServer() error {
	config, err := config.Load("config/gateway.yaml")
	if err != nil {
		log.Fatal(err)
	}

	mux := runtime.NewServeMux(runtime.WithMetadata(middleware.AuthMetadata))

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	ctx := context.Background()

	if err := authv1.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, config.Backends.AuthServiceAddr, opts); err != nil {
		return err
	}
	if err := userv1.RegisterUserServiceHandlerFromEndpoint(ctx, mux, config.Backends.UserServiceAddr, opts); err != nil {
		return err
	}

	handler := middleware.AuthMiddleware(mux, []byte(config.JWT.Secret))

	srv := &http.Server{
		Addr:         config.Server.HTTPAddr,
		Handler:      handler,
		ReadTimeout:  config.Server.ReadTimeout,
		WriteTimeout: config.Server.WriteTimeout,
	}
	log.Printf("gateway on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
	return nil
}
