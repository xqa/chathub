package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xqa/chathub/cmd/flags"
)

var RootCmd = &cobra.Command{
	Use:   "chathub",
	Short: "A chatgpt web app.",
	Long:  `A chatgpt web app`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&flags.Mode, "mode", "dev", `mode of server, can be "prod" or "dev"`)
	RootCmd.PersistentFlags().StringVar(&flags.DataDir, "data", "data", "config file")
	RootCmd.PersistentFlags().BoolVar(&flags.LogStd, "log-std", false, "Force to log to std")
}
