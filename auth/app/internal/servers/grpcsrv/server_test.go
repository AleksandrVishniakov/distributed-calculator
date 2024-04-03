package grpcsrv

import (
	"context"
	"errors"
	"testing"

	authv1 "github.com/AleksandrVishniakov/dc-protos/gen/go/auth/v1"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/servers"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/servers/grpcsrv/mocks"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/services/auth"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/pkg/fake"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestServerAPI_Register(t *testing.T) {
	type fields struct {
		ctx         context.Context
		authService Auth
		req         *authv1.RegisterRequest
	}

	var tests = []struct {
		name    string
		fields  *fields
		prepare func(f *fields)
		wantErr error
	}{
		{
			name: "test_register",
			fields: &fields{
				ctx: context.Background(),
				req: &authv1.RegisterRequest{
					Login:    gofakeit.Name(),
					Password: fake.RandPassword(),
				},
			},
			prepare: func(f *fields) {
				authService := mocks.NewAuth(t)
				authService.
					On("RegisterNewUser", mock.Anything, f.req.Login, f.req.Password).
					Once().
					Return(uint64(1), nil)

				f.authService = authService
			},
		},
		{
			name: "empty_login",
			fields: &fields{
				ctx: context.Background(),
				req: &authv1.RegisterRequest{
					Login:    "",
					Password: fake.RandPassword(),
				},
			},
			wantErr: errors.New(servers.MsgLoginIsRequired),
		},
		{
			name: "empty_password",
			fields: &fields{
				ctx: context.Background(),
				req: &authv1.RegisterRequest{
					Login:    gofakeit.Name(),
					Password: "",
				},
			},
			wantErr: errors.New(servers.MsgPasswordIsRequired),
		},
		{
			name: "user_already_exists",
			fields: &fields{
				ctx: context.Background(),
				req: &authv1.RegisterRequest{
					Login:    gofakeit.Name(),
					Password: fake.RandPassword(),
				},
			},
			prepare: func(f *fields) {
				authService := mocks.NewAuth(t)

				authService.
					On("RegisterNewUser", mock.Anything, f.req.Login, f.req.Password).
					Once().
					Return(uint64(0), auth.ErrUserAlreadyExists)

				f.authService = authService
			},
			wantErr: errors.New(servers.MsgUserAlreadyExists),
		},
		{
			name: "unexpected_internal_error",
			fields: &fields{
				ctx: context.Background(),
				req: &authv1.RegisterRequest{
					Login:    gofakeit.Name(),
					Password: fake.RandPassword(),
				},
			},
			prepare: func(f *fields) {
				authService := mocks.NewAuth(t)

				authService.
					On("RegisterNewUser", mock.Anything, f.req.Login, f.req.Password).
					Once().
					Return(uint64(0), errors.New("unexpected error"))

				f.authService = authService
			},
			wantErr: errors.New(servers.MsgInternalError),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.prepare != nil {
				test.prepare(test.fields)
			}

			server := &serverAPI{auth: test.fields.authService}

			_, err := server.Register(test.fields.ctx, test.fields.req)
			if err != nil {
				if test.wantErr == nil {
					t.Fatalf("unexpected error: %s", err.Error())
				}

				require.ErrorContains(t, err, test.wantErr.Error(), "")
			}
		})
	}
}

func TestServerAPI_Login(t *testing.T) {
	type fields struct {
		ctx         context.Context
		authService Auth
		req         *authv1.LoginRequest
	}

	var tests = []struct {
		name    string
		fields  *fields
		prepare func(f *fields)
		wantErr error
	}{
		{
			name: "test_login",
			fields: &fields{
				ctx: context.Background(),
				req: &authv1.LoginRequest{
					Login:    gofakeit.Name(),
					Password: fake.RandPassword(),
				},
			},
			prepare: func(f *fields) {
				authService := mocks.NewAuth(t)
				authService.
					On("Login", mock.Anything, f.req.Login, f.req.Password).
					Once().
					Return(gofakeit.UUID(), nil)

				f.authService = authService
			},
		},
		{
			name: "empty_login",
			fields: &fields{
				ctx: context.Background(),
				req: &authv1.LoginRequest{
					Login:    "",
					Password: fake.RandPassword(),
				},
			},
			wantErr: errors.New(servers.MsgLoginIsRequired),
		},
		{
			name: "empty_password",
			fields: &fields{
				ctx: context.Background(),
				req: &authv1.LoginRequest{
					Login:    gofakeit.Name(),
					Password: "",
				},
			},
			wantErr: errors.New(servers.MsgPasswordIsRequired),
		},
		{
			name: "user_not_found",
			fields: &fields{
				ctx: context.Background(),
				req: &authv1.LoginRequest{
					Login:    gofakeit.Name(),
					Password: fake.RandPassword(),
				},
			},
			prepare: func(f *fields) {
				authService := mocks.NewAuth(t)

				authService.
					On("Login", mock.Anything, f.req.Login, f.req.Password).
					Once().
					Return("", auth.ErrUserNotFound)

				f.authService = authService
			},
			wantErr: errors.New(servers.MsgUserNotFound),
		},
		{
			name: "invalid_credentials",
			fields: &fields{
				ctx: context.Background(),
				req: &authv1.LoginRequest{
					Login:    gofakeit.Name(),
					Password: fake.RandPassword(),
				},
			},
			prepare: func(f *fields) {
				authService := mocks.NewAuth(t)

				authService.
					On("Login", mock.Anything, f.req.Login, f.req.Password).
					Once().
					Return("", auth.ErrInvalidCredentials)

				f.authService = authService
			},
			wantErr: errors.New(servers.MsgInvalidCredentials),
		},
		{
			name: "unexpected_internal_error",
			fields: &fields{
				ctx: context.Background(),
				req: &authv1.LoginRequest{
					Login:    gofakeit.Name(),
					Password: fake.RandPassword(),
				},
			},
			prepare: func(f *fields) {
				authService := mocks.NewAuth(t)

				authService.
					On("Login", mock.Anything, f.req.Login, f.req.Password).
					Once().
					Return("", errors.New("unexpected error"))

				f.authService = authService
			},
			wantErr: errors.New(servers.MsgInternalError),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.prepare != nil {
				test.prepare(test.fields)
			}

			server := &serverAPI{auth: test.fields.authService}

			_, err := server.Login(test.fields.ctx, test.fields.req)
			if err != nil {
				if test.wantErr == nil {
					t.Fatalf("unexpected error: %s", err.Error())
				}

				require.ErrorContains(t, err, test.wantErr.Error(), "")
			}
		})
	}
}
