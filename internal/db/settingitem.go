package db

import (
	"github.com/pkg/errors"
	"github.com/xqa/chathub/internal/model"
)

func GetSettingItems() ([]model.SettingItem, error) {
	var settingItems []model.SettingItem
	if err := db.Find(&settingItems).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return settingItems, nil
}

func GetSettingItemByKey(key string) (*model.SettingItem, error) {
	var settingItem model.SettingItem
	if err := db.Where("key = ?", key).First(&settingItem).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return &settingItem, nil
}

func GetSettingItemsByGroup(group int) ([]model.SettingItem, error) {
	var settingItems []model.SettingItem
	if err := db.Where("`group` in ?", group).Find(&settingItems).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return settingItems, nil
}

func GetSettingItemsInGroups(groups []int) ([]model.SettingItem, error) {
	var settingItems []model.SettingItem
	if err := db.Where("`group` in ?", groups).Find(&settingItems).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return settingItems, nil
}

func SaveSettingItems(items []model.SettingItem) (err error) {
	return errors.WithStack(db.Save(items).Error)
}

func SaveSettingItem(item *model.SettingItem) error {
	return errors.WithStack(db.Save(item).Error)
}

func DeleteSettingItemByKey(key string) error {
	return errors.WithStack(db.Delete(&model.SettingItem{Key: key}).Error)
}
