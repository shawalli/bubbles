package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/shawalli/bubbles/calendar"
)

type Model struct {
	activeMonth int

	months []tea.Model
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.activeMonth = (m.activeMonth + 1) % 2
		}
	}

	n, cmd := m.months[m.activeMonth].Update(msg)
	m.months[m.activeMonth] = n

	return m, cmd
}

func (m Model) View() string {
	padding := gloss.NewStyle().Padding(2, 3)
	var months []string
	for i, n := range m.months {
		month := n.View()

		block := gloss.JoinVertical(
			gloss.Top,
			gloss.NewStyle().
				AlignHorizontal(gloss.Center).
				Width(gloss.Width(month)).
				Bold(true).
				Render(n.(calendar.MonthModel).Title(true)),
			"",
			n.View(),
		)

		border := gloss.HiddenBorder()
		if i == m.activeMonth {
			border = gloss.DoubleBorder()
		}
		block = gloss.NewStyle().Border(border, true).Render(
			padding.Render(block),
		)

		months = append(months, block)
	}

	return gloss.JoinHorizontal(gloss.Top, months...)
}

type contentModel struct {
	content string
}

func (m contentModel) Init() tea.Cmd                           { return nil }
func (m contentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m contentModel) View() string                            { return m.content }

func main() {
	now := time.Now()
	nextMonthFromNow := now.AddDate(0, 1, 0)

	// Manage a calendar for this month and one for next month
	m := Model{
		months: []tea.Model{
			calendar.NewMonth(now.Year(), now.Month()),
			calendar.NewMonth(nextMonthFromNow.Year(), nextMonthFromNow.Month()),
		},
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("could not run program: %v", err)
		os.Exit(1)
	}
}
