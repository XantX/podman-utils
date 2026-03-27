package components

import (
	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/podutil/podutil/internal/podman"
)

var (
	LabelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("107")).Bold(true)
	ValueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("254"))
)

type DetailModel struct {
	container podman.ContainerDetail
	onBack    func() tea.Model
}

func NewDetailModel(container podman.ContainerDetail, onBack func() tea.Model) *DetailModel {
	return &DetailModel{
		container: container,
		onBack:    onBack,
	}
}

func (m *DetailModel) Init() tea.Cmd {
	return nil
}

func (m *DetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc", "backspace", "b":
			return m.onBack(), nil
		}
	case tea.WindowSizeMsg:
		return m, nil
	}
	return m, nil
}

func (m *DetailModel) View() tea.View {
	s := "\n" + TitleStyle.Render("Detalles del Contenedor") + "\n\n"

	s += LabelStyle.Render("ID: ") + ValueStyle.Render(m.container.ID) + "\n"
	s += LabelStyle.Render("Nombre: ") + ValueStyle.Render(m.container.Name) + "\n"
	s += LabelStyle.Render("Imagen: ") + ValueStyle.Render(m.container.Image) + "\n"
	s += LabelStyle.Render("Estado: ") + ValueStyle.Render(m.container.State) + "\n"
	s += LabelStyle.Render("Status: ") + ValueStyle.Render(m.container.Status) + "\n"
	s += LabelStyle.Render("Puertos: ") + ValueStyle.Render(m.container.Ports) + "\n"
	s += LabelStyle.Render("Creado: ") + ValueStyle.Render(m.container.Created) + "\n"
	s += LabelStyle.Render("Comando: ") + ValueStyle.Render(m.container.Command) + "\n"

	s += "\n" + HelpStyle.Render("Presiona q/esc/backspace para volver\n")

	v := tea.NewView(s)
	return v
}
