package user

import (
	"context"
	"friend-service/internal/user-service/domain/user/entities"
	"friend-service/internal/user-service/domain/user/repo"

	"github.com/samber/mo"
)

type GetFriendsRequest struct {
	UserID string
}

type GetFriendsResponse struct {
	FriendIDs []string
}

type GetFriendsService struct {
	repo repo.FriendsRepo
}

func NewGetFriendsService(repo repo.FriendsRepo) *GetFriendsService {
	return &GetFriendsService{repo: repo}
}

func (s *GetFriendsService) GetFriends(ctx context.Context, req GetFriendsRequest) (GetFriendsResponse, error) {
	if req.UserID == "" {
		return GetFriendsResponse{}, ErrUserIDRequired
	}
	friendEnt, err := s.repo.Get(ctx, entities.FriendsFilter{UserID: mo.Some(req.UserID)})
	if err != nil {
		if err == entities.ErrUsersFriendsNotFound {
			return GetFriendsResponse{}, nil
		}
		return GetFriendsResponse{}, err
	}

	return GetFriendsResponse{FriendIDs: friendEnt.Friends}, nil
}
