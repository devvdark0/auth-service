package grpc

import (
	"context"
	authv1 "github.com/devvdark0/auth-service/api/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
import "google.golang.org/grpc"

const emptyValue = 0

type Auth interface {
	Login(ctx context.Context, email, password string) (string, error)
	RegisterNewUser(ctx context.Context, nick, email, password string) (int64, error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *authv1.LoginRequest,
) (*authv1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	//TODO: implement auth service
	token, err := s.auth.Login(ctx, req.Email, req.Password)
	if err != nil {
		//TODO: ...
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authv1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userId, err := s.auth.RegisterNewUser(ctx, req.Nickname, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authv1.RegisterResponse{UserId: userId}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *authv1.IsAdminRequest) (*authv1.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.UserId)
	if err != nil {
		//TODO: ...
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &authv1.IsAdminResponse{IsAdmin: isAdmin}, nil
}

func validateLogin(login *authv1.LoginRequest) error {
	if login.Email == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if login.Password == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	return nil
}

func validateRegister(reg *authv1.RegisterRequest) error {
	if reg.Nickname == "" {
		return status.Error(codes.InvalidArgument, "username is required")
	}
	if reg.Email == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if reg.Password == "" {
		return status.Error(codes.InvalidArgument, "password is rquired")
	}

	return nil
}

func validateIsAdmin(req *authv1.IsAdminRequest) error {
	if req.UserId <= emptyValue {
		return status.Error(codes.InvalidArgument, "invalid userID")
	}
	return nil
}
