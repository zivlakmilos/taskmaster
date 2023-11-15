package tui

import (
	"os"
	"path"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jmoiron/sqlx"
	"github.com/zivlakmilos/taskmaster/db"
)

type AppState int

const (
	AppStateProjects AppState = iota
	AppStateTasks

	AppStateLen
)

type Config struct {
	ConfigDir string
}

type MainModel struct {
	state  AppState
	models [AppStateLen]tea.Model

	cfg Config
}

func NewMainModel(cfg Config) *MainModel {
	models := [AppStateLen]tea.Model{}

	models[AppStateProjects] = NewProjectsModel(cfg)

	return &MainModel{
		state:  AppStateProjects,
		models: models,
		cfg:    cfg,
	}
}

func (m *MainModel) Init() tea.Cmd {
	return runMigrations(m.cfg)
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
		break
	}

	var cmd tea.Cmd
	m.models[m.state], cmd = m.models[m.state].Update(msg)

	return m, cmd
}

func (m *MainModel) View() string {
	return m.models[m.state].View()
}

func openDb(cfg Config) (*sqlx.DB, error) {
	dbUrl := path.Join(cfg.ConfigDir, "database.db")
	return db.Open(dbUrl)
}

func runMigrations(cfg Config) tea.Cmd {
	return func() tea.Msg {
		dbUrl := path.Join(cfg.ConfigDir, "database.db")

		if _, err := os.Stat(dbUrl); os.IsExist(err) {
			return nil
		}

		if _, err := os.Stat(cfg.ConfigDir); os.IsNotExist(err) {
			err := os.Mkdir(cfg.ConfigDir, os.ModePerm)
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
}
