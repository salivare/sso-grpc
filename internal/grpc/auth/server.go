package auth

import (
	"context"
	ssov1 "github.com/salivare/protos-sso/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appId int32,
	) (token string, err error)
	Register(
		ctx context.Context,
		email string,
		password string,
	) (userId int64, err error)
	IsAdmin(
		ctx context.Context,
		userId int64,
	) (bool, error)
}

type ServerAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &ServerAPI{auth: auth})
}

const (
	emptyValue = 0
)

func (s *ServerAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	// TODO: Add middleware
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		// TODO:...
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *ServerAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	// TODO: Add middleware
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		// TODO: ...
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &ssov1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *ServerAPI) IsAdmin(
	ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	// TODO: Add middleware
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func validateLogin(req *ssov1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is empty")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is empty")
	}

	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "appId is empty")
	}

	return nil
}

func validateRegister(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is empty")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is empty")
	}

	return nil
}

func validateIsAdmin(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == 0 {
		return status.Error(codes.InvalidArgument, "userId is empty")
	}

	return nil
}
