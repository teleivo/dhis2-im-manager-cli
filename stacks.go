package instance

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
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
	ready        bool
	viewport     viewport.Model
	curStack     *Stack
	curStackJson string
}

func NewStacks(im *Manager) model {
	var items []list.Item
	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.SetShowTitle(false)
	list.SetShowHelp(false)

	view := viewport.New(0, 0)
	view.KeyMap = viewport.KeyMap{
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", " ", "f"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "b"),
		),
		Up: key.NewBinding(
			key.WithKeys("u", "ctrl+u"),
		),
		Down: key.NewBinding(
			key.WithKeys("d", "ctrl+d"),
		),
	}

	return model{
		manager:  im,
		list:     list,
		viewport: view,
	}
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
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			st := m.stacks[m.list.Index()]
			return m, m.fetchStack(st.ID)
		}
	case stacksMsg:
		m.stacks = msg.stacks
		cmd := m.list.SetItems(msg.items)
		return m, cmd
	case stackMsg:
		m.curStack = msg.stack
		m.curStackJson = msg.stackJson
		m.viewport.SetContent(m.curStackJson)
		// TODO any cmd necessary?
		// TODO should curStack be cleared at any point?
		// TODO return model and nil or fall through?
		return m, nil
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		// TODO split space more "equally" between the list and the viewport
		m.list.SetSize(msg.Width-h, msg.Height-v)

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport.Width = msg.Width - h
			m.viewport.Height = msg.Height - v
			m.viewport.SetContent(m.curStackJson)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - v
		}
	}

	// Handle keyboard and mouse events
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var doc strings.Builder
	list := docStyle.Render(m.list.View())

	if m.curStackJson != "" {
		doc.WriteString(lipgloss.JoinHorizontal(
			lipgloss.Top,
			list,
			docStyle.Render(m.viewport.View()),
		))
	} else {
		doc.WriteString(list)
	}

	return doc.String()
}
