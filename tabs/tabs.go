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

type Model struct {
	keyMap KeyMap

	tabs      []Tab
	activeTab int

	width  int
	height int

	wraparound bool

	styles Styles
}

func New(tabs ...Tab) Model {
	m := Model{
		keyMap: DefaultKeyMap(len(tabs)),

		tabs: tabs,

		styles: DefaultStyles(),
	}
	m = m.DefaultDimensions()

	return m
}

func (m Model) Wraparound() Model {
	m.wraparound = true
	return m
}

func (m Model) Styles(styles Styles) Model {
	m.styles = styles
	return m
}

func (m Model) DefaultDimensions() Model {
	termWidth, termHieght, _ := term.GetSize(os.Stdout.Fd())

	m.width = termWidth - 2
	m.height = termHieght

	return m
}

func (m Model) Width(w int) Model {
	m.width = w
	return m
}

func (m Model) Height(h int) Model {
	m.height = h
	return m
}

func width(s string) int {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")

	var w int
	for _, line := range lines {
		w = max(w, gloss.Width(line))
	}

	return w
}

func height(s string) int {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return gloss.Height(s)
}

func (m Model) SetTab(i int) Model {
	if i < 0 {
		i = 0
	} else if i >= len(m.tabs) {
		i = len(m.tabs) - 1
	}
	m.activeTab = i

	return m
}

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

func (m Model) Init() tea.Cmd { return nil }

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
	}

	// Update children
	child, cmd := m.tabs[m.activeTab].child.Update(msg)
	m.tabs[m.activeTab].child = child

	cmds = append(cmds, cmd)

	return m, tea.Sequence(cmds...)
}

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

func (m Model) ViewTabs() string {
	// Render all tabs
	var tabs []string
	for i, tab := range m.tabs {
		style := m.styles.Tab
		row := tab.title
		if isActive := (i == m.activeTab); isActive {
			style = m.styles.ActiveTab
			row = fmt.Sprintf("%s %s %s",
				m.styles.TabIndicator.Render("╾"),
				row,
				m.styles.TabIndicator.Render("╼"),
			)
			if isFirst := (i == 0); isFirst {
				border, tb, rb, bb, lb := style.GetBorder()
				border.BottomLeft = "│"
				style = style.Border(border, tb, rb, bb, lb)
			}
			row = style.Render(row)
		} else {
			row = fmt.Sprintf("  %s  ", row)
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
