package model

import (
	"github.com/pkg/errors"
	"github.com/xqa/chathub/internal/errs"
)

type Role int

const (
	GENERAL Role = iota
	ADMIN
)

type User struct {
	ID        uint   `json:"id" gorm:"primaryKey"`                      // unique key
	Username  string `json:"username" gorm:"unique" binding:"required"` // username
	Password  string `json:"password"`                                  // password
	Role      Role   `json:"role"`                                      // user's role
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatarUrl"`
}

func (u User) IsGeneral() bool {
	return u.Role == GENERAL
}

func (u User) IsAdmin() bool {
	return u.Role == ADMIN
}

func (u User) ValidatePassword(password string) error {
	if password == "" {
		return errors.WithStack(errs.EmptyPassword)
	}
	if u.Password != password {
		return errors.WithStack(errs.WrongPassword)
	}
	return nil
}
