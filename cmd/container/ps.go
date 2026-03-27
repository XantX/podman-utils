package container

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/podutil/podutil/cmd/container/components"
	"github.com/podutil/podutil/internal/podman"
)

type PsListModel struct {
	client        *podman.Client
	containers    []podman.ContainerDetail
	filtered      []podman.ContainerDetail
	selectedIndex int
	filter        string
	err           error
}

func (m *PsListModel) Init() tea.Cmd {
	containers, err := m.client.ListRunningContainersDetailed()
	if err != nil {
		m.err = err
		return nil
	}

	m.containers = containers
	m.filtered = containers
	return nil
}

func (m *PsListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				detailModel := components.NewDetailModel(selected, func() tea.Model {
					return m
				})
				return detailModel, nil
			}
		}
	case tea.WindowSizeMsg:
		return m, nil
	}

	m.filterContainers()
	if m.selectedIndex >= len(m.filtered) {
		m.selectedIndex = len(m.filtered) - 1
	}
	if m.selectedIndex < 0 {
		m.selectedIndex = 0
	}
	return m, nil
}

func (m *PsListModel) filterContainers() {
	if m.filter == "" {
		m.filtered = m.containers
		return
	}
	var filtered []podman.ContainerDetail
	filterLower := toLower(m.filter)
	for _, c := range m.containers {
		if contains(toLower(c.Name), filterLower) {
			filtered = append(filtered, c)
		}
	}
	m.filtered = filtered
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

func (m *PsListModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	s := "\n" + components.TitleStyle.Render("Contenedores corriendo") + "\n\n"

	if len(m.filtered) == 0 {
		s += "No hay contenedores corriendo.\n"
	} else {
		for i, c := range m.filtered {
			cursor := "  "
			if i == m.selectedIndex {
				cursor = ">"
			}

			line := fmt.Sprintf("%s %s  %s", cursor, c.ID, c.Name)
			if i == m.selectedIndex {
				s += components.SelectedItemStyle.Render(line) + "\n"
			} else {
				s += components.ItemStyle.Render(line) + "\n"
			}
		}
	}

	s += "\n"
	s += components.HelpStyle.Render("Filtra por nombre: ")
	s += m.filter + "\n"
	s += "\n"
	s += components.HelpStyle.Render("Enter: Ver detalles | q/esc: Salir\n")

	return s
}

func PsCmd() error {
	model := &PsListModel{
		client: podman.New(),
	}

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
