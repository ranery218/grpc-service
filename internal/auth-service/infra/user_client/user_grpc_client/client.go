package usergrpcclient

import (
    "context"

    userv1 "friend-service/proto/gen/user/v1"
    "google.golang.org/grpc"
)

type Client struct {
    cli userv1.UserServiceClient
}

func New(conn *grpc.ClientConn) *Client {
    return &Client{cli: userv1.NewUserServiceClient(conn)}
}

func (c *Client) CreateProfile(ctx context.Context, userID, username string) error {
    _, err := c.cli.CreateProfile(ctx, &userv1.CreateProfileRequest{
        UserId:   userID,
        Username: username,
    })
    return err
}
