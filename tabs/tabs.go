package tabs

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

// TabSizeMsg is a message about the dimensions of the tabs.Model.
// It is used by a [Tab] to request a resizing of [Model] and all other [Tab] windows.
type TabSizeMsg struct {
	// Width is the tab's total width.
	Width int

	// Height is the tab's total height.
	Height int
}

// Model represents a group of tab headers and their content.
type Model struct {
	// keyMap is key bindings for tab navigation
	keyMap KeyMap

	// tabs to render
	tabs []Tab

	// activeTab is the index of the tab currently being presented
	activeTab int

	// width of all tabs in the tab group
	width int

	// height of all tabs in the tab group
	height int

	// wraparound enables wrapping around from the last tab to the first tab or from the first tab
	// back to the last tab
	wraparound bool

	// styles for the tab headers and bodies
	styles Styles
}

// New creates a new Model.
func New(tabs ...Tab) Model {
	m := Model{
		keyMap: DefaultKeyMap(len(tabs)),
		tabs:   tabs,
		styles: DefaultStyles(),
	}
	m = m.DefaultDimensions()

	return m
}

// Wraparound enables wraparound navigation from the last tab to the first tab and from the
// first tab back to the last tab.
func (m Model) Wraparound() Model {
	m.wraparound = true
	return m
}

// Styles enables custom styling.
func (m Model) Styles(styles Styles) Model {
	m.styles = styles
	return m
}

// DefaultDimensions applies the current terminal's dimensions to the tab group.
func (m Model) DefaultDimensions() Model {
	termWidth, termHieght, _ := term.GetSize(os.Stdout.Fd())

	// Pad from the right edge of the screen
	m.width = termWidth - 2
	m.height = termHieght

	return m
}

// Width sets the Model width.
func (m Model) Width(w int) Model {
	m.width = w
	return m
}

// Height sets the Model height.
func (m Model) Height(h int) Model {
	m.height = h
	return m
}

// width is a convenience function to calculate the width of the longest line.
func width(s string) int {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")

	var w int
	for _, line := range lines {
		w = max(w, gloss.Width(line))
	}

	return w
}

// height is a convenience function to calculate the height all lines.
func height(s string) int {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return gloss.Height(s)
}

// SetTab sets the active tab.
func (m Model) SetTab(i int) Model {
	if i < 0 {
		i = 0
	} else if i >= len(m.tabs) {
		i = len(m.tabs) - 1
	}
	m.activeTab = i

	return m
}

// NextTab moves the active tab forward.
func (m Model) NextTab() Model {
	i := m.activeTab + 1

	if i == len(m.tabs) {
		if m.wraparound {
			i = 0
		} else {
			i--
		}
	}

	return m.SetTab(i)
}

// PreviousTab moves the active tab backward.
func (m Model) PreviousTab() Model {
	i := m.activeTab - 1

	if i < 0 {
		if m.wraparound {
			i = len(m.tabs) - 1
		} else {
			i = 0
		}
	}

	return m.SetTab(i)
}

// Init the Model.
func (m Model) Init() tea.Cmd { return nil }

// Update the Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.TabLeft):
			m = m.PreviousTab()
		case key.Matches(msg, m.keyMap.TabRight):
			m = m.NextTab()
		default:
			for i, kb := range m.keyMap.TabNumbers {
				if key.Matches(msg, kb) {
					m = m.SetTab(i)
					break
				}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width - 2
		m.height = msg.Height
		msg.Width = m.width
		cmds = append(cmds, tea.ClearScreen)
	case TabSizeMsg:
		m.width = msg.Width - 2
		m.height = msg.Height
		cmds = append(cmds, tea.ClearScreen)
	}

	// Update and pass messages to all children
	for i := range m.tabs {
		child, cmd := m.tabs[i].child.Update(msg)
		m.tabs[i].child = child
		cmds = append(cmds, cmd)
	}

	return m, tea.Sequence(cmds...)
}

// View renders the Model.
func (m Model) View() string {
	doc := strings.Builder{}

	// Tabs
	tabs := m.ViewTabs()
	doc.WriteString(tabs + "\n")

	// Window
	w := width(tabs) - (m.styles.TabSpacer.GetPaddingLeft() +
		m.styles.TabSpacer.GetPaddingRight() +
		m.styles.TabWindow.GetPaddingLeft() +
		m.styles.TabWindow.GetPaddingRight())
	window := m.ViewWindow(w)
	doc.WriteString(window + "\n")

	return doc.String()
}

// ViewTabs renders tab headers.
func (m Model) ViewTabs() string {
	// Render all tabs
	var tabs []string
	for i, tab := range m.tabs {
		style := m.styles.Tab
		row := tab.title
		if isActive := (i == m.activeTab); isActive {
			style = m.styles.ActiveTab
			row = fmt.Sprintf("%s %s %s",
				m.styles.TabIndicator.Render(m.styles.TabIndicatorLeft),
				row,
				m.styles.TabIndicator.Render(m.styles.TabIndicatorRight),
			)
			if isFirst := (i == 0); isFirst {
				border, tb, rb, bb, lb := style.GetBorder()
				border.BottomLeft = "│"
				style = style.Border(border, tb, rb, bb, lb)
			}
			row = style.Render(row)
		} else {
			row = fmt.Sprintf("%s %s %s",
				strings.Repeat(" ", width(m.styles.TabIndicatorLeft)),
				row,
				strings.Repeat(" ", width(m.styles.TabIndicatorRight)),
			)
			if isFirst := (i == 0); isFirst {
				border, tb, rb, bb, lb := style.GetBorder()
				border.BottomLeft = "├"
				style = style.Border(border, tb, rb, bb, lb)
			}
			row = style.Render(row)
		}
		tabs = append(tabs, row)
	}
	row := gloss.JoinHorizontal(gloss.Top, tabs...)

	// Render tab spacer (along top right)
	//   Need at least one rune to get the bottom-right border to print
	tabSpacerSize := max(1, m.width-width(row)-1)
	row = gloss.JoinHorizontal(
		gloss.Bottom,
		row,
		m.styles.TabSpacer.Padding(0, 1).Render(strings.Repeat(" ", tabSpacerSize)),
	)

	return row
}

// ViewWindow renders the tab window.
func (m Model) ViewWindow(w int) string {
	content := m.tabs[m.activeTab].child.View()

	var padded []string
	for _, l := range strings.Split(content, "\n") {
		lenL := width(l)
		if lenL > w {
			l1 := l[:w-1] + "\r"
			l2 := l[w:lenL]
			l2 = l2 + strings.Repeat(" ", w-width(l2)) + "\r"
			padded = append(padded, l1, l2)
		} else {
			p := l + strings.Repeat(" ", w-lenL)
			padded = append(padded, p)
		}
	}
	content = strings.Join(padded, "\n")

	return m.styles.TabWindow.Render(content)
}
