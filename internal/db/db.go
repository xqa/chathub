package db

import (
	log "github.com/sirupsen/logrus"

	"github.com/xqa/chathub/internal/model"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(d *gorm.DB) {
	db = d
	err := db.AutoMigrate(new(model.User), new(model.SettingItem))
	if err != nil {
		log.Fatalf("failed migrate database: %s", err.Error())
	}
}

func GetDb() *gorm.DB {
	return db
}
