package middlewares

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/xqa/chathub/internal/model"
	"github.com/xqa/chathub/internal/op"
	"github.com/xqa/chathub/server/common"
)

// Auth is a middleware that checks if the user is logged in.
func Auth(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		common.ErrorStrResp(c, "token not found, login please", 401)
		c.Abort()
		return
	}
	userClaims, err := common.ParseToken(token)
	if err != nil {
		common.ErrorResp(c, err, 401)
		c.Abort()
		return
	}
	user, err := op.GetUserByName(userClaims.Username)
	if err != nil {
		common.ErrorResp(c, err, 401)
		c.Abort()
		return
	}
	c.Set("user", user)
	log.Debugf("use login token: %+v", user)
	c.Next()
}

func AuthAdmin(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	if !user.IsAdmin() {
		common.ErrorStrResp(c, "You are not an admin", 403)
		c.Abort()
	} else {
		c.Next()
	}
}
