package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zivlakmilos/taskmaster/db"
)

type ProjectsMode int

const (
	ProjectsModeNormal ProjectsMode = iota
	ProjectsModeNew
	ProjectsModeRename
	ProjectsModeDelete
)

type ProjectsModel struct {
	cfg  Config
	mode ProjectsMode

	list  list.Model
	input textinput.Model
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type ProjectItem struct {
	id          string
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
	list.DisableQuitKeybindings()
	list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("a"),
				key.WithHelp("a", "add project"),
			),
			key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "rename project"),
			),
			key.NewBinding(
				key.WithKeys("d"),
				key.WithHelp("d", "delete project"),
			),
			key.NewBinding(
				key.WithKeys("q"),
				key.WithHelp("q", "quit"),
			),
		}
	}

	return &ProjectsModel{
		cfg:   cfg,
		list:  list,
		input: textinput.New(),
	}
}

func (m *ProjectsModel) Init() tea.Cmd {
	return nil
}

func (m *ProjectsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		break
	}

	switch m.mode {
	case ProjectsModeNormal:
		cmd = m.handleUpdateNormal(msg)
		break
	case ProjectsModeNew:
		cmd = m.handleUpdateNew(msg)
		break
	case ProjectsModeRename:
		cmd = m.handleUpdateRename(msg)
		break
	case ProjectsModeDelete:
		cmd = m.handleUpdateDelete(msg)
		break
	}
	cmds = append(cmds, cmd)

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *ProjectsModel) View() string {
	res := m.list.View()

	if m.mode == ProjectsModeNew || m.mode == ProjectsModeRename {
		res += "\n\n"
		res += m.input.View()
	}

	if m.mode == ProjectsModeDelete {
		res += "\n\n"
		res += "press d again to delete or esc to abort"
	}

	return res
}

func (m *ProjectsModel) handleUpdateNormal(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			m.mode = ProjectsModeNew
			m.input.SetValue("")
			m.input.Focus()
			break
		case "r":
			m.mode = ProjectsModeRename
			m.input.SetValue("")
			m.input.Focus()
			break
		case "d":
			m.mode = ProjectsModeDelete
			break
		case "q":
			return tea.Quit
		}
		break
	}

	return nil
}

func (m *ProjectsModel) handleUpdateNew(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.mode = ProjectsModeNormal
			return nil
		case "enter":
			m.mode = ProjectsModeNormal
			return addProject(m.cfg, m.input.Value())
		}
		break
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return cmd
}

func (m *ProjectsModel) handleUpdateRename(msg tea.Msg) tea.Cmd {
	return nil
}

func (m *ProjectsModel) handleUpdateDelete(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "d":
			item := m.list.SelectedItem().(ProjectItem)
			m.mode = ProjectsModeNormal
			return removeProject(m.cfg, item.id)
		case "esc":
			m.mode = ProjectsModeNormal
			return nil
		}
	}

	return nil
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
		projectList = append(projectList, ProjectItem{id: project.Id, title: project.Name, description: string(project.Status)})
	}

	return projectList
}

func addProject(cfg Config, name string) tea.Cmd {
	return func() tea.Msg {
		con, err := openDb(cfg)
		if err != nil {
			return err
		}

		store := db.NewProjectStore(con)

		project := db.NewProject()
		project.Name = name
		err = store.Save(project)
		if err != nil {
			return err
		}

		return nil
	}
}

func removeProject(cfg Config, id string) tea.Cmd {
	return func() tea.Msg {
		con, err := openDb(cfg)
		if err != nil {
			return err
		}

		store := db.NewProjectStore(con)
		store.Delete(id)

		return nil
	}
}
