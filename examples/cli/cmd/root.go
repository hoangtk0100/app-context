package cmd

import (
	"fmt"
	"os"

	appctx "github.com/hoangtk0100/app-context"
	"github.com/spf13/cobra"
)

func newAppCtx() appctx.AppContext {
	return appctx.NewAppContext(
		appctx.WithName("demo-cli"),
	)
}

var outEnvCMD = &cobra.Command{
	Use:   "outenv",
	Short: "Output all environment variables to std",
	Run: func(cmd *cobra.Command, args []string) {
		newAppCtx().OutEnv()
	},
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start app",
	Run: func(cmd *cobra.Command, args []string) {
		appCtx := newAppCtx()
		log := appCtx.Logger("service")

		if err := appCtx.Load(); err != nil {
			log.Error(err)
		}
	},
}

func Execute() {
	rootCmd.AddCommand(outEnvCMD)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
