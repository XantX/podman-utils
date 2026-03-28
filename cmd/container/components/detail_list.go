package components

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

type DetailListModel struct {
	list          list.Model
	styles        styles
	darkBG        bool
	width, height int
	fetchItems    func() ([]ContainerItem, error)
	OnDetail      func(item ContainerItem) tea.Model
	err           error
	loaded        bool
	keys          *listKeyMap
	delegateKeys  *delegateKeyMap
}

func NewDetailListModel(
	title string,
	onDetail func(item ContainerItem) tea.Model,
	fetchItems func() ([]ContainerItem, error),
) *DetailListModel {
	m := &DetailListModel{
		styles:     newStyles(false),
		OnDetail:   onDetail,
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

func (m *DetailListModel) updateListProperties() {
	// Update list size.
	h, v := m.styles.app.GetFrameSize()
	m.list.SetSize(m.width-h, m.height-v)

	// Update the model and list styles.
	m.styles = newStyles(m.darkBG)
	m.list.Styles.Title = m.styles.title
}

func (m *DetailListModel) Init() tea.Cmd {
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

func (m *DetailListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if ok && m.OnDetail != nil {
				return m.OnDetail(selected), nil
			}

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

func (m *DetailListModel) View() tea.View {
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
