package auth

import (
	"context"
	"errors"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/repositories/usersrepo"
	"log/slog"

	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/domain/models"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/pkg/e"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/pkg/sl"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("auth: invalid credentials")
	ErrUserAlreadyExists  = errors.New("auth: user already exists")
	ErrUserNotFound       = errors.New("auth: user not found")
)

type Auth struct {
	log            *slog.Logger
	userSaver      UserSaver
	userProvider   UserProvider
	tokenGenerator TokenGenerator
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name=UserSaver
type UserSaver interface {
	Save(
		ctx context.Context,
		login string,
		passwordHash string,
	) (id uint64, err error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name=UserProvider
type UserProvider interface {
	User(
		ctx context.Context,
		login string,
	) (user *models.User, err error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name=TokenGenerator
type TokenGenerator interface {
	NewToken(
		userId uint64,
		login string,
	) (token string, err error)
}

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	tokenGenerator TokenGenerator,
) *Auth {
	return &Auth{
		log:            log,
		userSaver:      userSaver,
		userProvider:   userProvider,
		tokenGenerator: tokenGenerator,
	}
}

// RegisterNewUser registers new user and returns userID
//
// Returns ErrUserAlreadyExists, if user with provided login already exists in storage
func (a *Auth) RegisterNewUser(ctx context.Context, login string, password string) (userID uint64, err error) {
	const src = "Auth.RegisterNewUser"

	log := a.log.With(
		slog.String("src", src),
		slog.String("login", login),
	)

	log.Debug("registering new user...")

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))

		return 0, e.WrapErr(err, src)
	}

	id, err := a.userSaver.Save(ctx, login, string(passwordHash))
	if err != nil {
		if errors.Is(err, usersrepo.ErrUserAlreadyExists) {
			return 0, e.WrapErr(ErrUserAlreadyExists, src)
		}

		return 0, e.WrapErr(err, src)
	}

	log.Debug("user registered")

	return id, nil
}

// Login logs in an existing user and returns auth token
//
// # Returns ErrUserNotFound, if user with provided login not found
//
// Returns ErrInvalidCredentials, if stored password hash not equal to hash of provided password
func (a *Auth) Login(ctx context.Context, login string, password string) (token string, err error) {
	const src = "Auth.Login"

	log := a.log.With(
		slog.String("src", src),
		slog.String("login", login),
	)

	log.Debug("logging in user...")

	user, err := a.userProvider.User(ctx, login)
	if err != nil {
		if errors.Is(err, usersrepo.ErrUserNotFound) {
			return "", e.WrapErr(ErrUserNotFound, src)
		}

		return "", e.WrapErr(err, src)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		log.Error("failed to validate user password", sl.Err(err))
		return "", e.WrapErr(ErrInvalidCredentials, src)
	}

	token, err = a.tokenGenerator.NewToken(
		user.ID,
		user.Login,
	)

	if err != nil {
		log.Error("failed to generate token", sl.Err(err))

		return "", e.WrapErr(err, src)
	}

	log.Debug("user logged in")

	return token, nil
}
