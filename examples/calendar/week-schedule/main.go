package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/shawalli/bubbles/calendar"
)

type Appointment struct {
	StartTime time.Time
	EndTime   time.Time

	Title string

	Completed bool
}

type dayModel struct {
	Date time.Time

	Appointments []Appointment
}

func (m dayModel) Init() tea.Cmd                           { return nil }
func (m dayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m dayModel) View() string {
	slices.SortStableFunc(m.Appointments, func(a, b Appointment) int { return a.StartTime.Compare(b.EndTime) })

	var appts []string

	for _, a := range m.Appointments {
		style := gloss.NewStyle().Align(gloss.Left).Width(10)
		if a.Completed {
			style = style.
				Italic(true).
				Foreground(gloss.Color("#FFA200"))
		}

		appt := strings.Join([]string{
			style.Render(fmt.Sprintf("%s -", a.StartTime.Format("3:05PM"))),
			style.Render(fmt.Sprintf("  %s", a.EndTime.Format("3:05PM"))),
			style.Render(a.Title),
		}, "\n")

		appts = append(appts, appt)
	}

	return gloss.JoinVertical(gloss.Top, strings.Join(appts, "\n\n"))
}

type Model struct {
	schedule tea.Model

	activeDate time.Time
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

	n, cmd := m.schedule.Update(msg)
	m.schedule = n

	return m, cmd
}

func (m Model) View() string {
	return m.schedule.View()
}

func getDemoAppointments() map[time.Time][]Appointment {
	data := []map[string]interface{}{
		{
			"startTime": "2024-09-30T09:00:00Z",
			"endTime":   "2024-09-30T10:30:00Z",
			"title":     "Core Hours",
			"completed": true,
		},
		{
			"startTime": "2024-09-30T11:15:00Z",
			"endTime":   "2024-09-30T11:45:00Z",
			"title":     "Standup",
			"completed": true,
		},
		{
			"startTime": "2024-09-30T12:00:00Z",
			"endTime":   "2024-09-30T13:00:00Z",
			"title":     "Lunch",
			"completed": true,
		},
		{
			"startTime": "2024-09-30T13:05:00Z",
			"endTime":   "2024-09-30T14:00:00Z",
			"title":     "Mid-Level Weekly",
			"completed": true,
		},
		{
			"startTime": "2024-09-30T14:00:00Z",
			"endTime":   "2024-09-30T15:30:00Z",
			"title":     "Core Hours",
			"completed": true,
		},
		{
			"startTime": "2024-10-01T09:00:00Z",
			"endTime":   "2024-10-01T10:30:00Z",
			"title":     "Core Hours",
			"completed": true,
		},
		{
			"startTime": "2024-10-01T12:00:00Z",
			"endTime":   "2024-10-01T13:00:00Z",
			"title":     "Lunch",
			"completed": true,
		},
		{
			"startTime": "2024-10-01T14:00:00Z",
			"endTime":   "2024-10-01T15:30:00Z",
			"title":     "Core Hours",
		},
		{
			"startTime": "2024-10-02T09:00:00Z",
			"endTime":   "2024-10-02T10:30:00Z",
			"title":     "Core Hours",
		},
		{
			"startTime": "2024-10-02T10:30:00Z",
			"endTime":   "2024-10-02T11:10:00Z",
			"title":     "Refinement",
		},
		{
			"startTime": "2024-10-02T11:15:00Z",
			"endTime":   "2024-10-02T11:45:00Z",
			"title":     "Standup",
		},
		{
			"startTime": "2024-10-02T12:00:00Z",
			"endTime":   "2024-10-02T13:00:00Z",
			"title":     "Lunch",
		},
		{
			"startTime": "2024-10-02T14:00:00Z",
			"endTime":   "2024-10-02T15:30:00Z",
			"title":     "Core Hours",
		},
		{
			"startTime": "2024-10-02T15:00:00Z",
			"endTime":   "2024-10-02T16:00:00Z",
			"title":     "Metrics",
		},
		{
			"startTime": "2024-10-03T09:00:00Z",
			"endTime":   "2024-10-03T10:30:00Z",
			"title":     "Core Hours",
		},
		{
			"startTime": "2024-10-03T12:00:00Z",
			"endTime":   "2024-10-03T13:00:00Z",
			"title":     "Lunch",
		},
		{
			"startTime": "2024-10-03T13:00:00Z",
			"endTime":   "2024-10-03T13:50:00Z",
			"title":     "Product Alignment",
		},
		{
			"startTime": "2024-10-03T14:00:00Z",
			"endTime":   "2024-10-03T15:30:00Z",
			"title":     "Core Hours",
		},
		{
			"startTime": "2024-10-04T09:00:00Z",
			"endTime":   "2024-10-04T17:00:00Z",
			"title":     "No Meeting Day",
		},
	}

	days := make(map[time.Time][]Appointment)
	for _, d := range data {
		b, err := json.Marshal(&d)
		if err != nil {
			panic(fmt.Sprintf("unexpected error marshalling demo data: %v", err))
		}

		var appt Appointment
		if err := json.Unmarshal(b, &appt); err != nil {
			panic(fmt.Sprintf("unexpected error unmarshalling demo data: %v", err))
		}

		key := time.Date(appt.StartTime.Year(), appt.StartTime.Month(), appt.StartTime.Day(), 0, 0, 0, 0, time.UTC)
		day := days[key]
		day = append(day, appt)
		days[key] = day
	}

	return days
}

func main() {
	m := Model{
		schedule: calendar.NewWeek(
			time.Date(2024, time.October, 1, 0, 0, 0, 0, time.UTC),
		).Weekdays(calendar.Weekdays{
			time.Monday:    "Mon ",
			time.Tuesday:   "Tue",
			time.Wednesday: "Wed",
			time.Thursday:  "Thu",
			time.Friday:    "Fri",
		}),
	}

	for ts, d := range getDemoAppointments() {
		n, _ := m.Update(calendar.DayContentMsg{
			Date: ts,
			Content: dayModel{
				Date:         ts,
				Appointments: d,
			},
		})
		m = n.(Model)
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("could not run program: %v", err)
		os.Exit(1)
	}
}
