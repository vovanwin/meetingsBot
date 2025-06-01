package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"github.com/vovanwin/meetingsBot/cmd/dependency"
	"github.com/vovanwin/meetingsBot/internal/telegramBoll"
)

var (
	Version = "0.1"

	rootCmd = &cobra.Command{
		Use:     "server",
		Version: Version,
		Short:   "Запуск Http REST API",
		Run: func(cmd *cobra.Command, args []string) {
			fx.New(inject()).Run()
		},
	}
)

func inject() fx.Option {
	return fx.Options(
		//fx.NopLogger,

		fx.Provide(
			dependency.ProvideConfig,
			dependency.ProvideLogger,
			dependency.ProvidePgx,
		),

		telegramBoll.Module,
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
