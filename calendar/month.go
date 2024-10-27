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

		days: make(map[int]tea.Model),

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
	firstVisibleWeekday := m.weekdays.First(time.Date(m.year, m.month, 1, 0, 0, 0, 0, time.UTC))

	inititalizeActiveDay := func() bool {
		if m.activeDay != 0 {
			return false
		}
		// If uninitialized, set the active day as the first visible day in the month
		// so that cursor works as expected
		m.activeDay = 1 + int(firstVisibleWeekday)

		return true
	}

	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		oldActiveDay := m.activeDay
		daysInMonth := DaysInMonth(m.year, m.month)
		switch {
		case key.Matches(msg, m.keyMap.Left):
			inititalizeActiveDay()

			ad := time.Date(m.year, m.month, m.activeDay, 0, 0, 0, 0, time.UTC)
			diff := (7 + int(ad.Weekday()) - int(m.weekdays.Previous(ad))) % 7
			i := m.activeDay - diff
			if i <= 0 {
				i = daysInMonth
			}
			m.activeDay = i
		case key.Matches(msg, m.keyMap.Right):
			// If initializing the active day, assume the intent of the right event was to move down
			// into the first visible day of the month.
			if inititalizeActiveDay() {
				break
			}

			ad := time.Date(m.year, m.month, m.activeDay, 0, 0, 0, 0, time.UTC)
			diff := (7 + int(m.weekdays.Next(ad)) - int(ad.Weekday())) % 7
			i := m.activeDay + diff
			if i > daysInMonth {
				i = 1
			}
			m.activeDay = i
		case key.Matches(msg, m.keyMap.Up):
			inititalizeActiveDay()

			// No need to calculate visiblity because we can assume that the same weekday in the week prior to the
			// current week will also be visible
			i := m.activeDay - 7
			if i < 0 {
				i = m.activeDay + 28
				if i > daysInMonth {
					i = i - 7
				}
			}
			m.activeDay = i
		case key.Matches(msg, m.keyMap.Down):
			// If initializing the active day, assume the intent of the down event was to move down
			// into the first visible day of the month.
			if inititalizeActiveDay() {
				break
			}

			// No need to calculate visiblity because we can assume that the same weekday in the week following the
			// current week will also be visible
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

// ViewHeaders renders the weekday headers.
func (m MonthModel) ViewHeaders() string {
	startDate := m.StartOfFirstWeek()
	first := m.weekdays.First(startDate)
	last := m.weekdays.Last(startDate)

	var headers []string
	for i := 0; i < 7; i++ {
		day := startDate.AddDate(0, 0, i)
		label, ok := m.weekdays.Get(day.Weekday())
		if !ok {
			continue
		}
		style := m.styles.MiddleHeaderStyle
		switch day.Weekday() {
		case first:
			style = m.styles.LeftHeaderStyle
		case last:
			style = m.styles.RightHeaderStyle
		}
		style = style.Width(m.styles.DateStyles.Width)

		headers = append(headers, style.Render(label))
	}

	return gloss.JoinHorizontal(gloss.Top, headers...)
}

// ViewWeeks renders the calendar main block.
func (m MonthModel) ViewWeeks() string {
	daysInMonth := DaysInMonth(m.year, m.month)
	firstWeekdayOfMonth := FirstWeekdayOfMonth(m.year, m.month)
	weeksInMonth := CalendarRowsInMonth(m.year, m.month, m.startOfWeek)

	// If the first week starts in the previous month and the last visible day of that week is still in the previous
	// month, remove it from the number of weeks for the purpose of rendering the calendar
	calendarStartDate := m.StartOfFirstWeek()
	spread := m.weekdays.Spread(calendarStartDate)
	if calendarStartDate.AddDate(0, 0, spread).Month() < m.month {
		weeksInMonth -= 1
	}

	firstVisibleWeekday := m.weekdays.First(calendarStartDate)

	var weeks [][]string
	var week []string
	for i := 0; i < daysInMonth; i++ {
		date := time.Date(m.year, m.month, (i + 1), 0, 0, 0, 0, time.UTC)
		if !m.weekdays.IsVisible(date.Weekday()) {
			continue
		}

		// If beginning a new week, only add it if the week has content
		wd := time.Weekday((i + int(firstWeekdayOfMonth)) % 7)
		if (wd == firstVisibleWeekday) && (len(week) != 0) {
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
	padDay := time.Date(m.year, m.month+1, 0, 0, 0, 0, 0, time.UTC)
	for len(week) < len(m.weekdays) {
		padDay = padDay.AddDate(0, 0, 1)

		if !m.weekdays.IsVisible(padDay.Weekday()) {
			continue
		}
		day := m.ViewDay(padDay.Weekday(), 0, "", true)
		week = append(week, day)
	}
	weeks = append(weeks, week)

	// Pad start of month
	var pad []string
	padDay = m.StartOfFirstWeek().AddDate(0, 0, -1)
	week = weeks[0]
	for (len(week) + len(pad)) < len(m.weekdays) {
		padDay = padDay.AddDate(0, 0, 1)

		if !m.weekdays.IsVisible(padDay.Weekday()) {
			continue
		}
		day := m.ViewDay(padDay.Weekday(), 0, "", false)
		pad = append(pad, day)
	}
	weeks[0] = append(pad, week...)

	// Combine each week into a horizontal string
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

	startDate := m.StartOfFirstWeek()
	firstVisibleWeekday := m.weekdays.First(startDate)
	lastVisibleWeekday := m.weekdays.Last(startDate)
	switch weekday {
	case firstVisibleWeekday:
		style = m.styles.MiddleLeftDayStyle
		if lastRow {
			style = m.styles.BottomLeftDayStyle
		}
	case lastVisibleWeekday:
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
	style = style.Width(m.styles.DateStyles.Width)

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

// StartOfFirstWeek calculates the first day of the first full week of the month.
func (m MonthModel) StartOfFirstWeek() time.Time {
	startDate := time.Date(m.year, m.month, 1, 0, 0, 0, 0, time.UTC)
	diff := int(startDate.Weekday()) - int(m.startOfWeek)
	if diff < 0 {
		diff = diff + 7
	}
	return startDate.AddDate(0, 0, (-1 * diff))
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
