package cmd

import (
	"fmt"

	"github.com/nicjohnson145/poke/config"
	"github.com/spf13/cobra"
)

func Root() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "poke",
		Short: "",
		Long:  "",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// So we don't print usage messages on execution errors
			cmd.SilenceUsage = true
			// So we dont double report errors
			cmd.SilenceErrors = true
			return config.InitializeConfig(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("hello world")
			return nil
		},
	}
	rootCmd.PersistentFlags().BoolP(config.Debug, "d", false, "Enable debug logging")

	return rootCmd
}
