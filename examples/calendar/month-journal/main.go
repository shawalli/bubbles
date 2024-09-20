package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/shawalli/bubbles/calendar"
)

type Log struct {
	timestamp time.Time

	exercise string
}

type Model struct {
	calendar tea.Model

	activeDate time.Time

	log map[string][]Log
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			return m, tea.Quit
		}
	case calendar.ActiveDateMsg:
		m.activeDate = msg.Date
	}

	n, cmd := m.calendar.Update(msg)
	m.calendar = n

	return m, cmd
}

func (m Model) viewLog() string {
	width := gloss.NewStyle().Width(24)

	header := m.activeDate.Format("January 2, 2006")
	if m.activeDate == (time.Time{}) {
		// Suppress initial-state active-date
		header = ""
	}
	header = width.AlignHorizontal(gloss.Center).Bold(true).Render(header)

	key := m.activeDate.Format("2006-01-02")
	entries := m.log[key]
	slices.SortStableFunc(entries, func(a, b Log) int { return a.timestamp.Compare(b.timestamp) })

	var lines []string
	for _, e := range entries {
		l := width.Italic(true).Render(fmt.Sprintf("• %s: %s\n\n", e.timestamp.Format("3:04 PM"), e.exercise))
		lines = append(lines, l)
	}
	lines = append([]string{"\n", header, "\n\n"}, lines...)

	return gloss.JoinVertical(gloss.Top, lines...)
}

func (m Model) View() string {
	window := gloss.JoinHorizontal(
		gloss.Top,
		m.calendar.View(),
		strings.Repeat(" ", 6),
		m.viewLog(),
	)

	return gloss.NewStyle().
		Border(gloss.RoundedBorder(), true).
		Render(window)
}

type contentModel struct {
	numExercises int
}

func (m contentModel) Init() tea.Cmd                           { return nil }
func (m contentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m contentModel) View() string {
	style := gloss.NewStyle()
	switch m.numExercises {
	case 1:
		style = style.Foreground(gloss.Color("#FFA200"))
	case 2:
		style = style.Foreground(gloss.Color("#003CFF"))
	default:
		style = style.Foreground(gloss.Color("#22C11D"))
	}
	return style.Render(strings.Repeat("•", m.numExercises))
}

func getDemoLogs() map[string][]Log {
	logs := make(map[string][]Log)

	data := map[string]string{
		"2024-09-01T08:04:05Z": "walk",
		"2024-09-01T18:24:56Z": "bike",
		"2024-09-04T07:43:22Z": "bike",
		"2024-09-05T06:15:30Z": "run",
		"2024-09-09T09:58:01Z": "weights",
		"2024-09-09T11:07:20Z": "run",
		"2024-09-09T19:45:12Z": "walk",
		"2024-09-15T21:02:54Z": "weights",
		"2024-09-18T04:30:32Z": "swim",
		"2024-09-18T16:53:38Z": "walk",
		"2024-09-20T08:18:08Z": "run",
		"2024-09-25T06:45:34Z": "swim",
		"2024-09-25T12:00:44Z": "walk",
		"2024-09-25T17:32:32Z": "walk",
		"2024-09-28T10:02:46Z": "run",
		"2024-09-28T10:22:23Z": "weights",
		"2024-09-28T18:27:12Z": "walk",
	}

	for ts, d := range data {
		t, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			panic(fmt.Sprintf("unexpected invalid timestamp %q!!!", ts))
		}

		key := t.Format("2006-01-02")

		daysLogs := logs[key]

		daysLogs = append(daysLogs, Log{
			timestamp: t,
			exercise:  d,
		})
		logs[key] = daysLogs
	}

	return logs
}

func main() {
	m := Model{
		calendar: calendar.NewMonth(2024, time.September),
		log:      getDemoLogs(),
	}

	for ts, l := range m.log {
		d, _ := time.Parse("2006-01-02", ts)
		n, _ := m.Update(calendar.DayContentMsg{
			Date:    d,
			Content: contentModel{numExercises: len(l)},
		})
		m = n.(Model)
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("could not run program: %v", err)
		os.Exit(1)
	}
}
