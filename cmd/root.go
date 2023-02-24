package cmd

import (
	"net/http"
	"time"

	"github.com/nicjohnson145/poke/config"
	"github.com/nicjohnson145/poke/internal"
	"github.com/spf13/cobra"
)

func Root() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "poke",
		Args: cobra.ExactArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// So we don't print usage messages on execution errors
			cmd.SilenceUsage = true
			// So we dont double report errors
			cmd.SilenceErrors = true
			return config.InitializeConfig(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := config.InitLogger()
			runner := internal.NewRunner(internal.RunnerOpts{
				Logger: config.WithComponent(logger, "runner"),
				HttpExecutor: internal.NewHTTPExecutor(internal.HTTPExecutorOpts{
					Logger: config.WithComponent(logger, "httpexecutor"),
					Client: &http.Client{
						Timeout: 10 * time.Second,
					},
				}),
				Parser: internal.NewFSParser(internal.FSParserOpts{
					Logger: config.WithComponent(logger, "fsparser"),
				}),
			})
			return runner.Run(args[0])
		},
	}
	rootCmd.PersistentFlags().BoolP(config.Debug, "d", false, "Enable debug logging")

	return rootCmd
}
