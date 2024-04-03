package auth

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/domain/models"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/repositories/usersrepo"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/services/auth/mocks"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/pkg/fake"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/pkg/sl"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

const defaultUserID uint64 = 1

func TestAuth_Login(t *testing.T) {
	type fields struct {
		log          *slog.Logger
		userSaver    UserSaver
		userProvider UserProvider
		tokenParser  TokenGenerator
	}
	type args struct {
		ctx      context.Context
		login    string
		password string
	}

	tests := []struct {
		name    string
		fields  *fields
		prepare func(t *testing.T, a *args, f *fields, wantToken string)
		args    *args

		wantErr   bool
		targetErr error

		wantToken string
	}{
		{
			name: "test_login",
			args: &args{
				login:    gofakeit.Name(),
				password: fake.RandPassword(),
			},
			fields: &fields{
				log: sl.NewDiscardLogger(),
			},
			prepare: func(t *testing.T, a *args, f *fields, wantToken string) {
				userProvider := mocks.NewUserProvider(t)
				tokenParser := mocks.NewTokenGenerator(t)

				passHash, err := bcrypt.GenerateFromPassword([]byte(a.password), bcrypt.DefaultCost)
				if err != nil {
					t.Fatalf("password hash generating error: %s", err.Error())
				}
				userProvider.
					On("User", mock.Anything, a.login).
					Once().
					Return(&models.User{
						ID:           defaultUserID,
						Login:        a.login,
						PasswordHash: string(passHash),
					}, nil)

				tokenParser.
					On("NewToken", defaultUserID, a.login).
					Once().
					Return(wantToken, nil)

				f.userProvider = userProvider
				f.tokenParser = tokenParser
			},
			wantToken: gofakeit.UUID(),
		},
		{
			name: "err_user_not_found",
			args: &args{
				login:    gofakeit.Name(),
				password: fake.RandPassword(),
			},
			fields: &fields{
				log: sl.NewDiscardLogger(),
			},
			prepare: func(t *testing.T, a *args, f *fields, wantToken string) {
				userProvider := mocks.NewUserProvider(t)
				userProvider.
					On("User", mock.Anything, a.login).
					Once().
					Return(nil, usersrepo.ErrUserNotFound)

				f.userProvider = userProvider
			},
			wantErr:   true,
			targetErr: ErrUserNotFound,
		},
		{
			name: "err_get_user_unexpected_error",
			args: &args{
				login:    gofakeit.Name(),
				password: fake.RandPassword(),
			},
			fields: &fields{
				log: sl.NewDiscardLogger(),
			},
			prepare: func(t *testing.T, a *args, f *fields, wantToken string) {
				userProvider := mocks.NewUserProvider(t)
				userProvider.
					On("User", mock.Anything, a.login).
					Once().
					Return(nil, errors.New("unexpected error"))

				f.userProvider = userProvider
			},
			wantErr: true,
		},
		{
			name: "err_password_validation_fail",
			args: &args{
				login:    gofakeit.Name(),
				password: fake.RandPassword(),
			},
			fields: &fields{
				log: sl.NewDiscardLogger(),
			},
			prepare: func(t *testing.T, a *args, f *fields, wantToken string) {
				userProvider := mocks.NewUserProvider(t)
				userProvider.
					On("User", mock.Anything, a.login).
					Once().
					Return(&models.User{
						ID:           defaultUserID,
						Login:        a.login,
						PasswordHash: "incorrect_password_hash",
					}, nil)

				f.userProvider = userProvider
			},
			wantErr:   true,
			targetErr: ErrInvalidCredentials,
		},
		{
			name: "err_token_generation_error",
			args: &args{
				login:    gofakeit.Name(),
				password: fake.RandPassword(),
			},
			fields: &fields{
				log: sl.NewDiscardLogger(),
			},
			prepare: func(t *testing.T, a *args, f *fields, wantToken string) {
				userProvider := mocks.NewUserProvider(t)
				tokenParser := mocks.NewTokenGenerator(t)

				passHash, err := bcrypt.GenerateFromPassword([]byte(a.password), bcrypt.DefaultCost)
				if err != nil {
					t.Fatalf("password hash generating error: %s", err.Error())
				}
				userProvider.
					On("User", mock.Anything, a.login).
					Once().
					Return(&models.User{
						ID:           defaultUserID,
						Login:        a.login,
						PasswordHash: string(passHash),
					}, nil)

				tokenParser.
					On("NewToken", defaultUserID, a.login).
					Once().
					Return("", errors.New("unexpected error"))

				f.userProvider = userProvider
				f.tokenParser = tokenParser
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.ctx == nil {
				tt.args.ctx = context.Background()
			}

			if tt.prepare != nil {
				tt.prepare(t, tt.args, tt.fields, tt.wantToken)
			}

			a := &Auth{
				log:            tt.fields.log,
				userSaver:      tt.fields.userSaver,
				userProvider:   tt.fields.userProvider,
				tokenGenerator: tt.fields.tokenParser,
			}
			token, err := a.Login(tt.args.ctx, tt.args.login, tt.args.password)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("unexpected error: %s", err.Error())
				}

				if tt.targetErr != nil {
					assert.ErrorIs(t, err, tt.targetErr)
				}
			}

			assert.Equal(t, tt.wantToken, token)
		})
	}
}

func TestAuth_RegisterNewUser(t *testing.T) {
	type fields struct {
		log          *slog.Logger
		userSaver    UserSaver
		userProvider UserProvider
		tokenParser  TokenGenerator
	}
	type args struct {
		ctx      context.Context
		login    string
		password string
	}

	tests := []struct {
		name    string
		fields  *fields
		prepare func(t *testing.T, a *args, f *fields, wantUserID uint64)
		args    *args

		wantErr   bool
		targetErr error

		wantUserID uint64
	}{
		{
			name: "test_register",
			fields: &fields{
				log: sl.NewDiscardLogger(),
			},
			args: &args{
				login:    gofakeit.Name(),
				password: fake.RandPassword(),
			},
			prepare: func(t *testing.T, a *args, f *fields, wantUserID uint64) {
				userSaver := mocks.NewUserSaver(t)
				userSaver.
					On("Save", mock.Anything, a.login, mock.AnythingOfType("string")).
					Once().
					Return(wantUserID, nil)

				f.userSaver = userSaver
			},

			wantUserID: gofakeit.Uint64(),
		},
		{
			name: "err_user_already_exists",
			fields: &fields{
				log: sl.NewDiscardLogger(),
			},
			args: &args{
				login:    gofakeit.Name(),
				password: fake.RandPassword(),
			},
			prepare: func(t *testing.T, a *args, f *fields, wantUserID uint64) {
				userSaver := mocks.NewUserSaver(t)
				userSaver.
					On("Save", mock.Anything, a.login, mock.AnythingOfType("string")).
					Once().
					Return(uint64(0), usersrepo.ErrUserAlreadyExists)

				f.userSaver = userSaver
			},

			wantErr:   true,
			targetErr: ErrUserAlreadyExists,
		},
		{
			name: "err_user_creation_unexpected_error",
			fields: &fields{
				log: sl.NewDiscardLogger(),
			},
			args: &args{
				login:    gofakeit.Name(),
				password: fake.RandPassword(),
			},
			prepare: func(t *testing.T, a *args, f *fields, wantUserID uint64) {
				userSaver := mocks.NewUserSaver(t)
				userSaver.
					On("Save", mock.Anything, a.login, mock.AnythingOfType("string")).
					Once().
					Return(uint64(0), errors.New("unexpected error"))

				f.userSaver = userSaver
			},

			wantErr: true,
		},
		{
			name: "err_pass_hash_generation_unexpected_error",
			fields: &fields{
				log: sl.NewDiscardLogger(),
			},
			args: &args{
				login: gofakeit.Name(),
				// if password is longer than 72, bcrypt.GenerateFromPassword returns bcrypt.ErrPasswordTooLong
				password: gofakeit.Password(true, true, true, true, true, 100),
			},

			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.ctx == nil {
				tt.args.ctx = context.Background()
			}

			if tt.prepare != nil {
				tt.prepare(t, tt.args, tt.fields, tt.wantUserID)
			}

			a := &Auth{
				log:            tt.fields.log,
				userSaver:      tt.fields.userSaver,
				userProvider:   tt.fields.userProvider,
				tokenGenerator: tt.fields.tokenParser,
			}

			userID, err := a.RegisterNewUser(tt.args.ctx, tt.args.login, tt.args.password)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("unexpected error: %s", err.Error())
				}

				if tt.targetErr != nil {
					assert.ErrorIs(t, err, tt.targetErr)
				}
			}

			assert.Equal(t, tt.wantUserID, userID)
		})
	}
}

func TestNew(t *testing.T) {
	authService := New(sl.NewDiscardLogger(), nil, nil, nil)
	require.NotNil(t, authService, "cannot get nil from Auth.New constructor")
}
