package main

import (
	"encoding/json"
	"fmt"
	"log"
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

type selectItemMsg struct {
	index int
}

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
	curIndex int
	curStack string
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
		}
	case selectItemMsg:
		if msg.index != m.curIndex {
			log.Println("cur index: ", m.curIndex, "new index: ", msg.index)
			m.curIndex = msg.index
			// st := m.stacks[msg.index]
			// cmds = append(cmds, m.fetchStack(st.ID))
			m.viewport.SetContent(fmt.Sprintf("index change event: %d", msg.index))
			return m, nil
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		// TODO split space "equally" between the list and the viewport
		m.list.SetSize(msg.Width-h, msg.Height-v)

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport.Width = msg.Width - h
			m.viewport.Height = msg.Height - v
			m.viewport.SetContent(m.curStack)
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

	if m.curStack != "" {
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
	if err := run(); err != nil {
		fmt.Printf("Program failed: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	// TODO hacky since it only works running this mockup from within stacks
	// folder
	sts, err := os.ReadFile("stacks.json")
	if err != nil {
		return err
	}
	var stacks []instance.Stacks
	err = json.Unmarshal(sts, &stacks)
	if err != nil {
		return err
	}
	var items []list.Item
	for _, v := range stacks {
		items = append(items, item{title: fmt.Sprintf("%s (%d)", v.Name, v.ID)})
	}
	dl := list.NewDefaultDelegate()
	dl.ShowDescription = false
	dl.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		// get the currently selected item
		// cannot just listen to mouse down/up events as list model has not
		// been updated as stacks model receives the event before. Thus, rely
		// on a delegate which is called after the list model was updated.
		return func() tea.Msg {
			return selectItemMsg{index: m.Index()}
		}
	}
	list := list.New(items, dl, 0, 0)
	list.SetShowTitle(false)
	list.SetShowHelp(false)

	curStack, err := os.ReadFile("stack.json")
	if err != nil {
		return err
	}

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

	// TODO turn off list help
	// TODO create separate help combining list and viewport keys
	m := model{
		list:     list,
		curStack: string(curStack),
		viewport: view,
		curIndex: -1,
	}
	logfile := os.Getenv("BUBBLETEA_LOG")
	if logfile != "" {
		if _, err := tea.LogToFile(logfile, "simple"); err != nil {
			return err
		}
	}
	p := tea.NewProgram(m, tea.WithAltScreen())

	return p.Start()
}
