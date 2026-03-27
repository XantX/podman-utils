package container

import (
	"charm.land/bubbletea/v2"
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
		func() ([]components.ContainerItem, error) {
			containers, err := podman.New().ListRunningContainers()
			if err != nil {
				return nil, err
			}
			var items []components.ContainerItem
			for _, c := range containers {
				items = append(items, components.ContainerItem{
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
