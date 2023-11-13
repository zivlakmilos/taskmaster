package tui

import tea "github.com/charmbracelet/bubbletea"

type ProjectsModel struct {
}

func NewProjectsModel() *ProjectsModel {
	return &ProjectsModel{}
}

func (m *ProjectsModel) Init() tea.Cmd {
	return nil
}

func (m *ProjectsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *ProjectsModel) View() string {
	return ""
}
