package grpcsrv

import (
	"context"
	"errors"

	authgrpc "github.com/AleksandrVishniakov/dc-protos/gen/go/auth/v1"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/servers"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/services/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name=Auth
type Auth interface {
	RegisterNewUser(
		ctx context.Context,
		login string,
		password string,
	) (userID uint64, err error)

	Login(
		ctx context.Context,
		login string,
		password string,
	) (token string, err error)
}

type serverAPI struct {
	authgrpc.UnimplementedAuthServer
	auth Auth
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	authgrpc.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) Register(ctx context.Context, r *authgrpc.RegisterRequest) (*authgrpc.RegisterResponse, error) {
	if r.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, servers.MsgLoginIsRequired)
	}

	if r.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, servers.MsgPasswordIsRequired)
	}

	id, err := s.auth.RegisterNewUser(ctx, r.GetLogin(), r.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, servers.MsgUserAlreadyExists)
		}

		return nil, status.Error(codes.Internal, servers.MsgInternalError)
	}

	return &authgrpc.RegisterResponse{ID: id}, nil
}

func (s *serverAPI) Login(ctx context.Context, r *authgrpc.LoginRequest) (*authgrpc.LoginResponse, error) {
	if r.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, servers.MsgLoginIsRequired)
	}

	if r.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, servers.MsgPasswordIsRequired)
	}

	token, err := s.auth.Login(ctx, r.GetLogin(), r.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, servers.MsgUserNotFound)
		}
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, servers.MsgInvalidCredentials)
		}

		return nil, status.Error(codes.Internal, servers.MsgInternalError)
	}

	return &authgrpc.LoginResponse{Token: token}, nil
}
