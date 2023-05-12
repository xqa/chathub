package bootstrap

import (
	"fmt"
	stdlog "log"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/xqa/chathub/cmd/flags"
	"github.com/xqa/chathub/internal/conf"
	"github.com/xqa/chathub/internal/db"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitDB() {
	logLevel := logger.Silent
	if flags.Mode == "dev" {
		logLevel = logger.Info
	}
	newLogger := logger.New(
		stdlog.New(log.StandardLogger().Out, "\r\n", stdlog.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: conf.Conf.Database.TablePrefix,
		},
		Logger: newLogger,
	}
	var dB *gorm.DB
	var err error
	if flags.Mode == "dev" {
		dB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), gormConfig)
	} else {
		database := conf.Conf.Database
		if !(strings.HasSuffix(database.DBFile, ".db") && len(database.DBFile) > 3) {
			log.Fatalf("db name error.")
		}
		dB, err = gorm.Open(sqlite.Open(fmt.Sprintf("%s?_journal=WAL&_vacuum=incremental",
			database.DBFile)), gormConfig)
	}
	if err != nil {
		log.Fatalf("failed to connect database:%s", err.Error())
	}
	db.Init(dB)
}
