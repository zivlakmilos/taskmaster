package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zivlakmilos/taskmaster/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "TUI",
	Run: func(cmd *cobra.Command, args []string) {
		showTui()
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}

func showTui() {
	mainModel := tui.NewMainModel(tui.Config{ConfigDir: cfg.configDir})
	p := tea.NewProgram(mainModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		exitWithError(err)
	}
}
