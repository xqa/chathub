package handles

import (
	"time"

	"github.com/Xhofe/go-cache"
	"github.com/gin-gonic/gin"
	"github.com/xqa/chathub/internal/model"
	"github.com/xqa/chathub/internal/op"
	"github.com/xqa/chathub/server/common"
)

var loginCache = cache.NewMemCache[int]()
var (
	defaultDuration = time.Minute * 5
	defaultTimes    = 5
)

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
	OtpCode  string `json:"otp_code"`
}

func Login(c *gin.Context) {
	// check count of login
	ip := c.ClientIP()
	count, ok := loginCache.Get(ip)
	if ok && count >= defaultTimes {
		common.ErrorStrResp(c, "Too many unsuccessful sign-in attempts have been made using an incorrect username or password, Try again later.", 429)
		loginCache.Expire(ip, defaultDuration)
		return
	}
	// check username
	var req LoginReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	user, err := op.GetUserByName(req.Username)
	if err != nil {
		common.ErrorResp(c, err, 400)
		loginCache.Set(ip, count+1)
		return
	}
	// validate password
	if err := user.ValidatePassword(req.Password); err != nil {
		common.ErrorResp(c, err, 400)
		loginCache.Set(ip, count+1)
		return
	}
	// generate token
	token, err := common.GenerateToken(user.Username)
	if err != nil {
		common.ErrorResp(c, err, 400, true)
		return
	}
	common.SuccessResp(c, gin.H{"token": token})
	loginCache.Del(ip)
}

type UserResp struct {
	model.User
}

// CurrentUser get current user by token
// if token is empty, return guest user
func CurrentUser(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	userResp := UserResp{
		User: *user,
	}
	userResp.Password = ""
	common.SuccessResp(c, userResp)
}

func UpdateCurrent(c *gin.Context) {
	var req model.User
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	user := c.MustGet("user").(*model.User)
	user.Username = req.Username
	if req.Password != "" {
		user.Password = req.Password
	}
	if err := op.UpdateUser(user); err != nil {
		common.ErrorResp(c, err, 500)
	} else {
		common.SuccessResp(c)
	}
}
