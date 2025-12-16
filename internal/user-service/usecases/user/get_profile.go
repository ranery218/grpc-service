package user

import (
	"context"
	"friend-service/internal/user-service/domain/user/entities"
	"friend-service/internal/user-service/domain/user/repo"

	"github.com/samber/mo"
)

type GetProfileRequest struct {
	ID string
}

type GetProfileResponse struct {
	ID       string
	Username string
}

type GetProfileService struct {
	repo repo.UserRepo
}

func NewGetProfileService(repo repo.UserRepo) *GetProfileService {
	return &GetProfileService{repo: repo}
}

func (s *GetProfileService) GetProfile(ctx context.Context, req GetProfileRequest) (GetProfileResponse, error) {
	user, err := s.repo.Get(ctx, entities.ProfileFilter{ID: mo.Some(req.ID)})
	if err != nil {
		return GetProfileResponse{}, err
	}

	return GetProfileResponse{ID: user.ID, Username: user.Username}, nil
}
