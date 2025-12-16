package transport

import (
	"context"
	"errors"
	"friend-service/internal/user-service/usecases/user"
	userv1 "friend-service/proto/gen/user/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	userv1.UnimplementedUserServiceServer
	createProfileService  *user.CreateProfileService
	getProfileService     *user.GetProfileService
	updateProfileService  *user.UpdateProfileService
	deleteProfileService  *user.DeleteProfileService
	getAllProfilesService *user.GetAllProfilesService
	addFriendService      *user.AddFriendService
	getFriendsService     *user.GetFriendsService
	deleteFriendService   *user.DeleteFriendService
}

func NewUserServer(
	createProfileService *user.CreateProfileService,
	getProfileService *user.GetProfileService,
	updateProfileService *user.UpdateProfileService,
	deleteProfileService *user.DeleteProfileService,
	getAllProfilesService *user.GetAllProfilesService,
	addFriendService *user.AddFriendService,
	getFriendsService *user.GetFriendsService,
	deleteFriendService *user.DeleteFriendService,
) *UserServer {
	return &UserServer{
		createProfileService:  createProfileService,
		getProfileService:     getProfileService,
		updateProfileService:  updateProfileService,
		deleteProfileService:  deleteProfileService,
		getAllProfilesService: getAllProfilesService,
		addFriendService:      addFriendService,
		getFriendsService:     getFriendsService,
		deleteFriendService:   deleteFriendService,
	}
}

func (s *UserServer) CreateProfile(ctx context.Context, req *userv1.CreateProfileRequest) (*userv1.CreateProfileResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	userID := req.GetUserId()
	username := req.GetUsername()

	if userID == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}
	if username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}

	serviceReq := user.CreatProfileRequest{
		ID:       userID,
		Username: username,
	}

	serviceResp, err := s.createProfileService.CreateProfile(ctx, serviceReq)
	if err != nil {
		return nil, mapErr(err)
	}

	return &userv1.CreateProfileResponse{
		UserId:   serviceResp.ID,
		Username: serviceResp.Username,
	}, nil
}

func (s *UserServer) GetProfile(ctx context.Context, req *userv1.GetProfileRequest) (*userv1.GetProfileResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	userID := req.GetUserId()
	if userID == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	serviceReq := user.GetProfileRequest{
		ID: userID,
	}

	serviceResp, err := s.getProfileService.GetProfile(ctx, serviceReq)
	if err != nil {
		return nil, mapErr(err)
	}

	return &userv1.GetProfileResponse{
		UserId:   serviceResp.ID,
		Username: serviceResp.Username,
	}, nil
}

func (s *UserServer) UpdateProfile(ctx context.Context, req *userv1.UpdateProfileRequest) (*userv1.UpdateProfileResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	userID := req.GetUserId()
	username := req.GetUsername()

	if userID == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}
	if username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}

	serviceReq := user.UpdateProfileRequest{
		ID:       userID,
		Username: username,
	}

	serviceResp, err := s.updateProfileService.UpdateProfile(ctx, serviceReq)
	if err != nil {
		return nil, mapErr(err)
	}

	return &userv1.UpdateProfileResponse{
		UserId:   serviceResp.ID,
		Username: serviceResp.Username,
	}, nil
}

func (s *UserServer) DeleteProfile(ctx context.Context, req *userv1.DeleteProfileRequest) (*userv1.DeleteProfileResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	userID := req.GetUserId()
	if userID == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	serviceReq := user.DeleteProfileRequest{
		ID: userID,
	}

	err := s.deleteProfileService.DeleteProfile(ctx, serviceReq)
	if err != nil {
		return nil, mapErr(err)
	}

	return &userv1.DeleteProfileResponse{}, nil
}

func (s *UserServer) GetProfileList(ctx context.Context, req *userv1.GetProfileListRequest) (*userv1.GetProfileListResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	serviceResp, err := s.getAllProfilesService.GetAllProfiles(ctx, user.GetAllProfilesRequest{})
	if err != nil {
		return nil, mapErr(err)
	}

	var profiles []*userv1.Profile
	for _, p := range serviceResp.Profiles {
		profiles = append(profiles, &userv1.Profile{
			UserId:   p.ID,
			Username: p.Username,
		})
	}

	return &userv1.GetProfileListResponse{
		Profiles: profiles,
	}, nil
}

func (s *UserServer) AddFriend(ctx context.Context, req *userv1.AddFriendRequest) (*userv1.AddFriendResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	userID := req.GetUserId()
	friendID := req.GetFriendId()

	if userID == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}
	if friendID == "" {
		return nil, status.Error(codes.InvalidArgument, "friend ID is required")
	}

	serviceReq := user.AddFriendRequest{
		UserID:   userID,
		FriendID: friendID,
	}

	serviceResp, err := s.addFriendService.AddFriend(ctx, serviceReq)
	if err != nil {
		return nil, mapErr(err)
	}

	friendIDs := []string{}
	for _, f := range serviceResp.FriendIDs {
		friendIDs = append(friendIDs, f)
	}

	return &userv1.AddFriendResponse{FriendIds: friendIDs}, nil
}

func (s *UserServer) GetFriends(ctx context.Context, req *userv1.GetFriendsRequest) (*userv1.GetFriendsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	userID := req.GetUserId()
	if userID == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	serviceReq := user.GetFriendsRequest{
		UserID: userID,
	}

	serviceResp, err := s.getFriendsService.GetFriends(ctx, serviceReq)
	if err != nil {
		return nil, mapErr(err)
	}

	friendIDs := []string{}
	for _, f := range serviceResp.FriendIDs {
		friendIDs = append(friendIDs, f)
	}

	return &userv1.GetFriendsResponse{FriendIds: friendIDs}, nil
}

func (s *UserServer) DeleteFriend(ctx context.Context, req *userv1.DeleteFriendRequest) (*userv1.DeleteFriendResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	userID := req.GetUserId()
	friendID := req.GetFriendId()

	if userID == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}
	if friendID == "" {
		return nil, status.Error(codes.InvalidArgument, "friend ID is required")
	}

	serviceReq := user.DeleteFriendRequest{
		UserID:   userID,
		FriendID: friendID,
	}

	err := s.deleteFriendService.DeleteFriend(ctx, serviceReq)
	if err != nil {
		return nil, mapErr(err)
	}

	return &userv1.DeleteFriendResponse{}, nil
}

func mapErr(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, user.ErrUserIDRequired):
		return status.Error(codes.Unauthenticated, "unauthorized: user ID is required")
	case errors.Is(err, user.ErrNoneFriends):
		return status.Error(codes.Unauthenticated, "friends not found")
	case errors.Is(err, user.ErrFriendIDsEmpty):
		return status.Error(codes.Unauthenticated, "friend IDs list cannot be empty")
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
