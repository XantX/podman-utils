package container

import (
	"fmt"
	"io"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/podutil/podutil/cmd/container/components"
	"github.com/podutil/podutil/internal/podman"
)

var (
	PsItemStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("254")).Background(lipgloss.Color("236"))
	PsSelectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("232")).Background(lipgloss.Color("107"))
)

type PsContainerItem struct {
	ID      string
	Name    string
	Details podman.ContainerDetail
}

func (i PsContainerItem) FilterValue() string { return i.Name }

func (i PsContainerItem) Title() string       { return i.Name }
func (i PsContainerItem) Description() string { return i.ID }

type PsDelegate struct{}

func (d PsDelegate) Height() int                               { return 1 }
func (d PsDelegate) Spacing() int                              { return 0 }
func (d PsDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d PsDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	c := item.(PsContainerItem)
	cursor := m.Cursor() == index
	if cursor {
		fmt.Fprintf(w, "  %s  %s", c.ID, c.Name)
	} else {
		fmt.Fprintf(w, "  %s  %s", c.ID, c.Name)
	}
}

type PsListModel struct {
	list       list.Model
	client     *podman.Client
	containers []PsContainerItem
	onDetail   func(podman.ContainerDetail) tea.Model
	err        error
	loaded     bool
}

func NewPsListModel(onDetail func(podman.ContainerDetail) tea.Model) *PsListModel {
	items := []list.Item{}
	l := list.New(items, PsDelegate{}, 0, 0)
	l.Title = "Contenedores corriendo"
	l.SetShowFilter(true)
	l.SetFilteringEnabled(true)

	return &PsListModel{
		list:     l,
		client:   podman.New(),
		onDetail: onDetail,
	}
}

func (m *PsListModel) Init() tea.Cmd {
	return func() tea.Msg {
		containers, err := m.client.ListRunningContainersDetailed()
		if err != nil {
			return err
		}
		m.containers = make([]PsContainerItem, len(containers))
		for i, c := range containers {
			m.containers[i] = PsContainerItem{
				ID:      c.ID,
				Name:    c.Name,
				Details: c,
			}
		}
		m.loaded = true
		listItems := make([]list.Item, len(m.containers))
		for i, item := range m.containers {
			listItems[i] = item
		}
		return m.list.SetItems(listItems)
	}
}

func (m *PsListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.err != nil {
		return m, nil
	}

	switch msg := msg.(type) {
	case error:
		m.err = msg
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "enter":
			selected, ok := m.list.SelectedItem().(PsContainerItem)
			if ok && m.onDetail != nil {
				return m.onDetail(selected.Details), nil
			}
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *PsListModel) View() tea.View {
	if m.err != nil {
		v := tea.NewView("")
		v.SetContent("\n" + components.ErrorMsgStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n\nPresiona q para salir.\n")
		return v
	}

	if !m.loaded {
		v := tea.NewView("")
		v.SetContent("\nCargando...\n")
		return v
	}

	s := m.list.View()
	s += "\n" + components.HelpStyle.Render("Enter: Ver detalles | ↑↓: Navegar | /: Filtrar | q/esc: Salir\n")
	v := tea.NewView(s)
	return v
}

func PsCmd() error {
	model := NewPsListModel(func(detail podman.ContainerDetail) tea.Model {
		return components.NewDetailModel(detail, func() tea.Model {
			return &PsListModel{}
		})
	})

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
