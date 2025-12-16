package main

import (
	"database/sql"
	"log"
	"net"
	"time"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"friend-service/internal/user-service/config"
	"friend-service/internal/user-service/infra/repository/postrges"
	transport "friend-service/internal/user-service/transport"
	"friend-service/internal/user-service/usecases/user"
	userv1 "friend-service/proto/gen/user/v1"
)

func main() {
	if err := initUserApp(); err != nil {
		log.Fatal(err)
	}
}

func initUserApp() error {
	cfg, err := config.Load("config/user.yaml")
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

	profileRepo := postrges.NewProfileRepository(db)
	friendsRepo := postrges.NewFriendsRepository(db)

	createProfileService := user.NewCreateProfileService(profileRepo)
	getProfileService := user.NewGetProfileService(profileRepo)
	updateProfileService := user.NewUpdateProfileService(profileRepo)
	deleteProfileService := user.NewDeleteProfileService(profileRepo)
	getAllProfilesService := user.NewGetAllProfilesService(profileRepo)
	addFriendService := user.NewAddFriendService(friendsRepo)
	getFriendsService := user.NewGetFriendsService(friendsRepo)
	deleteFriendService := user.NewDeleteFriendService(friendsRepo)

	userServer := transport.NewUserServer(
		createProfileService,
		getProfileService,
		updateProfileService,
		deleteProfileService,
		getAllProfilesService,
		addFriendService,
		getFriendsService,
		deleteFriendService,
	)

	lis, err := net.Listen("tcp", cfg.Server.GRPCAddr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	userv1.RegisterUserServiceServer(grpcServer, userServer)

	return grpcServer.Serve(lis)
}
