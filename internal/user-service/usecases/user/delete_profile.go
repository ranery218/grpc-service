package user

import (
	"context"
	"friend-service/internal/user-service/domain/user/entities"
	"friend-service/internal/user-service/domain/user/repo"

	"github.com/samber/mo"
)

type DeleteProfileRequest struct {
	ID string
}

type DeleteProfileService struct {
	repo repo.UserRepo
}

func NewDeleteProfileService(repo repo.UserRepo) *DeleteProfileService {
	return &DeleteProfileService{repo: repo}
}

func (s *DeleteProfileService) DeleteProfile(ctx context.Context, req DeleteProfileRequest) error {
	err := s.repo.Delete(ctx, entities.ProfileFilter{ID: mo.Some(req.ID)})
	if err != nil {
		return err
	}

	return nil
}
