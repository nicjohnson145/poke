package cmd

import (
	"os"
	"time"

	"github.com/nicjohnson145/poke/config"
	"github.com/nicjohnson145/poke/internal"
	"github.com/spf13/cobra"
)

func Root() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:  "poke",
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
			client := internal.NewHttpClient(internal.HttpClientConfig{
				Logger: config.WithComponent(logger, "httpclient"),
			})
			client.SetTimeout(10 * time.Second)

			runner := internal.NewRunner(internal.RunnerOpts{
				Logger: config.WithComponent(logger, "runner"),
				HttpExecutor: internal.NewHTTPExecutor(internal.HTTPExecutorOpts{
					Logger: config.WithComponent(logger, "httpexecutor"),
					Client: client,
				}),
				GrpcExecutor: internal.NewGRPCExecutor(internal.GRPCExecutorOpts{
					Logger: config.WithComponent(logger, "grpcexecutor"),
				}),
				Parser: internal.NewFSParser(internal.FSParserOpts{
					Logger: config.WithComponent(logger, "fsparser"),
				}),
				Output: os.Stdout,
			})
			return runner.Run(args[0])
		},
	}
	rootCmd.PersistentFlags().BoolP(config.Debug, "d", false, "Enable debug logging")
	rootCmd.Flags().BoolP(config.FailFast, "f", false, "Stop execution on first sequence failure")

	rootCmd.AddCommand(
		versionCmd(),
	)

	return rootCmd
}
