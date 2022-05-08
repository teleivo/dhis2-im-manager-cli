package instance

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	manager *Manager
	list    list.Model
}

func NewStacks(im *Manager) model {
	var items []list.Item
	m := model{
		manager: im,
		list:    list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	m.list.Title = "Stacks"
	return m
}

func (m model) Init() tea.Cmd {
	return m.fetchStacks()
}

type stacksMsg []list.Item

func (m model) fetchStacks() tea.Cmd {
	return func() tea.Msg {
		sts, err := m.manager.Stacks()
		if err != nil {
			return err
		}
		var items []list.Item
		for _, st := range sts {
			items = append(items, item{title: fmt.Sprintf("%s (%d)", st.Name, st.ID)})
		}
		return stacksMsg(items)
	}
}

func (m model) fetchStack(id int) tea.Cmd {
	return func() tea.Msg {
		st, err := m.manager.Stack(id)
		if err != nil {
			return err
		}
		// TODO what type of msg do I need for a pager?
		_ = st
		return nil
	}
}

// TODO show details in pager on the right
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case stacksMsg:
		cmd := m.list.SetItems(msg)
		return m, cmd
	}

	newList, cmd := m.list.Update(msg)
	m.list = newList
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}
