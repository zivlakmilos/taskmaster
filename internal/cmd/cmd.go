package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
)

type Config struct {
	configDir string
}

var cfg Config

var rootCmd = &cobra.Command{
	Use:     "taskmaster",
	Short:   "Taskmaster is CLI task manager",
	Version: "1.0.0",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		exitWithError(err)
	}
	configDir = path.Join(configDir, "taskmaster")

	rootCmd.PersistentFlags().StringVarP(&cfg.configDir, "configDir", "c", configDir, "Config directory")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		exitWithError(err)
	}
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err)
	os.Exit(1)
}
