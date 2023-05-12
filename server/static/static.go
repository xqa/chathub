package static

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xqa/chathub/internal/conf"
	"github.com/xqa/chathub/internal/setting"
	"github.com/xqa/chathub/pkg/utils"
	"github.com/xqa/chathub/public"
)

func InitIndex() {
	index, err := public.Public.ReadFile("dist/index.html")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			utils.Log.Fatalf("index.html not exist, you may forget to put dist of frontend to public/dist")
		}
		utils.Log.Fatalf("failed to read index.html: %v", err)
	}
	conf.RawIndexHtml = string(index)
	UpdateIndex()
}

func UpdateIndex() {
	title := setting.GetStr(conf.SiteTitle)
	conf.ManageHtml = conf.RawIndexHtml
	replaceMap1 := map[string]string{
		"Loading...": title,
	}
	for k, v := range replaceMap1 {
		conf.ManageHtml = strings.Replace(conf.ManageHtml, k, v, 1)
	}
	conf.IndexHtml = conf.ManageHtml
}

func Static(r *gin.Engine, noRoute func(handlers ...gin.HandlerFunc)) {
	InitIndex()
	folders := []string{"assets", "images", "streamer", "static"}
	r.Use(func(c *gin.Context) {
		for i := range folders {
			if strings.HasPrefix(c.Request.RequestURI, fmt.Sprintf("/%s/", folders[i])) {
				c.Header("Cache-Control", "public, max-age=15552000")
			}
		}
	})
	for i, folder := range folders {
		folder = "dist/" + folder
		sub, err := fs.Sub(public.Public, folder)
		if err != nil {
			utils.Log.Fatalf("can't find folder: %s", folder)
		}
		r.StaticFS(fmt.Sprintf("/%s/", folders[i]), http.FS(sub))
	}

	noRoute(func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.Status(200)
		_, _ = c.Writer.WriteString(conf.IndexHtml)
		c.Writer.Flush()
		c.Writer.WriteHeaderNow()
	})
}
