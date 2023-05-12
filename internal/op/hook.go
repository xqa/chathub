package op

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/xqa/chathub/internal/conf"
	"github.com/xqa/chathub/internal/model"
)

// Setting
type SettingItemHook func(item *model.SettingItem) error

var settingItemHooks = map[string]SettingItemHook{
	conf.PrivacyRegs: func(item *model.SettingItem) error {
		regStrs := strings.Split(item.Value, "\n")
		regs := make([]*regexp.Regexp, 0, len(regStrs))
		for _, regStr := range regStrs {
			reg, err := regexp.Compile(regStr)
			if err != nil {
				return errors.WithStack(err)
			}
			regs = append(regs, reg)
		}
		conf.PrivacyReg = regs
		return nil
	},
}

func RegisterSettingItemHook(key string, hook SettingItemHook) {
	settingItemHooks[key] = hook
}

func HandleSettingItemHook(item *model.SettingItem) (hasHook bool, err error) {
	if hook, ok := settingItemHooks[item.Key]; ok {
		return true, hook(item)
	}
	return false, nil
}
