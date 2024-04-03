package httpsrv

import (
	"errors"

	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/servers"
)

type UserRequestDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (u *UserRequestDTO) Valid() error {
	if u.Login == "" {
		return errors.New(servers.MsgLoginIsRequired)
	}

	if u.Password == "" {
		return errors.New(servers.MsgPasswordIsRequired)
	}

	if len(u.Password) > 64 {
		return errors.New(servers.MsgTooLongPassword)
	}

	return nil
}

type RegisterResponseDTO struct {
	ID uint64 `json:"id"`
}

type LoginResponseDTO struct {
	Token string `json:"token"`
}
