package model

type SettingItem struct {
	Key     string `json:"key" gorm:"primaryKey" binding:"required"` // unique key
	Value   string `json:"value"`                                    // value
	Help    string `json:"help"`                                     // help message
	Type    string `json:"type"`                                     // string, number, bool, select
	Options string `json:"options"`                                  // values for select
}
