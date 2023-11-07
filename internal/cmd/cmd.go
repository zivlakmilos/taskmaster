package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"github.com/zivlakmilos/taskmaster/db"
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

func openDb() (*sqlx.DB, error) {
	err := runMigrations()
	if err != nil {
		return nil, err
	}

	dbUrl := path.Join(cfg.configDir, "database.db")
	return db.Open(dbUrl)
}

func runMigrations() error {
	dbUrl := path.Join(cfg.configDir, "database.db")

	if _, err := os.Stat(dbUrl); os.IsExist(err) {
		return nil
	}

	if _, err := os.Stat(cfg.configDir); os.IsNotExist(err) {
		err := os.Mkdir(cfg.configDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	con, err := db.Open(dbUrl)
	if err != nil {
		return err
	}
	defer con.Close()

	err = db.RunMigrations(con)
	if err != nil {
		return err
	}

	return nil
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err)
	os.Exit(1)
}
