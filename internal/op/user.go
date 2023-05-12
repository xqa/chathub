package op

import (
	"time"

	"github.com/Xhofe/go-cache"
	"github.com/xqa/chathub/internal/db"
	"github.com/xqa/chathub/internal/errs"
	"github.com/xqa/chathub/internal/model"
	"github.com/xqa/chathub/pkg/singleflight"
)

var userCache = cache.NewMemCache(cache.WithShards[*model.User](2))
var userG singleflight.Group[*model.User]
var guestUser *model.User
var adminUser *model.User

func GetAdmin() (*model.User, error) {
	if adminUser == nil {
		user, err := db.GetUserByRole(model.ADMIN)
		if err != nil {
			return nil, err
		}
		adminUser = user
	}
	return adminUser, nil
}

func GetGeneral() (*model.User, error) {
	if guestUser == nil {
		user, err := db.GetUserByRole(model.GENERAL)
		if err != nil {
			return nil, err
		}
		guestUser = user
	}
	return guestUser, nil
}

func GetUserByRole(role model.Role) (*model.User, error) {
	return db.GetUserByRole(role)
}

func GetUserByName(username string) (*model.User, error) {
	if username == "" {
		return nil, errs.EmptyUsername
	}
	if user, ok := userCache.Get(username); ok {
		return user, nil
	}
	user, err, _ := userG.Do(username, func() (*model.User, error) {
		_user, err := db.GetUserByName(username)
		if err != nil {
			return nil, err
		}
		userCache.Set(username, _user, cache.WithEx[*model.User](time.Hour))
		return _user, nil
	})
	return user, err
}

func GetUserById(id uint) (*model.User, error) {
	return db.GetUserById(id)
}

func GetUsers() (users []model.User, err error) {
	return db.GetUsers()
}

func CreateUser(u *model.User) error {
	return db.CreateUser(u)
}

func DeleteUserById(id uint) error {
	old, err := db.GetUserById(id)
	if err != nil {
		return err
	}
	if old.IsAdmin() || old.IsGeneral() {
		return errs.DeleteAdminOrGuest
	}
	userCache.Del(old.Username)
	return db.DeleteUserById(id)
}

func UpdateUser(u *model.User) error {
	old, err := db.GetUserById(u.ID)
	if err != nil {
		return err
	}
	if u.IsAdmin() {
		adminUser = nil
	}
	if u.IsGeneral() {
		guestUser = nil
	}
	userCache.Del(old.Username)
	return db.UpdateUser(u)
}
