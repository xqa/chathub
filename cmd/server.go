package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/xqa/chathub/cmd/flags"
	"github.com/xqa/chathub/internal/conf"
	"github.com/xqa/chathub/pkg/utils"
	"github.com/xqa/chathub/server"
)

// ServerCmd represents the server command
var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server at the specified address",
	Long: `Start the server at the specified address
the address is defined in config file`,
	Run: func(cmd *cobra.Command, args []string) {
		Init()
		if flags.Mode == "prod" {
			gin.SetMode(gin.ReleaseMode)
		}
		r := gin.New()
		r.Use(gin.LoggerWithWriter(log.StandardLogger().Out), gin.RecoveryWithWriter(log.StandardLogger().Out))
		server.Init(r)
		base := fmt.Sprintf("%s:%d", conf.Conf.Address, conf.Conf.Port)
		utils.Log.Infof("start server @ %s", base)
		srv := &http.Server{Addr: base, Handler: r}
		go func() {
			err := srv.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				utils.Log.Fatalf("failed to start: %s", err.Error())
			}
		}()
		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 5 seconds.
		quit := make(chan os.Signal)
		// kill (no param) default send syscanll.SIGTERM
		// kill -2 is syscall.SIGINT
		// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		utils.Log.Println("Shutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			utils.Log.Fatal("Server Shutdown:", err)
		}
		<-ctx.Done()
		utils.Log.Println("Server exiting")
	},
}

func init() {
	RootCmd.AddCommand(ServerCmd)
}
