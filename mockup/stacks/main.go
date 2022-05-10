package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	instance "github.com/teleivo/dhis2-im-manager-cli"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	ready    bool
	list     list.Model
	viewport viewport.Model
	curStack *instance.Stack
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
			// case "ctrl+d", "ctrl+u":
			// TODO pass on to viewport?
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

		// TODO set size of viewport
	}

	newList, cmd := m.list.Update(msg)
	cmds = append(cmds, cmd)
	m.list = newList

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var doc strings.Builder
	list := docStyle.Render(m.list.View())

	if m.curStack != nil {
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

func main() {
	items := []list.Item{
		item{title: "DHIS2 (1)"},
		item{title: "DHIS2 DB (2)"},
	}

	curStack := &instance.Stack{
		ID:   1,
		Name: "DHIS2",
		OptionalParams: []instance.OptionalParam{
			{
				ID:   12,
				Name: "LOG4J2_CONFIGURATION_FILE",
			},
			{
				ID:   12,
				Name: "LOG4J2_CONFIGURATION_FILE",
			},
			{
				ID:   12,
				Name: "LOG4J2_CONFIGURATION_FILE",
			},
			{
				ID:   12,
				Name: "LOG4J2_CONFIGURATION_FILE",
			},
			{
				ID:   12,
				Name: "LOG4J2_CONFIGURATION_FILE",
			},
		},
		RequiredParams: []instance.RequiredParam{
			{
				ID:   34,
				Name: "DB_HOST",
			},
		},
	}
	stackDetails, _ := json.MarshalIndent(curStack, "", "  ")
	view := viewport.New(40, 30)
	view.SetContent(string(stackDetails))
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

	m := model{
		curStack: curStack,
		list:     list.New(items, list.NewDefaultDelegate(), 0, 0),
		viewport: view,
	}
	m.list.Title = "Stacks"

	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Program failed: %s", err)
		os.Exit(1)
	}
}
