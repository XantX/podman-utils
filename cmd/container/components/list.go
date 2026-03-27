package components

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	ID   string
	Name string
}

func (i ContainerItem) FilterValue() string { return i.Name }

type ListModel struct {
	client         *podman.Client
	items          []ContainerItem
	filtered       []ContainerItem
	selectedIndex  int
	filter         string
	title          string
	actionName     string
	onAction       func(id string) error
	fetchItems     func() ([]ContainerItem, error)
	err            error
	successMessage string
}

func NewListModel(
	title string,
	actionName string,
	onAction func(id string) error,
	fetchItems func() ([]ContainerItem, error),
) *ListModel {
	return &ListModel{
		title:      title,
		actionName: actionName,
		onAction:   onAction,
		fetchItems: fetchItems,
		client:     podman.New(),
	}
}

func (m *ListModel) Init() tea.Cmd {
	items, err := m.fetchItems()
	if err != nil {
		m.err = err
		return nil
	}

	m.items = items
	m.filtered = items
	return nil
}

func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case "down", "j":
			if m.selectedIndex < len(m.filtered)-1 {
				m.selectedIndex++
			}
		case "enter":
			if len(m.filtered) > 0 {
				selected := m.filtered[m.selectedIndex]
				if err := m.onAction(selected.ID); err != nil {
					m.err = err
				} else {
					m.successMessage = fmt.Sprintf("%s '%s' (ID: %s) %s", m.actionName, selected.Name, selected.ID, "ejecutado exitosamente")
				}
				return m, tea.Quit
			}
		case "backspace":
			if len(m.filter) > 0 {
				m.filter = m.filter[:len(m.filter)-1]
			}
		default:
			if len(msg.String()) == 1 && msg.String() != "/" {
				m.filter += msg.String()
			}
		}
	case tea.WindowSizeMsg:
		return m, nil
	}

	m.filtered = m.filterItems()
	if m.selectedIndex >= len(m.filtered) {
		m.selectedIndex = len(m.filtered) - 1
	}
	if m.selectedIndex < 0 {
		m.selectedIndex = 0
	}
	return m, nil
}

func (m *ListModel) filterItems() []ContainerItem {
	if m.filter == "" {
		return m.items
	}
	var filtered []ContainerItem
	filterLower := toLower(m.filter)
	for _, item := range m.items {
		if contains(toLower(item.Name), filterLower) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}

func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func (m *ListModel) View() string {
	if m.successMessage != "" {
		return "\n" + SuccessMsgStyle.Render(m.successMessage) + "\n\nPresiona q para salir.\n"
	}

	if m.err != nil {
		return "\n" + ErrorMsgStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n\nPresiona q para salir.\n"
	}

	s := "\n" + TitleStyle.Render(m.title) + "\n\n"

	if len(m.filtered) == 0 {
		s += "No hay contenedores.\n"
	} else {
		for i, item := range m.filtered {
			cursor := "  "
			if i == m.selectedIndex {
				cursor = ">"
			}

			line := fmt.Sprintf("%s %s  %s", cursor, item.ID, item.Name)
			if i == m.selectedIndex {
				s += SelectedItemStyle.Render(line) + "\n"
			} else {
				s += ItemStyle.Render(line) + "\n"
			}
		}
	}

	s += "\n"
	s += HelpStyle.Render("Filtra por nombre (/ para activar): ")
	s += m.filter + "\n"
	s += "\n"
	s += HelpStyle.Render(fmt.Sprintf("Enter: %s | ↑↓: Navegar | /: Filtrar | q/esc: Salir\n", m.actionName))

	return s
}
