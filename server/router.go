package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/xqa/chathub/internal/conf"
	"github.com/xqa/chathub/server/common"
	"github.com/xqa/chathub/server/handles"
	"github.com/xqa/chathub/server/middlewares"
	"github.com/xqa/chathub/server/static"
)

func Init(e *gin.Engine) {
	e.Use(cors.Default())

	e.Any("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	common.SecretKey = []byte(conf.Conf.JwtSecret)
	// e.Use(middlewares.StoragesLoaded)
	if conf.Conf.MaxConnections > 0 {
		e.Use(middlewares.MaxAllowed(conf.Conf.MaxConnections))
	}

	api := e.Group("/api")
	api.POST("/auth/login", handles.Login)

	auth := api.Group("", middlewares.Auth)
	auth.GET("/me", handles.CurrentUser)
	auth.POST("/me", handles.UpdateCurrent)

	admin := auth.Group("/admin", middlewares.AuthAdmin)
	{
		user := admin.Group("/user")
		user.GET("/list", handles.ListUsers)
		user.GET("/get", handles.GetUser)
		user.POST("/create", handles.CreateUser)
		user.POST("/update", handles.UpdateUser)
		user.POST("/delete", handles.DeleteUser)
	}

	static.Static(e, func(handlers ...gin.HandlerFunc) {
		e.NoRoute(handlers...)
	})
}
