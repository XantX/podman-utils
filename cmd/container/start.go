package container

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/podutil/podutil/cmd/container/components"
	"github.com/podutil/podutil/internal/podman"
)

func StartCmd() error {
	model := components.NewListModel(
		"Contenedores detenidos",
		"Iniciar",
		func(id string) error {
			return podman.New().StartContainer(id)
		},
		func() ([]components.ContainerItem, error) {
			containers, err := podman.New().ListStoppedContainers()
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
