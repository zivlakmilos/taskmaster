package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zivlakmilos/taskmaster/db"
)

type ProjectsModel struct {
	cfg  Config
	list list.Model
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type ProjectItem struct {
	title       string
	description string
}

func (i ProjectItem) Title() string       { return i.title }
func (i ProjectItem) Description() string { return i.description }
func (i ProjectItem) FilterValue() string { return i.title }

func NewProjectsModel(cfg Config) *ProjectsModel {
	items := []list.Item{
		ProjectItem{title: "All", description: "All tasks"},
	}
	items = append(items, readProjects(cfg)...)
	list := list.New(items, list.NewDefaultDelegate(), 0, 0)

	list.Title = "Projects"

	return &ProjectsModel{
		cfg:  cfg,
		list: list,
	}
}

func (m *ProjectsModel) Init() tea.Cmd {
	return nil
}

func (m *ProjectsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		break
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *ProjectsModel) View() string {
	return m.list.View()
}

func readProjects(cfg Config) []list.Item {
	projectList := []list.Item{}

	con, err := openDb(cfg)
	if err != nil {
		return projectList
	}

	store := db.NewProjectStore(con)
	projects, err := store.GetAll()

	for _, project := range projects {
		projectList = append(projectList, ProjectItem{title: project.Name, description: string(project.Status)})
	}

	return projectList
}
