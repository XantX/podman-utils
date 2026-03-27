package container

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/podutil/podutil/cmd/container/components"
	"github.com/podutil/podutil/internal/podman"
)

func StopCmd() error {
	model := components.NewListModel(
		"Contenedores corriendo",
		"Detener",
		func(id string) error {
			return podman.New().StopContainer(id)
		},
		func() ([]components.ListItem, error) {
			containers, err := podman.New().ListRunningContainers()
			if err != nil {
				return nil, err
			}
			var items []components.ListItem
			for _, c := range containers {
				items = append(items, components.ListItem{
					ID:   c.ID,
					Name: c.Name,
				})
			}
			return items, nil
		},
	)

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
