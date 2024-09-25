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
func (l Weekdays) Get(day time.Weekday) (string, bool) {
	s, ok := l[day]

	return s, ok
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

// NewMonth creates a new WeekModel.
func NewWeek(sampleDate time.Time) WeekModel {
	startDay := sampleDate.AddDate(0, 0, -(int(sampleDate.Weekday()) - int(time.Sunday)))

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
//   - The [MonthModelweekdays] map is used to determine the previous date that has a day label and should therefore
//     be visible.
//   - If moving backwards from the first visible day, the method rolls over to the last visible day of the week.
func (m WeekModel) PreviousDate() WeekModel {
	ad := m.activeDate
	if m.activeDate == (time.Time{}) {
		ad = m.startDate
	}

	var prevWeekday time.Weekday
	for d := 1; d < 7; d++ {
		pw := time.Weekday((7 + int(ad.Weekday()) - d) % 7)
		if _, ok := m.weekdays.Get(pw); ok {
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
//   - The [MonthModelweekdays] map is used to determine the next date that has a day label and should therefore be
//     visible.
//   - If moving forwards from the last visible day, the method rolls over to the first visible day of the week.
//   - If the active date is unset (initial state), this method will set the next
func (m WeekModel) NextDate() WeekModel {
	ad := m.activeDate
	if m.activeDate == (time.Time{}) {
		ad = m.startDate
	}

	var nextWeekday time.Weekday
	for d := 0; d < 7; d++ {
		nw := time.Weekday((int(ad.Weekday()) + d) % 7)
		if _, ok := m.weekdays.Get(nw); ok {
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

// ViewHeaders renders the weekdays headers.
func (m WeekModel) ViewHeaders() string {
	var headers []string

	// Simplify successive loop by pre-caclulating first and last visible days
	firstVisibleDay := time.Weekday(-1)
	lastVisibleDay := time.Weekday(-1)
	for i := 0; i < 7; i++ {
		day := m.startDate.AddDate(0, 0, i)
		if _, ok := m.weekdays.Get(day.Weekday()); ok {
			if firstVisibleDay == time.Weekday(-1) {
				firstVisibleDay = day.Weekday()
			}
			lastVisibleDay = day.Weekday()
		}
	}

	for i := 0; i < 7; i++ {
		day := m.startDate.AddDate(0, 0, i)
		label, ok := m.weekdays.Get(day.Weekday())
		if !ok {
			continue
		}

		style := m.styles.MiddleHeaderStyle
		switch day.Weekday() {
		case firstVisibleDay:
			style = m.styles.LeftHeaderStyle
		case lastVisibleDay:
			style = m.styles.RightHeaderStyle
		}

		style = style.Width(m.styles.CellStyles.Width)

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
	style := m.styles.CellStyles.BodyStyle.
		Width(m.styles.CellStyles.Width).
		Height(m.styles.CellStyles.Height)

	var days []string
	var maxHeight int
	for i := 0; i < 7; i++ {
		day := m.startDate.AddDate(0, 0, i)
		if _, ok := m.weekdays.Get(day.Weekday()); !ok {
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
		style := m.styles.MiddleCellStyle
		switch i {
		case 0:
			style = m.styles.LeftCellStyle
		case len(days) - 1:
			style = m.styles.RightCellStyle
		}

		h := maxHeight + style.GetPaddingTop() + style.GetPaddingBottom()
		style = style.Height(h)

		days[i] = style.Render(days[i])
	}

	return gloss.JoinHorizontal(gloss.Top, days...)
}
