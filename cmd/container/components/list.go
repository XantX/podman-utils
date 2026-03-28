package components

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/podutil/podutil/internal/podman"
)

type styles struct {
	app           lipgloss.Style
	title         lipgloss.Style
	statusMessage lipgloss.Style
	success       lipgloss.Style
	error         lipgloss.Style
	help          lipgloss.Style
}

func newStyles(darkBG bool) styles {
	lightDark := lipgloss.LightDark(darkBG)

	return styles{
		app: lipgloss.NewStyle().
			Padding(1, 2),
		title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lightDark(lipgloss.Color("45"), lipgloss.Color("45"))).
			Padding(0, 1),
		statusMessage: lipgloss.NewStyle().
			Foreground(lightDark(lipgloss.Color("04B575"), lipgloss.Color("04B575"))),
		success: lipgloss.NewStyle().
			Foreground(lightDark(lipgloss.Color("46"), lipgloss.Color("46"))).
			Bold(true),
		error: lipgloss.NewStyle().
			Foreground(lightDark(lipgloss.Color("196"), lipgloss.Color("196"))).
			Bold(true),
		help: lipgloss.NewStyle().
			Foreground(lightDark(lipgloss.Color("245"), lipgloss.Color("245"))),
	}
}

// Global styles for backward compatibility with other files
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

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

type delegateKeyMap struct {
	choose key.Binding
	remove key.Binding
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
		remove: key.NewBinding(
			key.WithKeys("x", "backspace"),
			key.WithHelp("x", "delete"),
		),
	}
}

// Additional short help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.remove,
	}
}

// Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.remove,
		},
	}
}

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	// Customize selection styles to match the original colors
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		Foreground(lipgloss.Color("232")).
		Background(lipgloss.Color("107"))
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.
		Foreground(lipgloss.Color("232")).
		Background(lipgloss.Color("107"))

	// No custom UpdateFunc; key handling is done in the ListModel.Update
	d.UpdateFunc = nil

	help := []key.Binding{keys.choose, keys.remove}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

type ListModel struct {
	list           list.Model
	styles         styles
	darkBG         bool
	width, height  int
	actionName     string
	onAction       func(id string) error
	fetchItems     func() ([]ContainerItem, error)
	err            error
	successMessage string
	loaded         bool
	keys           *listKeyMap
	delegateKeys   *delegateKeyMap
}

func NewListModel(
	title string,
	actionName string,
	onAction func(id string) error,
	fetchItems func() ([]ContainerItem, error),
) *ListModel {
	m := &ListModel{
		styles:     newStyles(false),
		actionName: actionName,
		onAction:   onAction,
		fetchItems: fetchItems,
	}

	m.delegateKeys = newDelegateKeyMap()
	m.keys = newListKeyMap()

	delegate := newItemDelegate(m.delegateKeys)
	groceryList := list.New([]list.Item{}, delegate, 0, 0)
	groceryList.Title = title
	groceryList.Styles.Title = m.styles.title
	groceryList.SetShowFilter(true)
	groceryList.SetFilteringEnabled(true)
	groceryList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			m.keys.toggleSpinner,
			m.keys.toggleTitleBar,
			m.keys.toggleStatusBar,
			m.keys.togglePagination,
			m.keys.toggleHelpMenu,
		}
	}

	m.list = groceryList
	return m
}

func (m *ListModel) updateListProperties() {
	// Update list size.
	h, v := m.styles.app.GetFrameSize()
	m.list.SetSize(m.width-h, m.height-v)

	// Update the model and list styles.
	m.styles = newStyles(m.darkBG)
	m.list.Styles.Title = m.styles.title
}

func (m *ListModel) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestBackgroundColor,
		func() tea.Msg {
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
		},
	)
}

func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.successMessage != "" {
		return m, tea.Quit
	}

	if m.err != nil {
		return m, nil
	}

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		m.darkBG = msg.IsDark()
		m.updateListProperties()
		return m, nil

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.updateListProperties()
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case msg.String() == "ctrl+c" || msg.String() == "q" || msg.String() == "esc":
			return m, tea.Quit
		case key.Matches(msg, m.keys.toggleSpinner):
			cmd := m.list.ToggleSpinner()
			return m, cmd

		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil

		case key.Matches(msg, m.delegateKeys.choose):
			selected, ok := m.list.SelectedItem().(ContainerItem)
			if ok && m.onAction != nil {
				if err := m.onAction(selected.ID); err != nil {
					m.err = err
				} else {
					m.successMessage = fmt.Sprintf("%s '%s' (ID: %s) %s", m.actionName, selected.Name, selected.ID, "ejecutado exitosamente")
				}
			}
			return m, tea.Quit

		case key.Matches(msg, m.delegateKeys.remove):
			// Optionally handle remove action for containers
			// For now, just ignore or show a message
			selected, ok := m.list.SelectedItem().(ContainerItem)
			if ok {
				statusCmd := m.list.NewStatusMessage(m.styles.statusMessage.Render("Cannot delete container " + selected.Name))
				return m, statusCmd
			}
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *ListModel) View() tea.View {
	if m.successMessage != "" {
		v := tea.NewView("")
		v.SetContent("\n" + m.styles.success.Render(m.successMessage) + "\n\nPresiona q para salir.\n")
		return v
	}

	if m.err != nil {
		v := tea.NewView("")
		v.SetContent("\n" + m.styles.error.Render(fmt.Sprintf("Error: %v", m.err)) + "\n\nPresiona q para salir.\n")
		return v
	}

	if !m.loaded {
		v := tea.NewView("")
		v.SetContent("\nCargando...\n")
		return v
	}

	v := tea.NewView(m.styles.app.Render(m.list.View()))
	v.AltScreen = true
	return v
}
