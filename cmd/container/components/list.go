package components

import (
	"fmt"
	"io"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/podutil/podutil/internal/podman"
)

var (
	ItemStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("254")).Background(lipgloss.Color("236"))
	SelectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("232")).Background(lipgloss.Color("107"))
	TitleStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("45"))
	SuccessMsgStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true)
	ErrorMsgStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	HelpStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
)

type ContainerItem struct {
	ID      string
	Name    string
	Details podman.ContainerDetail
}

func (i ContainerItem) FilterValue() string { return i.Name }

func (i ContainerItem) Title() string       { return i.Name }
func (i ContainerItem) Description() string { return i.ID }

type ListDelegate struct {
	styles list.DefaultItemStyles
}

func NewListDelegate() *ListDelegate {
	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(lipgloss.Color("232")).Background(lipgloss.Color("107"))
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.Foreground(lipgloss.Color("232")).Background(lipgloss.Color("107"))
	return &ListDelegate{styles: d.Styles}
}

func (d *ListDelegate) Height() int                               { return 1 }
func (d *ListDelegate) Spacing() int                              { return 0 }
func (d *ListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d *ListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	c := item.(ContainerItem)
	cursor := m.Cursor() == index
	if cursor {
		fmt.Fprintf(w, "  %s  %s", c.ID, c.Name)
	} else {
		fmt.Fprintf(w, "  %s  %s", c.ID, c.Name)
	}
}

type ListModel struct {
	list           list.Model
	actionName     string
	onAction       func(id string) error
	fetchItems     func() ([]ContainerItem, error)
	err            error
	successMessage string
	loaded         bool
}

func NewListModel(
	title string,
	actionName string,
	onAction func(id string) error,
	fetchItems func() ([]ContainerItem, error),
) *ListModel {
	items := []list.Item{}
	l := list.New(items, NewListDelegate(), 0, 0)
	l.Title = title
	l.SetShowFilter(true)
	l.SetFilteringEnabled(true)

	return &ListModel{
		list:       l,
		actionName: actionName,
		onAction:   onAction,
		fetchItems: fetchItems,
	}
}

func (m *ListModel) Init() tea.Cmd {
	return func() tea.Msg {
		items, err := m.fetchItems()
		if err != nil {
			return err
		}
		m.loaded = true
		listItems := make([]list.Item, len(items))
		for i, item := range items {
			listItems[i] = item
		}
		return m.list.SetItems(listItems)
	}
}

func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.successMessage != "" {
		return m, tea.Quit
	}

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
			selected, ok := m.list.SelectedItem().(ContainerItem)
			if ok && m.onAction != nil {
				if err := m.onAction(selected.ID); err != nil {
					m.err = err
				} else {
					m.successMessage = fmt.Sprintf("%s '%s' (ID: %s) %s", m.actionName, selected.Name, selected.ID, "ejecutado exitosamente")
				}
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *ListModel) View() tea.View {
	if m.successMessage != "" {
		v := tea.NewView("")
		v.SetContent("\n" + SuccessMsgStyle.Render(m.successMessage) + "\n\nPresiona q para salir.\n")
		return v
	}

	if m.err != nil {
		v := tea.NewView("")
		v.SetContent("\n" + ErrorMsgStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n\nPresiona q para salir.\n")
		return v
	}

	if !m.loaded {
		v := tea.NewView("")
		v.SetContent("\nCargando...\n")
		return v
	}

	s := m.list.View()
	s += "\n" + HelpStyle.Render(fmt.Sprintf("Enter: %s | ↑↓: Navegar | /: Filtrar | q/esc: Salir\n", m.actionName))
	v := tea.NewView(s)
	return v
}
