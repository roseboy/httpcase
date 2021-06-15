package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "hc",
		Short: "A very NB tool for auto api test.",
	}
)

func init() {
	rootCmd.AddCommand(newRunCmd().cmd, newDemoCmd().cmd, newVersionCmd().cmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
