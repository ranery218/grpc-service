package user

import (
	"context"
	"friend-service/internal/user-service/domain/user/entities"
	"friend-service/internal/user-service/domain/user/repo"
)

type CreatProfileRequest struct {
	ID       string
	Username string
}

type CreateProfileResponse struct {
	ID       string
	Username string
}

type CreateProfileService struct {
	repo repo.UserRepo
}

func NewCreateProfileService(repo repo.UserRepo) *CreateProfileService {
	return &CreateProfileService{repo: repo}
}

func (s *CreateProfileService) CreateProfile(ctx context.Context, req CreatProfileRequest) (CreateProfileResponse, error) {
	user, err := s.repo.Create(ctx, entities.ProfileAttrs{ID: req.ID, Username: req.Username})
	if err != nil {
		return CreateProfileResponse{}, err
	}

	return CreateProfileResponse{ID: user.ID, Username: user.Username}, nil
}
