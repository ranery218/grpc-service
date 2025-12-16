package user

import (
	"context"
	"friend-service/internal/user-service/domain/user/entities"
	"friend-service/internal/user-service/domain/user/repo"
)

type GetAllProfilesRequest struct{}

type GetAllProfilesResponse struct {
	Profiles []struct{ ID, Username string }
}

type GetAllProfilesService struct {
	repo repo.UserRepo
}

func NewGetAllProfilesService(repo repo.UserRepo) *GetAllProfilesService {
	return &GetAllProfilesService{repo: repo}
}

func (s *GetAllProfilesService) GetAllProfiles(ctx context.Context, req GetAllProfilesRequest) (GetAllProfilesResponse, error) {
	users, err := s.repo.GetAll(ctx, entities.ProfileFilter{})
	if err != nil {
		return GetAllProfilesResponse{}, err
	}

	var resp GetAllProfilesResponse
	for _, user := range users {
		resp.Profiles = append(resp.Profiles, struct{ ID, Username string }{ID: user.ID, Username: user.Username})
	}

	return resp, nil
}
