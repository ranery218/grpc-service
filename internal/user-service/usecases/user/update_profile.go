package user

import (
	"context"
	"friend-service/internal/user-service/domain/user/entities"
	"friend-service/internal/user-service/domain/user/repo"

	"github.com/samber/mo"
)

type UpdateProfileRequest struct {
	ID       string
	Username string
}

type UpdateProfileResponse struct {
	ID       string
	Username string
}

type UpdateProfileService struct {
	repo repo.UserRepo
}

func NewUpdateProfileService(repo repo.UserRepo) *UpdateProfileService {
	return &UpdateProfileService{repo: repo}
}

func (s *UpdateProfileService) UpdateProfile(ctx context.Context, req UpdateProfileRequest) (UpdateProfileResponse, error) {
	user, err := s.repo.Update(ctx, entities.ProfileFilter{ID: mo.Some(req.ID)}, entities.ProfileAttrs{ID: req.ID, Username: req.Username})
	if err != nil {
		return UpdateProfileResponse{}, err
	}

	return UpdateProfileResponse{ID: user.ID, Username: user.Username}, nil
}
