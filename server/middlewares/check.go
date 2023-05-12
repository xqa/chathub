package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xqa/chathub/internal/conf"
	"github.com/xqa/chathub/pkg/utils"
	"github.com/xqa/chathub/server/common"
)

func StoragesLoaded(c *gin.Context) {
	if conf.StoragesLoaded {
		c.Next()
	} else {
		if utils.SliceContains([]string{"", "/", "/favicon.ico"}, c.Request.URL.Path) {
			c.Next()
			return
		}
		paths := []string{"/assets", "/images", "/streamer", "/static"}
		for _, path := range paths {
			if strings.HasPrefix(c.Request.URL.Path, path) {
				c.Next()
				return
			}
		}
		common.ErrorStrResp(c, "Loading storage, please wait", 500)
		c.Abort()
	}
}
