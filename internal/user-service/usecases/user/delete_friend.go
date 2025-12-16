package user

import (
	"context"
	"friend-service/internal/user-service/domain/user/entities"
	"friend-service/internal/user-service/domain/user/repo"

	"github.com/samber/mo"
)

type DeleteFriendRequest struct {
	UserID   string
	FriendID string
}

type DeleteFriendService struct {
	repo repo.FriendsRepo
}

func NewDeleteFriendService(repo repo.FriendsRepo) *DeleteFriendService {
	return &DeleteFriendService{repo: repo}
}

func (s *DeleteFriendService) DeleteFriend(ctx context.Context, req DeleteFriendRequest) error {
	if req.UserID == "" {
		return ErrUserIDRequired
	}

	if req.FriendID == "" {
		return ErrFriendIDsEmpty
	}

	friendEnt, err := s.repo.Get(ctx, entities.FriendsFilter{UserID: mo.Some(req.UserID)})
	if err != nil {
		if err == entities.ErrUsersFriendsNotFound {
			return ErrNoneFriends
		}
		return err
	}
	friendsList := friendEnt.Friends

	updatedFriends := []string{}
	for _, fID := range friendsList {
		if fID != req.FriendID {
			updatedFriends = append(updatedFriends, fID)
		}
	}

	_, err = s.repo.Update(ctx, entities.FriendsFilter{UserID: mo.Some(req.UserID)}, entities.FriendsAttrs{
		Friends: updatedFriends,
	})
	if err != nil {
		return err
	}

	return nil
}
