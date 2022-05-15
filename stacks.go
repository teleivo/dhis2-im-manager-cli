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

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type stacks struct {
	manager       *Manager
	ready         bool
	list          list.Model
	viewport      viewport.Model
	curIndex      int
	curStackJson  string
	stacks        []Stacks
	stacksDetails []*Stack
	stacksJson    []string
}

type selectItemMsg struct {
	index int
}

func NewStacks(im *Manager) stacks {
	d := list.NewDefaultDelegate()
	d.ShowDescription = false
	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		// get the currently selected item
		// cannot just listen to mouse down/up events as list model has not
		// been updated as stacks model receives the event before. Thus, rely
		// on a delegate which is called after the list model was updated.
		return func() tea.Msg {
			return selectItemMsg{index: m.Index()}
		}
	}

	var items []list.Item
	list := list.New(items, d, 0, 0)
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

	return stacks{
		manager:  im,
		list:     list,
		viewport: view,
		curIndex: -1,
	}
}

func (m stacks) Init() tea.Cmd {
	return m.fetchStacks()
}

type stacksMsg struct {
	stacks []Stacks
	items  []list.Item
}

func (m stacks) fetchStacks() tea.Cmd {
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

type stacksDetailsMsg struct {
	stacks     []*Stack
	stacksJson []string
}

func (m stacks) fetchStacksDetails() tea.Cmd {
	return func() tea.Msg {
		var ids []int
		for _, st := range m.stacks {
			ids = append(ids, st.ID)
		}
		sts, err := m.manager.StackDetails(ids...)
		// TODO put into message and handle in view
		if err != nil {
			return err
		}

		var stackJson []string
		for _, st := range sts {
			sj, err := json.MarshalIndent(st, "", "  ")
			if err != nil {
				return err
			}
			stackJson = append(stackJson, string(sj))
		}
		return stacksDetailsMsg{
			stacks:     sts,
			stacksJson: stackJson,
		}
	}
}

func (m stacks) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case stacksMsg:
		m.stacks = msg.stacks
		cmds = append(cmds, m.list.SetItems(msg.items))
		cmds = append(cmds, m.fetchStacksDetails())
		// TODO return or let it fall through?
		return m, tea.Batch(cmds...)
	case stacksDetailsMsg:
		m.stacksDetails = msg.stacks
		m.stacksJson = msg.stacksJson
		// TODO should I set the content here as well as in selectItemMsg?
		if m.curIndex >= 0 {
			m.curStackJson = m.stacksJson[m.curIndex]
			m.viewport.SetContent(m.curStackJson)
		}
		return m, nil
	case selectItemMsg:
		if msg.index != m.curIndex {
			m.curIndex = msg.index
			if len(m.stacksJson) > 0 {
				m.curStackJson = m.stacksJson[m.curIndex]
				m.viewport.SetContent(m.curStackJson)
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		// TODO this should not take up all the space
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

func (m stacks) View() string {
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
