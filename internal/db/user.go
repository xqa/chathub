package db

import (
	"github.com/pkg/errors"
	"github.com/xqa/chathub/internal/model"
)

func GetUserByRole(role model.Role) (*model.User, error) {
	user := model.User{Role: role}
	if err := db.Where(user).Take(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByName(username string) (*model.User, error) {
	user := model.User{Username: username}
	if err := db.Where(user).First(&user).Error; err != nil {
		return nil, errors.Wrapf(err, "failed find user")
	}
	return &user, nil
}

func GetUserById(id uint) (*model.User, error) {
	var u model.User
	if err := db.First(&u, id).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get old user")
	}
	return &u, nil
}

func CreateUser(u *model.User) error {
	return errors.WithStack(db.Create(u).Error)
}

func UpdateUser(u *model.User) error {
	return errors.WithStack(db.Save(u).Error)
}

func GetUsers() (users []model.User, err error) {
	userDB := db.Model(&model.User{})

	if err := userDB.Find(&users).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get find users")
	}
	return users, nil
}

func DeleteUserById(id uint) error {
	return errors.WithStack(db.Delete(&model.User{}, id).Error)
}
