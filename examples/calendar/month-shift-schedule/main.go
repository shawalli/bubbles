package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/shawalli/bubbles/calendar"
)

type Shift struct {
	timestamp time.Time

	employee string
}

type Model struct {
	calendar tea.Model

	activeDate time.Time

	log map[string]Shift
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

func (m Model) View() string {
	titleStyle := gloss.NewStyle().
		Bold(true).
		Underline(true).
		Foreground(gloss.Color("#22C11D"))
	window := gloss.JoinVertical(
		gloss.Center,
		titleStyle.Render(m.calendar.(calendar.MonthModel).Title(true)),
		"",
		m.calendar.View(),
	)
	return gloss.NewStyle().
		Border(gloss.RoundedBorder(), true).
		Render(window)
}

type contentModel struct {
	employee string
}

func (m contentModel) Init() tea.Cmd                           { return nil }
func (m contentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m contentModel) View() string {
	return m.employee
}

func getDemoShifts() map[string]Shift {
	data := map[string]string{
		"2024-09-01": "Alice",
		"2024-09-05": "Carol",
		"2024-09-06": "Carol",
		"2024-09-07": "Bob",
		"2024-09-08": "Ted",
		"2024-09-12": "Alice",
		"2024-09-13": "Carol",
		"2024-09-14": "Carol",
		"2024-09-20": "Carol",
		"2024-09-21": "Ted",
		"2024-09-22": "Bob",
		"2024-09-26": "Bob",
		"2024-09-27": "Alice",
		"2024-09-28": "Carol",
		"2024-09-29": "Ted",
	}

	logs := make(map[string]Shift)
	for ts, d := range data {
		t, err := time.Parse("2006-01-02", ts)
		if err != nil {
			panic(fmt.Sprintf("unexpected invalid timestamp %q!!!", ts))
		}

		key := t.Format("2006-01-02")

		logs[key] = Shift{
			timestamp: t,
			employee:  d,
		}
	}

	return logs
}

func main() {
	s := calendar.DefaultMonthStyles()
	s.DateStyles.Width = 15
	s.DateStyles.BodyStyle = s.DateStyles.BodyStyle.Width(15)

	m := Model{
		calendar: calendar.NewMonth(2024, time.September).
			StartOfWeek(time.Thursday).
			Weekdays(calendar.Weekdays{
				time.Thursday: "Thursday",
				time.Friday:   "Friday",
				time.Saturday: "Saturday",
				time.Sunday:   "Sunday",
			}).
			Styles(s),
		log: getDemoShifts(),
	}

	for _, shift := range m.log {
		n, _ := m.Update(calendar.DayContentMsg{
			Date:    shift.timestamp,
			Content: contentModel{employee: shift.employee},
		})
		m = n.(Model)
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("could not run program: %v", err)
		os.Exit(1)
	}
}
