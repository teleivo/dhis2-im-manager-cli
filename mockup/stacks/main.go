package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
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
	list     list.Model
	curStack *instance.Stack
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	newList, cmd := m.list.Update(msg)
	m.list = newList
	return m, cmd
}

func (m model) View() string {
	var doc strings.Builder
	list := docStyle.Render(m.list.View())
	curStack, _ := json.MarshalIndent(m.curStack, "", "  ")

	if m.curStack != nil {
		doc.WriteString(lipgloss.JoinHorizontal(
			lipgloss.Top,
			list,
			docStyle.Render(string(curStack)),
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
	m := model{
		curStack: &instance.Stack{
			ID:   1,
			Name: "DHIS2",
			OptionalParams: []instance.OptionalParam{
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
		},
		list: list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	m.list.Title = "Stacks"

	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Program failed: %s", err)
		os.Exit(1)
	}
}
