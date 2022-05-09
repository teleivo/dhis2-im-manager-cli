package instance

import (
	"encoding/json"
	"fmt"
	"strings"

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
	manager      *Manager
	list         list.Model
	stacks       []Stacks
	curStack     *Stack
	curStackJson string
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

type stacksMsg struct {
	stacks []Stacks
	items  []list.Item
}

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
		return stacksMsg{stacks: sts, items: items}
	}
}

type stackMsg struct {
	stack     *Stack
	stackJson string
}

func (m model) fetchStack(id int) tea.Cmd {
	return func() tea.Msg {
		st, err := m.manager.Stack(id)
		// TODO put into message and handle in view
		if err != nil {
			return err
		}
		stackJson, err := json.MarshalIndent(st, "", "  ")
		if err != nil {
			return err
		}
		return stackMsg{
			stack:     st,
			stackJson: string(stackJson),
		}
	}
}

// TODO load JSON when moving between list entries
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			st := m.stacks[m.list.Index()]
			return m, m.fetchStack(st.ID)
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case stacksMsg:
		m.stacks = msg.stacks
		cmd := m.list.SetItems(msg.items)
		return m, cmd
	case stackMsg:
		m.curStack = msg.stack
		m.curStackJson = msg.stackJson
		// TODO any cmd necessary?
		// TODO should curStack be cleared at any point?
		// TODO return model and nil or fall through?
		return m, nil
	}

	newList, cmd := m.list.Update(msg)
	m.list = newList
	return m, cmd
}

func (m model) View() string {
	var doc strings.Builder
	list := docStyle.Render(m.list.View())

	// TODO use a pager for the JSON
	if m.curStack != nil {
		doc.WriteString(lipgloss.JoinHorizontal(
			lipgloss.Top,
			list,
			docStyle.Render(m.curStackJson),
		))
	} else {
		doc.WriteString(list)
	}

	return doc.String()
}
