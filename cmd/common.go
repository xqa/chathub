package cmd

import (
	"github.com/xqa/chathub/internal/bootstrap"
	"github.com/xqa/chathub/internal/bootstrap/data"
)

func Init() {
	bootstrap.InitConfig()
	bootstrap.Log()
	bootstrap.InitDB()
	data.InitData()
}
