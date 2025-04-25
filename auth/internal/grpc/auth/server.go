package auth

import (
	"context"
	"errors"

	auth "gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/auth/internal/services"
	authpb "gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/auth/pkg/pb/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Auth interface {
	Register(ctx context.Context, telegramLogin, login, password string) (string, error)
	Login(ctx context.Context, telegramLogin, login, password string) (string, error)
	IsLogged(ctx context.Context, telegramLogin string) (string, error)
	Logout(ctx context.Context, telegramLogin string) error
}

type ServerAPI struct {
	authpb.UnimplementedAuth_V1Server
	auth Auth
}

var (
	ErrUserExists         = "user already exists"
	ErrUserNotFound       = "user not found"
	ErrInternal           = "internal error"
	ErrInvalidCredentials = "invalid credentials"
	ErrTokenNotFound      = "token not found"
	ErrTokenExists        = "token already exists"
)

func Register(grpc *grpc.Server, auth Auth) {
	authpb.RegisterAuth_V1Server(grpc, &ServerAPI{auth: auth})
}

func (s *ServerAPI) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.Register(ctx, req.GetTelegramLogin(), req.GetLogin(), req.GetPassword())
	if err != nil {
		return nil, formatError(err)
	}

	return &authpb.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *ServerAPI) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetTelegramLogin(), req.GetLogin(), req.GetPassword())
	if err != nil {
		return nil, formatError(err)
	}

	return &authpb.LoginResponse{
		Token: token,
	}, nil
}

func (s *ServerAPI) IsLogged(ctx context.Context, req *authpb.IsLoggedRequest) (*authpb.IsLoggedResponse, error) {
	if err := validateIsLogged(req); err != nil {
		return nil, err
	}
	token, err := s.auth.IsLogged(ctx, req.GetTelegramLogin())
	if err != nil {
		return nil, formatError(err)
	}

	return &authpb.IsLoggedResponse{
		Token: token,
	}, nil
}

func (s *ServerAPI) Logout(ctx context.Context, req *authpb.LogoutRequest) (*emptypb.Empty, error) {
	if err := validateLogout(req); err != nil {
		return nil, err
	}

	if err := s.auth.Logout(ctx, req.GetTelegramLogin()); err != nil {
		return nil, formatError(err)
	}

	return &emptypb.Empty{}, nil
}

func validateRegister(req *authpb.RegisterRequest) error {
	if req.GetLogin() == "" {
		return status.Error(codes.InvalidArgument, "login is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	if req.GetTelegramLogin() == "" {
		return status.Error(codes.InvalidArgument, "telegram login is required")
	}
	return nil
}

func validateLogin(req *authpb.LoginRequest) error {
	if req.GetLogin() == "" {
		return status.Error(codes.InvalidArgument, "login is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	if req.GetTelegramLogin() == "" {
		return status.Error(codes.InvalidArgument, "telegram login is required")
	}
	return nil
}

func validateIsLogged(req *authpb.IsLoggedRequest) error {
	if req.GetTelegramLogin() == "" {
		return status.Error(codes.InvalidArgument, "telegram login is required")
	}

	return nil
}

func validateLogout(req *authpb.LogoutRequest) error {
	if req.GetTelegramLogin() == "" {
		return status.Error(codes.InvalidArgument, "telegram login is required")
	}

	return nil
}

func formatError(err error) error {
	if errors.Is(err, auth.ErrUserExists) {
		return status.Error(codes.AlreadyExists, ErrUserExists)
	} else if errors.Is(err, auth.ErrUserNotFound) {
		return status.Error(codes.NotFound, ErrUserNotFound)
	} else if errors.Is(err, auth.ErrInvalidCredentials) {
		return status.Error(codes.PermissionDenied, ErrInvalidCredentials)
	} else if errors.Is(err, auth.ErrTokenNotFound) {
		return status.Error(codes.NotFound, ErrTokenNotFound)
	} else if errors.Is(err, auth.ErrTokenExists) {
		return status.Error(codes.AlreadyExists, ErrTokenExists)
	}

	return status.Error(codes.Internal, ErrInternal)
}
