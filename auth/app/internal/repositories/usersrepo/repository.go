package usersrepo

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/domain/models"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/pkg/e"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/pkg/sl"
)

var (
	ErrUserNotFound      = errors.New("usersrepo: user not found")
	ErrUserAlreadyExists = errors.New("usersrepo: users already exists")
)

type UserRepository struct {
	log *slog.Logger
	db  *sql.DB
}

func New(
	log *slog.Logger,
	db *sql.DB,
) *UserRepository {
	return &UserRepository{
		log: log,
		db:  db,
	}
}

func (u *UserRepository) Save(_ context.Context, login string, passwordHash string) (id uint64, err error) {
	const src = "UserRepository.save"

	log := u.log.With(
		slog.String("src", src),
		slog.String("login", login),
	)

	tx, err := u.db.Begin()
	if err != nil {
		log.Error("failed to begin transaction", sl.Err(err))
	}

	row := tx.QueryRow(
		`SELECT u.* FROM users u
				WHERE login=$1`,
		login,
	)

	var user = &models.User{}
	err = row.Scan(&user.ID, &user.Login, &user.PasswordHash)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Error("failed to get user", sl.Err(err))
		err := tx.Rollback()
		if err != nil {
			log.Error("failed to rollback", sl.Err(err))
			return 0, err
		}
		return 0, err
	}

	if err == nil {
		log.Warn("user already exists")
		err := tx.Rollback()
		if err != nil {
			log.Error("failed to rollback", sl.Err(err))
			return 0, err
		}
		return 0, ErrUserAlreadyExists
	}

	row = tx.QueryRow(
		`INSERT INTO users (login, password_hash)
				VALUES ($1, $2)
				RETURNING id`,
		login,
		passwordHash,
	)

	id = 0
	err = row.Scan(&id)

	if err != nil {
		log.Error("failed to create user", sl.Err(err))
		err := tx.Rollback()
		if err != nil {
			log.Error("failed to rollback", sl.Err(err))
			return 0, err
		}
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		log.Error("failed to commit", sl.Err(err))
		return 0, err
	}

	log.Debug("user saved")

	return id, nil
}

func (u *UserRepository) User(_ context.Context, login string) (user *models.User, err error) {
	const src = "UserRepository.User"

	log := u.log.With(
		slog.String("src", src),
		slog.String("login", login),
	)

	row := u.db.QueryRow(
		`SELECT u.* FROM users u
				WHERE login=$1`,
		login,
	)

	user = &models.User{}
	err = row.Scan(&user.ID, &user.Login, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("user not found")
			return nil, e.WrapErr(ErrUserNotFound, src)
		}
		log.Error("failed to get user", sl.Err(err))
		return nil, e.WrapErr(err, src)
	}

	log.Debug("user got")

	return user, nil
}
