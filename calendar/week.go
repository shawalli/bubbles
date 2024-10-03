package calendar

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

// Weekdays maps weekdays to their labels.
type Weekdays map[time.Weekday]string

// Get the label for a weekday if it exists in the map.
func (w Weekdays) Get(day time.Weekday) (string, bool) {
	s, ok := w[day]

	return s, ok
}

// IsVisible determines whether the weekday has a label, and should therefore be visible.
func (w Weekdays) IsVisible(weekday time.Weekday) bool {
	_, ok := w.Get(weekday)
	return ok
}

// First returns the first visible weekday based on the start date.
func (w Weekdays) First(startDate time.Time) time.Weekday {
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	for i := 0; i < 7; i++ {
		wd := startDate.AddDate(0, 0, i).Weekday()
		if w.IsVisible(wd) {
			return wd
		}
	}
	return time.Weekday(-1)
}

// Last returns the last visible weekday based on the start date.
func (w Weekdays) Last(startDate time.Time) time.Weekday {
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	for i := 1; i < 7; i++ {
		wd := startDate.AddDate(0, 0, (-1 * i)).Weekday()
		if w.IsVisible(wd) {
			return wd
		}
	}
	return time.Weekday(-1)
}

// Next returns the next visible weekday based on the start date.
func (w Weekdays) Next(startDate time.Time) time.Weekday {
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	startDate = startDate.AddDate(0, 0, 1)
	return w.First((startDate))
}

// Previous returns the prior visible weekday based on the start date.
func (w Weekdays) Previous(startDate time.Time) time.Weekday {
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	return w.Last(startDate)
}

// DefaultWeekdays returns default weekday labels.
func DefaultWeekdays() Weekdays {
	return Weekdays{
		time.Sunday:    "Sun",
		time.Monday:    "Mon",
		time.Tuesday:   "Tue",
		time.Wednesday: "Wed",
		time.Thursday:  "Thu",
		time.Friday:    "Fri",
		time.Saturday:  "Sat",
	}
}

// DefaultWeekdaysShort returns default short weekday labels.
func DefaultWeekdaysShort() Weekdays {
	return Weekdays{
		time.Sunday:    "U",
		time.Monday:    "M",
		time.Tuesday:   "T",
		time.Wednesday: "W",
		time.Thursday:  "R",
		time.Friday:    "F",
		time.Saturday:  "S",
	}
}

// WeekModel represents a calendar week.
type WeekModel struct {
	// keyMap is key bindings for calendar navigation
	keyMap KeyMap

	// startOfWeek is the day that represents the beginning of the week
	startOfWeek time.Weekday

	// weekdays manages labels for weekdays
	weekdays Weekdays

	// week to represent
	startDate time.Time

	// days contains user-provided information about each day
	days map[time.Time]tea.Model

	activeDate time.Time

	// Styles
	styles WeekStyles
}

// NewWeek creates a new WeekModel.
func NewWeek(sampleDate time.Time) WeekModel {
	startDay := sampleDate.AddDate(0, 0, (-1 * int(sampleDate.Weekday())))

	m := WeekModel{
		keyMap: DefaultWeekKeyMap(),

		startOfWeek: time.Sunday,
		weekdays:    DefaultWeekdays(),

		startDate:  startDay,
		days:       make(map[time.Time]tea.Model),
		activeDate: time.Time{},

		styles: DefaultWeekStyles(),
	}

	return m
}

// StartOfWeek sets the first day of the week.
func (m WeekModel) StartOfWeek(weekday time.Weekday) WeekModel {
	// Slide startDay forward or backward to match new start of week
	daysDiff := int(weekday) - int(m.startOfWeek)
	m.startDate = m.startDate.AddDate(0, 0, daysDiff)

	m.startOfWeek = weekday

	return m
}

// Weekdays sets custom weekday labels.
func (m WeekModel) Weekdays(weekdays Weekdays) WeekModel {
	m.weekdays = weekdays
	return m
}

// Styles sets custom styling.
func (m WeekModel) Styles(styles WeekStyles) WeekModel {
	m.styles = styles
	return m
}

// PreviousDate sets the activeDate to the previous visible date.
//
// Notes:
//   - The Weekdays map is used to determine the previous date that has a day label and should therefore
//     be visible.
//   - If moving backwards from the first visible day, the method rolls over to the last visible day of the week.
//   - If the active date is unset (initial state), this method will set the last visible weekday as the "previous"
//     date.
func (m WeekModel) PreviousDate() WeekModel {
	ad := m.activeDate
	if m.activeDate == (time.Time{}) {
		ad = m.startDate
	}

	var prevWeekday time.Weekday
	for d := 1; d < 7; d++ {
		pw := time.Weekday((7 + int(ad.Weekday()) - d) % 7)
		if m.weekdays.IsVisible(pw) {
			prevWeekday = pw
			break
		}
	}

	daysDiff := int(prevWeekday) - int(ad.Weekday())
	m.activeDate = ad.AddDate(0, 0, daysDiff)

	return m
}

// NextDate sets the activeDate to the next visible date.
//
// Notes:
//   - The Weekdays map is used to determine the next date that has a day label and should therefore be visible.
//   - If moving forwards from the last visible day, the method rolls over to the first visible day of the week.
//   - If the active date is unset (initial state), this method will set the first visible weekday as the "next"
//     date.
func (m WeekModel) NextDate() WeekModel {
	ad := m.activeDate
	if m.activeDate == (time.Time{}) {
		ad = m.startDate
	}

	var nextWeekday time.Weekday
	for d := 0; d < 7; d++ {
		nw := time.Weekday((int(ad.Weekday()) + d) % 7)
		if m.weekdays.IsVisible(nw) {
			// If it is the "same" day and it is not directly after being initialized, skip it.
			// Placing the initialization logic inside the loop allows [WeekModel] to have a startDate that is before
			// the first visible day. For example, the week starts on Sunday but the first visible day is Monday.
			if (d == 0) && (m.activeDate != (time.Time{})) {
				continue
			}
			nextWeekday = nw
			break
		}
	}

	daysDiff := int(nextWeekday) - int(ad.Weekday())
	m.activeDate = ad.AddDate(0, 0, daysDiff)

	return m
}

// Init the WeekModel.
func (m WeekModel) Init() tea.Cmd { return nil }

// Update the WeekModel.
func (m WeekModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		oldActiveDate := m.activeDate
		switch {
		case key.Matches(msg, m.keyMap.Left):
			m = m.PreviousDate()
		case key.Matches(msg, m.keyMap.Right):
			m = m.NextDate()
		}

		if oldActiveDate != m.activeDate {
			cmds = append(cmds, func() tea.Msg {
				return ActiveDateMsg{
					Date: time.Date(m.activeDate.Year(), m.activeDate.Month(), m.activeDate.Day(), 0, 0, 0, 0, time.UTC),
				}
			})
		}
	case DayContentMsg:
		if m.startDate.Compare(msg.Date) > 0 {
			break
		}
		if m.startDate.AddDate(0, 0, 6).Compare(msg.Date) < 0 {
			break
		}

		i := time.Date(msg.Date.Year(), msg.Date.Month(), msg.Date.Day(), 0, 0, 0, 0, time.UTC)

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

// View renders the WeekModel.
func (m WeekModel) View() string {
	return gloss.JoinVertical(
		gloss.Top,
		m.ViewHeaders(),
		m.ViewDates(),
	)
}

// ViewHeaders renders the weekday headers.
func (m WeekModel) ViewHeaders() string {
	first := m.weekdays.First(m.startDate)
	last := m.weekdays.Last(m.startDate)

	var headers []string
	for i := 0; i < 7; i++ {
		day := m.startDate.AddDate(0, 0, i)
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

		label = fmt.Sprintf(
			"%s\n%s",
			label,
			day.Format(m.styles.DateFormat),
		)

		if day.Compare(m.activeDate) == 0 {
			headerStyle := m.styles.ActiveHeaderStyle

			label = headerStyle.Render(label)
		}

		headers = append(headers, style.Render(label))
	}

	return gloss.JoinHorizontal(gloss.Top, headers...)
}

// ViewDates renders the individual dates.
func (m WeekModel) ViewDates() string {
	style := m.styles.DateStyles.BodyStyle.
		Width(m.styles.DateStyles.Width).
		Height(m.styles.DateStyles.Height)

	var days []string
	var maxHeight int
	for i := 0; i < 7; i++ {
		day := m.startDate.AddDate(0, 0, i)
		if !m.weekdays.IsVisible(day.Weekday()) {
			// Don't render anything for that day if it isn't in the header list.
			continue
		}

		body := style.Render("")
		if content, ok := m.days[day]; ok {
			body = style.Render(content.View())
		}

		maxHeight = max(maxHeight, gloss.Height(body))

		days = append(days, body)
	}

	for i := 0; i < len(days); i++ {
		// Figure out if the border style is left, middle, or right
		style := m.styles.MiddleDayStyle
		switch i {
		case 0:
			style = m.styles.LeftDayStyle
		case len(days) - 1:
			style = m.styles.RightDayStyle
		}

		h := maxHeight + style.GetPaddingTop() + style.GetPaddingBottom()
		style = style.Height(h)

		days[i] = style.Render(days[i])
	}

	return gloss.JoinHorizontal(gloss.Top, days...)
}
