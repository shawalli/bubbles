package calendar

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

// DayContentMsg enables updates for the content of a single day.
type DayContentMsg struct {
	// The day to update
	Date time.Time

	// The day model
	Content tea.Model
}

// ActiveDateMsg notifies to other models which date is set as the active date.
type ActiveDateMsg struct {
	// The day to update
	Date time.Time
}

// MonthModel represents a full calendar month.
type MonthModel struct {
	// keyMap is key bindings for calendar navigation
	keyMap KeyMap

	// startOfWeek is the day that represents the beginning of the week
	startOfWeek time.Weekday

	// weekdays manages labels for weekdays
	weekdays Weekdays

	// year of the month represented
	year int
	// month to represent
	month time.Month

	// days contains user-provided information about each day
	days map[int]tea.Model

	activeDay int

	// Styles
	styles MonthStyles
}

// NewMonth creates a new MonthModel.
func NewMonth(year int, month time.Month) MonthModel {
	m := MonthModel{
		keyMap: DefaultMonthKeyMap(),

		year:  year,
		month: month,

		startOfWeek: time.Sunday,
		weekdays:    DefaultWeekdays(),

		days:      make(map[int]tea.Model),
		activeDay: 0,

		styles: DefaultMonthStyles(),
	}

	return m
}

// StartOfWeek sets the first day of a week.
func (m MonthModel) StartOfWeek(weekday time.Weekday) MonthModel {
	m.startOfWeek = weekday
	return m
}

// Weekdays sets custom weekday labels.
func (m MonthModel) Weekdays(weekdays Weekdays) MonthModel {
	m.weekdays = weekdays
	return m
}

// Styles sets custom styling.
func (m MonthModel) Styles(styles MonthStyles) MonthModel {
	m.styles = styles
	return m
}

// Init the MonthModel.
func (m MonthModel) Init() tea.Cmd { return nil }

// Update the MonthModel.
func (m MonthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		oldActiveDay := m.activeDay
		daysInMonth := DaysInMonth(m.year, m.month)
		switch {
		case key.Matches(msg, m.keyMap.Left):
			i := m.activeDay - 1
			if i <= 0 {
				i = daysInMonth
			}
			m.activeDay = i
		case key.Matches(msg, m.keyMap.Right):
			i := m.activeDay + 1
			if i > daysInMonth {
				i = 1
			}
			m.activeDay = i
		case key.Matches(msg, m.keyMap.Up):
			i := m.activeDay - 7
			if i < 0 {
				i = m.activeDay + 28
				if i > daysInMonth {
					i = i - 7
				}
			}
			m.activeDay = i
		case key.Matches(msg, m.keyMap.Down):
			i := m.activeDay + 7
			if i > daysInMonth {
				i = m.activeDay - 28
				if m.activeDay < 28 {
					i = i + 7
				}
			}
			m.activeDay = i
		}

		if oldActiveDay != m.activeDay {
			cmds = append(cmds, func() tea.Msg {
				return ActiveDateMsg{
					Date: time.Date(m.year, m.month, m.activeDay, 0, 0, 0, 0, time.UTC),
				}
			})
		}
	case DayContentMsg:
		if msg.Date.Year() != m.year {
			break
		}
		if msg.Date.Month() != m.month {
			break
		}
		// Translate from 1-indexed date to 0-indexed array
		i := msg.Date.Day() - 1
		m.days[i] = msg.Content
	default:
		for i, d := range m.days {
			n, cmd := d.Update(msg)
			m.days[i] = n
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

// View renders the MonthModel.
func (m MonthModel) View() string {
	return gloss.JoinVertical(
		gloss.Top,
		m.ViewHeaders(),
		m.ViewWeeks(),
	)
}

// ViewHeaders renders the weekdays headers.
func (m MonthModel) ViewHeaders() string {
	var headers []string

	for i := 0; i < 7; i++ {
		wd := time.Weekday((7 + int(m.startOfWeek) + i) % 7)

		style := m.styles.MiddleHeaderStyle
		switch i {
		case 0:
			style = m.styles.LeftHeaderStyle
		case 6:
			style = m.styles.RightHeaderStyle
		}
		label, _ := m.weekdays.Get(wd)
		headers = append(headers, style.Render(label))
	}

	return gloss.JoinHorizontal(gloss.Top, headers...)
}

// ViewWeeks renders the calendar main block.
func (m MonthModel) ViewWeeks() string {
	daysInMonth := DaysInMonth(m.year, m.month)
	firstWeekdayOfMonth := FirstWeekdayOfMonth(m.year, m.month)
	weeksInMonth := CalendarRowsInMonth(m.year, m.month, m.startOfWeek)

	var weeks [][]string
	var week []string
	for i := 0; i < daysInMonth; i++ {
		wd := time.Weekday((i + int(firstWeekdayOfMonth)) % 7)
		// If it is the start of the week, but is not the first day of the month,
		// add week to calendar, reset week slice, and begin new week
		if wd == m.startOfWeek && i != 0 {
			weeks = append(weeks, week)
			week = nil
		}

		// Render day number and day body into one block of text
		body := m.styles.DateStyles.BodyStyle.Render("")
		if dayBodyModel, ok := m.days[i]; ok {
			body = m.styles.DateStyles.BodyStyle.Render(dayBodyModel.View())
		}

		lastWeek := len(weeks) == (weeksInMonth - 1)
		day := m.ViewDay(wd, i+1, body, lastWeek)

		// Add bordered day to week
		week = append(week, day)

	}

	// Pad end of month
	for i := len(week); i < 7; i++ {
		wd := time.Weekday((int(m.startOfWeek) + i) % 7)
		day := m.ViewDay(wd, 0, "", true)
		week = append(week, day)
	}
	weeks = append(weeks, week)

	// Pad start of month
	firstWeek := weeks[0]
	for i := (7 - len(firstWeek) - 1); i >= 0; i-- {
		wd := time.Weekday((int(m.startOfWeek) + i) % 7)
		firstWeek = append([]string{m.ViewDay(wd, 0, "", false)}, firstWeek...)
	}
	weeks[0] = firstWeek

	// Combine each week into a horiztonal string
	var rows []string
	for _, week := range weeks {
		rows = append(
			rows,
			gloss.JoinHorizontal(gloss.Top, week...),
		)
	}

	// Combine individual week rows together into a vertical month
	return gloss.JoinVertical(gloss.Top, rows...)
}

// ViewDay renders a single day.
//
// If zero is passed in for the day, an empty date block will be rendered.
func (m MonthModel) ViewDay(weekday time.Weekday, day int, body string, lastRow bool) string {
	num := m.styles.DateStyles.NumberStyle.Render("")
	if day > 0 {
		style := m.styles.DateStyles.NumberStyle
		if day == m.activeDay {
			style = m.styles.DateStyles.ActiveNumberStyle
		}
		num = style.Render(fmt.Sprintf("%d", day))
	}

	dateBlock := gloss.JoinVertical(
		gloss.Top,
		num,
		body,
	)

	// Figure out if the border style is left, middle, right
	// or bottom-left, bottom-middle, or bottom-right
	var style gloss.Style
	endOfWeek := time.Weekday((int(m.startOfWeek) + 6) % 7)
	switch weekday {
	case m.startOfWeek:
		style = m.styles.MiddleLeftDayStyle
		if lastRow {
			style = m.styles.BottomLeftDayStyle
		}
	case endOfWeek:
		style = m.styles.MiddleRightDayStyle
		if lastRow {
			style = m.styles.BottomRightDayStyle
		}
	default:
		style = m.styles.MiddleDayStyle
		if lastRow {
			style = m.styles.BottomDayStyle
		}
	}

	// Put border around day content and return
	return style.Render(dateBlock)
}

// Title generates a title for the calendar that may be used during rendering.
func (m MonthModel) Title(includeYear bool) string {
	d := time.Date(m.year, m.month, 1, 0, 0, 0, 0, time.UTC)

	if includeYear {
		return d.Format("January 2006")
	}
	return d.Format("January")
}

// DaysInMonth calculates the number of days in a given month and year.
func DaysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// FirstWeekdayOfMonth calculates the weekday of the first day of the month.
func FirstWeekdayOfMonth(year int, month time.Month) time.Weekday {
	d := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	return d.Weekday()
}

// CalendarRowsInMonth calculates the number of calendar rows for a month.
func CalendarRowsInMonth(year int, month time.Month, startOfWeek time.Weekday) int {
	dim := DaysInMonth(year, month)

	n := 1
	fdom := FirstWeekdayOfMonth(year, month)
	for i := 1; i <= dim; i++ {
		wd := time.Weekday(((i - 1) + int(fdom)) % 7)

		if wd == startOfWeek && i != 1 {
			n++
		}
	}

	return n
}
