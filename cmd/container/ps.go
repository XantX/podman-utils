package container

import (
	tea "charm.land/bubbletea/v2"
	"github.com/podutil/podutil/cmd/container/components"
	"github.com/podutil/podutil/internal/podman"
)

func PsCmd() error {
	var model *components.DetailListModel
	model = components.NewDetailListModel(
		"Todos los contenedores",
		func(item components.ContainerItem) tea.Model {
			// When an item is selected, return a detail model
			// The onBack function returns the same list model (captured via closure)
			return components.NewDetailModel(item.Details, func() tea.Model {
				return model
			})
		},
		func() ([]components.ContainerItem, error) {
			client := podman.New()
			containers, err := client.ListRunningContainersDetailed()
			if err != nil {
				return nil, err
			}
			items := make([]components.ContainerItem, len(containers))
			for i, c := range containers {
				items[i] = components.ContainerItem{
					ID:      c.ID,
					Name:    c.Name,
					Details: c,
				}
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
