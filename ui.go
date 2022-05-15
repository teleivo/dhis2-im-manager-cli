package instance

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const (
	// In real life situations we'd adjust the document to fit the width we've
	// detected. In the case of this example we're hardcoding the width, and
	// later using the detected width only to truncate in order to avoid jaggy
	// wrapping.
	width = 96
)

// Style definitions.
var (
	// General.
	// docStyle = lipgloss.NewStyle().Margin(1, 2)
	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)

	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	// Tabs.

	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	tab = lipgloss.NewStyle().
		Border(tabBorder, true).
		BorderForeground(highlight).
		Padding(0, 1)

	activeTab = tab.Copy().Border(activeTabBorder, true)

	tabGap = tab.Copy().
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)

	// Status Bar.

	statusNugget = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 1).
			MarginRight(1)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle)

	managerUrlStyle = statusNugget.Copy().Background(lipgloss.Color("#6124DF"))
)

type UI struct {
	manager   *Manager
	component tea.Model
}

func NewUI(im *Manager, component tea.Model) tea.Model {
	return &UI{
		manager:   im,
		component: component,
	}
}

func (ui UI) Init() tea.Cmd {
	return ui.component.Init()
}

func (ui UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO use WindowSizeMsg to set width and all

	var cmd tea.Cmd
	ui.component, cmd = ui.component.Update(msg)
	return ui, cmd
}

func (ui UI) View() string {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	doc := strings.Builder{}

	// Tabs
	{
		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			activeTab.Render("Stacks"),
			tab.Render("Instances"),
		)
		gap := tabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(row)-2)))
		row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
		doc.WriteString(row + "\n\n")
	}

	// Current Component
	{
		doc.WriteString(ui.component.View())
	}

	// Status bar
	{
		w := lipgloss.Width

		// TODO fill in real data
		auth := statusStyle.Render("user@some.com")
		managerUrl := managerUrlStyle.Render("@ instance.test.com")
		statusVal := statusText.Copy().
			Width(width - w(auth) - w(managerUrl)).
			Render("Ravishing")

		bar := lipgloss.JoinHorizontal(lipgloss.Top,
			auth,
			statusVal,
			managerUrl,
		)

		doc.WriteString(statusBarStyle.Width(width).Render(bar))
	}

	if physicalWidth > 0 {
		docStyle = docStyle.MaxWidth(physicalWidth)
	}

	return docStyle.Render(doc.String())
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
