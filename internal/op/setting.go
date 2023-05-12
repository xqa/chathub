package op

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Xhofe/go-cache"
	"github.com/xqa/chathub/internal/db"
	"github.com/xqa/chathub/internal/model"
	"github.com/xqa/chathub/pkg/singleflight"
	"github.com/xqa/chathub/pkg/utils"
)

var settingCache = cache.NewMemCache(cache.WithShards[*model.SettingItem](4))
var settingG singleflight.Group[*model.SettingItem]
var settingCacheF = func(item *model.SettingItem) {
	settingCache.Set(item.Key, item, cache.WithEx[*model.SettingItem](time.Hour))
}

var settingGroupCache = cache.NewMemCache(cache.WithShards[[]model.SettingItem](4))
var settingGroupG singleflight.Group[[]model.SettingItem]
var settingGroupCacheF = func(key string, item []model.SettingItem) {
	settingGroupCache.Set(key, item, cache.WithEx[[]model.SettingItem](time.Hour))
}

func settingCacheUpdate() {
	settingCache.Clear()
	settingGroupCache.Clear()
}

func GetSettingsMap() map[string]string {
	items, _ := GetSettingItems()
	settings := make(map[string]string)
	for _, item := range items {
		settings[item.Key] = item.Value
	}
	return settings
}

func GetSettingItems() ([]model.SettingItem, error) {
	if items, ok := settingGroupCache.Get("ALL_SETTING_ITEMS"); ok {
		return items, nil
	}
	items, err, _ := settingGroupG.Do("ALL_SETTING_ITEMS", func() ([]model.SettingItem, error) {
		_items, err := db.GetSettingItems()
		if err != nil {
			return nil, err
		}
		settingGroupCacheF("ALL_SETTING_ITEMS", _items)
		return _items, nil
	})
	return items, err
}

func GetSettingItemByKey(key string) (*model.SettingItem, error) {
	if item, ok := settingCache.Get(key); ok {
		return item, nil
	}

	item, err, _ := settingG.Do(key, func() (*model.SettingItem, error) {
		_item, err := db.GetSettingItemByKey(key)
		if err != nil {
			return nil, err
		}
		settingCacheF(_item)
		return _item, nil
	})
	return item, err
}

func GetSettingItemInKeys(keys []string) ([]model.SettingItem, error) {
	var items []model.SettingItem
	for _, key := range keys {
		item, err := GetSettingItemByKey(key)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	return items, nil
}

func GetSettingItemsByGroup(group int) ([]model.SettingItem, error) {
	key := strconv.Itoa(group)
	if items, ok := settingGroupCache.Get(key); ok {
		return items, nil
	}
	items, err, _ := settingGroupG.Do(key, func() ([]model.SettingItem, error) {
		_items, err := db.GetSettingItemsByGroup(group)
		if err != nil {
			return nil, err
		}
		settingGroupCacheF(key, _items)
		return _items, nil
	})
	return items, err
}

func GetSettingItemsInGroups(groups []int) ([]model.SettingItem, error) {
	sort.Ints(groups)
	key := strings.Join(utils.MustSliceConvert(groups, func(i int) string {
		return strconv.Itoa(i)
	}), ",")

	if items, ok := settingGroupCache.Get(key); ok {
		return items, nil
	}
	items, err, _ := settingGroupG.Do(key, func() ([]model.SettingItem, error) {
		_items, err := db.GetSettingItemsInGroups(groups)
		if err != nil {
			return nil, err
		}
		settingGroupCacheF(key, _items)
		return _items, nil
	})
	return items, err
}

func SaveSettingItems(items []model.SettingItem) error {
	err := db.SaveSettingItems(items)
	return err
}

func SaveSettingItem(item *model.SettingItem) (err error) {
	if err = db.SaveSettingItem(item); err != nil {
		return err
	}
	settingCacheUpdate()
	return nil
}

func DeleteSettingItemByKey(key string) error {
	settingCacheUpdate()
	return db.DeleteSettingItemByKey(key)
}
