package container

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/podutil/podutil/internal/podman"
)

var (
	itemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("254")).
			Background(lipgloss.Color("236"))

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("232")).
				Background(lipgloss.Color("107"))

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("45"))
)

type ContainerItem struct {
	id   string
	name string
}

func (i ContainerItem) FilterValue() string { return i.name }

type Model struct {
	client        *podman.Client
	containers    []ContainerItem
	filtered      []ContainerItem
	selectedIndex int
	filter        string
	err           error
}

func NewModel() *Model {
	return &Model{
		client: podman.New(),
	}
}

func (m *Model) Init() tea.Cmd {
	containers, err := m.client.ListStoppedContainers()
	if err != nil {
		m.err = err
		return nil
	}

	var items []ContainerItem
	for _, c := range containers {
		items = append(items, ContainerItem{
			id:   c.ID,
			name: c.Name,
		})
	}

	m.containers = items
	m.filtered = items
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				if err := m.client.StartContainer(selected.id); err != nil {
					m.err = err
				}
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		return m, nil
	}

	m.filtered = m.filterContainers()
	if m.selectedIndex >= len(m.filtered) {
		m.selectedIndex = len(m.filtered) - 1
	}
	return m, nil
}

func (m *Model) filterContainers() []ContainerItem {
	if m.filter == "" {
		return m.containers
	}
	var filtered []ContainerItem
	for _, c := range m.containers {
		if containsIgnoreCase(c.name, m.filter) {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

func containsIgnoreCase(s, substr string) bool {
	s = toLower(s)
	substr = toLower(substr)
	return contains(s, substr)
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
	return len(s) >= len(substr) && (s == substr || len(s) == 0 ||
		(containsAt(s, substr) != -1))
}

func containsAt(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func (m *Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	s := "\n" + titleStyle.Render("Contenedores detenidos") + "\n\n"

	if len(m.filtered) == 0 {
		s += "No hay contenedores detenidos.\n"
	} else {
		for i, c := range m.filtered {
			cursor := "  "
			if i == m.selectedIndex {
				cursor = ">"
			}

			line := fmt.Sprintf("%s %s  %s", cursor, c.id, c.name)
			if i == m.selectedIndex {
				s += selectedItemStyle.Render(line) + "\n"
			} else {
				s += itemStyle.Render(line) + "\n"
			}
		}
	}

	s += "\n"
	s += "Filtra por nombre: "
	s += m.filter + "\n"
	s += "\nEnter: Iniciar | q/esc: Salir\n"

	return s
}

func StartCmd() error {
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
