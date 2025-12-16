package user

import (
	"context"
	"friend-service/internal/user-service/domain/user/entities"
	"friend-service/internal/user-service/domain/user/repo"

	"github.com/samber/mo"
)

type AddFriendRequest struct {
	UserID   string
	FriendID string
}

type AddFriendResponse struct {
	FriendIDs []string
}

type AddFriendService struct {
	repo repo.FriendsRepo
}

func NewAddFriendService(repo repo.FriendsRepo) *AddFriendService {
	return &AddFriendService{repo: repo}
}

func (s *AddFriendService) AddFriend(ctx context.Context, req AddFriendRequest) (AddFriendResponse, error) {
	if req.UserID == "" {
		return AddFriendResponse{}, ErrUserIDRequired
	}

	if req.FriendID == "" {
		return AddFriendResponse{}, ErrFriendIDsEmpty
	}

	friendEnt, err := s.repo.Get(ctx, entities.FriendsFilter{UserID: mo.Some(req.UserID)})
	if err != nil {
		if err == entities.ErrUsersFriendsNotFound {
			_, err = s.repo.Create(ctx, entities.FriendsAttrs{
				UserID:  req.UserID,
				Friends: []string{req.FriendID},
			})
			if err != nil {
				return AddFriendResponse{}, err
			}
			return AddFriendResponse{FriendIDs: []string{req.FriendID}}, nil
		}
		return AddFriendResponse{}, err
	}

	friendList := friendEnt.Friends

	friendsMap := make(map[string]struct{})
	for _, fID := range friendList {
		friendsMap[fID] = struct{}{}
	}

	if _, exists := friendsMap[req.FriendID]; !exists {
		friendList = append(friendList, req.FriendID)
	}

	updatedFriendEnt, err := s.repo.Update(ctx, entities.FriendsFilter{UserID: mo.Some(req.UserID)}, entities.FriendsAttrs{
		UserID:  req.UserID,
		Friends: friendList,
	})
	if err != nil {
		return AddFriendResponse{}, err
	}

	return AddFriendResponse{FriendIDs: updatedFriendEnt.Friends}, nil
}
