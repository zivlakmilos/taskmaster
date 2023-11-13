package tui

import tea "github.com/charmbracelet/bubbletea"

type AppState int

const (
	AppStateProjects AppState = iota
	AppStateTasks

	AppStateLen
)

type MainModel struct {
	state AppState

	models [AppStateLen]tea.Model
}

func NewMainModel() *MainModel {
	models := [AppStateLen]tea.Model{}

	models[AppStateProjects] = NewProjectsModel()

	return &MainModel{
		state:  AppStateProjects,
		models: models,
	}
}

func (m *MainModel) Init() tea.Cmd {
	return nil
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
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
