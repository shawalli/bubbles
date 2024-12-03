package calendar

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/golden"
	"github.com/stretchr/testify/assert"
)

func Test_NewMonth(t *testing.T) {
	// Test
	got := NewMonth(2024, time.September)

	// Assertions
	assert.Equal(t, 2024, got.year)
	assert.Equal(t, time.September, got.month)
	assert.Equal(t, time.Sunday, got.startOfWeek)
	assert.Zero(t, got.activeDay)
	assert.Empty(t, got.days)
}

func TestMonthModel_StartOfWeek(t *testing.T) {
	// Setup
	tm := NewMonth(2024, time.September)

	// Test
	got := tm.StartOfWeek(time.Tuesday)

	// Assertions
	assert.Equal(t, time.Tuesday, got.startOfWeek)
}

func TestMonthModel_Weekdays(t *testing.T) {
	// Setup
	tm := NewMonth(2024, time.September)

	weekdays := DefaultWeekdaysShort()

	// Test
	got := tm.Weekdays(weekdays)

	// Assertions
	assert.Equal(t, weekdays, got.weekdays)
}

func TestMonthModel_Styles(t *testing.T) {
	// Setup
	tm := NewMonth(2024, time.September)

	styles := MonthStyles{}

	// Test
	got := tm.Styles(styles)

	// Assertions
	assert.Equal(t, styles, got.styles)
}

func TestMonthModel_Update(t *testing.T) {
	tests := []struct {
		name          string
		model         MonthModel
		msgs          []tea.Msg
		wantActiveDay int
		wantMsgs      []tea.Msg
	}{
		{
			name:  "first-right",
			model: NewMonth(2024, time.September),
			msgs: []tea.Msg{
				tea.KeyMsg{Type: tea.KeyRight},
			},
			wantActiveDay: 1,
			wantMsgs: []tea.Msg{
				ActiveDateMsg{Date: time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC)},
			},
		},
		{
			name:  "first-down",
			model: NewMonth(2024, time.September),
			msgs: []tea.Msg{
				tea.KeyMsg{Type: tea.KeyRight},
			},
			wantActiveDay: 1,
			wantMsgs: []tea.Msg{
				ActiveDateMsg{Date: time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC)},
			},
		},
		{
			name:  "first-up",
			model: NewMonth(2024, time.September),
			msgs: []tea.Msg{
				tea.KeyMsg{Type: tea.KeyUp},
			},
			wantActiveDay: 29,
			wantMsgs: []tea.Msg{
				ActiveDateMsg{Date: time.Date(2024, time.September, 29, 0, 0, 0, 0, time.UTC)},
			},
		},
		{
			name:  "first-left",
			model: NewMonth(2024, time.September),
			msgs: []tea.Msg{
				tea.KeyMsg{Type: tea.KeyLeft},
			},
			wantActiveDay: 30,
			wantMsgs: []tea.Msg{
				ActiveDateMsg{Date: time.Date(2024, time.September, 30, 0, 0, 0, 0, time.UTC)},
			},
		},
		{
			name:  "right-up",
			model: NewMonth(2024, time.September),
			msgs: []tea.Msg{
				tea.KeyMsg{Type: tea.KeyRight},
				tea.KeyMsg{Type: tea.KeyUp},
			},
			wantActiveDay: 29,
			wantMsgs: []tea.Msg{
				ActiveDateMsg{Date: time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC)},
				ActiveDateMsg{Date: time.Date(2024, time.September, 29, 0, 0, 0, 0, time.UTC)},
			},
		},
		{
			name: "month-wraparound-right",
			model: func() MonthModel {
				m := NewMonth(2024, time.September)
				m.activeDay = 30
				return m
			}(),
			msgs: []tea.Msg{
				tea.KeyMsg{Type: tea.KeyRight},
			},
			wantActiveDay: 1,
			wantMsgs: []tea.Msg{
				ActiveDateMsg{Date: time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC)},
			},
		},
		{
			name: "month-wraparound-left",
			model: func() MonthModel {
				m := NewMonth(2024, time.September)
				m.activeDay = 1
				return m
			}(),
			msgs: []tea.Msg{
				tea.KeyMsg{Type: tea.KeyLeft},
			},
			wantActiveDay: 30,
			wantMsgs: []tea.Msg{
				ActiveDateMsg{Date: time.Date(2024, time.September, 30, 0, 0, 0, 0, time.UTC)},
			},
		},
		{
			name: "week-wraparound-right",
			model: func() MonthModel {
				m := NewMonth(2024, time.September)
				m.activeDay = 6
				return m
			}(),
			msgs: []tea.Msg{
				tea.KeyMsg{Type: tea.KeyRight},
			},
			wantActiveDay: 7,
			wantMsgs: []tea.Msg{
				ActiveDateMsg{Date: time.Date(2024, time.September, 7, 0, 0, 0, 0, time.UTC)},
			},
		},
		{
			name: "week-wraparound-left",
			model: func() MonthModel {
				m := NewMonth(2024, time.September)
				m.activeDay = 14
				return m
			}(),
			msgs: []tea.Msg{
				tea.KeyMsg{Type: tea.KeyLeft},
			},
			wantActiveDay: 13,
			wantMsgs: []tea.Msg{
				ActiveDateMsg{Date: time.Date(2024, time.September, 13, 0, 0, 0, 0, time.UTC)},
			},
		},
		{
			name: "week-skip-up",
			model: func() MonthModel {
				m := NewMonth(2024, time.September)
				m.activeDay = 11
				return m
			}(),
			msgs: []tea.Msg{
				tea.KeyMsg{Type: tea.KeyUp},
			},
			wantActiveDay: 4,
			wantMsgs: []tea.Msg{
				ActiveDateMsg{Date: time.Date(2024, time.September, 4, 0, 0, 0, 0, time.UTC)},
			},
		},
		{
			name: "week-skip-down",
			model: func() MonthModel {
				m := NewMonth(2024, time.September)
				m.activeDay = 11
				return m
			}(),
			msgs: []tea.Msg{
				tea.KeyMsg{Type: tea.KeyDown},
			},
			wantActiveDay: 18,
			wantMsgs: []tea.Msg{
				ActiveDateMsg{Date: time.Date(2024, time.September, 18, 0, 0, 0, 0, time.UTC)},
			},
		},
		{
			name: "week",
			model: func() MonthModel {
				m := NewMonth(2024, time.September)
				m.activeDay = 10
				return m
			}(),
			msgs: []tea.Msg{
				tea.KeyMsg{Type: tea.KeyRight},
				tea.KeyMsg{Type: tea.KeyRight},
				tea.KeyMsg{Type: tea.KeyRight},
				tea.KeyMsg{Type: tea.KeyRight},
				tea.KeyMsg{Type: tea.KeyRight},
				tea.KeyMsg{Type: tea.KeyRight},
				tea.KeyMsg{Type: tea.KeyRight},
			},
			wantActiveDay: 17,
			wantMsgs: []tea.Msg{
				ActiveDateMsg{Date: time.Date(2024, time.September, 11, 0, 0, 0, 0, time.UTC)},
				ActiveDateMsg{Date: time.Date(2024, time.September, 12, 0, 0, 0, 0, time.UTC)},
				ActiveDateMsg{Date: time.Date(2024, time.September, 13, 0, 0, 0, 0, time.UTC)},
				ActiveDateMsg{Date: time.Date(2024, time.September, 14, 0, 0, 0, 0, time.UTC)},
				ActiveDateMsg{Date: time.Date(2024, time.September, 15, 0, 0, 0, 0, time.UTC)},
				ActiveDateMsg{Date: time.Date(2024, time.September, 16, 0, 0, 0, 0, time.UTC)},
				ActiveDateMsg{Date: time.Date(2024, time.September, 17, 0, 0, 0, 0, time.UTC)},
			},
		},
		{
			name: "random",
			model: func() MonthModel {
				m := NewMonth(2024, time.September)
				m.activeDay = 4
				return m
			}(),
			msgs: []tea.Msg{
				tea.KeyMsg{Type: tea.KeyRight},
				tea.KeyMsg{Type: tea.KeyDown},
				tea.KeyMsg{Type: tea.KeyDown},
				tea.KeyMsg{Type: tea.KeyLeft},
				tea.KeyMsg{Type: tea.KeyLeft},
				tea.KeyMsg{Type: tea.KeyLeft},
				tea.KeyMsg{Type: tea.KeyLeft},
				tea.KeyMsg{Type: tea.KeyLeft},
				tea.KeyMsg{Type: tea.KeyLeft},
				tea.KeyMsg{Type: tea.KeyUp},
				tea.KeyMsg{Type: tea.KeyUp},
				tea.KeyMsg{Type: tea.KeyRight},
				tea.KeyMsg{Type: tea.KeyRight},
				tea.KeyMsg{Type: tea.KeyLeft},
				tea.KeyMsg{Type: tea.KeyRight},
			},
			wantActiveDay: 29,
			wantMsgs: []tea.Msg{
				ActiveDateMsg{Date: time.Date(2024, time.September, 5, 0, 0, 0, 0, time.UTC)},
				ActiveDateMsg{Date: time.Date(2024, time.September, 12, 0, 0, 0, 0, time.UTC)},
				ActiveDateMsg{Date: time.Date(2024, time.September, 19, 0, 0, 0, 0, time.UTC)},
				ActiveDateMsg{Date: time.Date(2024, time.September, 18, 0, 0, 0, 0, time.UTC)},
				ActiveDateMsg{Date: time.Date(2024, time.September, 17, 0, 0, 0, 0, time.UTC)},
				ActiveDateMsg{Date: time.Date(2024, time.September, 16, 0, 0, 0, 0, time.UTC)},
				ActiveDateMsg{Date: time.Date(2024, time.September, 15, 0, 0, 0, 0, time.UTC)},
				ActiveDateMsg{Date: time.Date(2024, time.September, 14, 0, 0, 0, 0, time.UTC)},
				ActiveDateMsg{Date: time.Date(2024, time.September, 13, 0, 0, 0, 0, time.UTC)},
				ActiveDateMsg{Date: time.Date(2024, time.September, 6, 0, 0, 0, 0, time.UTC)},
				ActiveDateMsg{Date: time.Date(2024, time.September, 27, 0, 0, 0, 0, time.UTC)},
				ActiveDateMsg{Date: time.Date(2024, time.September, 28, 0, 0, 0, 0, time.UTC)},
				ActiveDateMsg{Date: time.Date(2024, time.September, 29, 0, 0, 0, 0, time.UTC)},
				ActiveDateMsg{Date: time.Date(2024, time.September, 28, 0, 0, 0, 0, time.UTC)},
				ActiveDateMsg{Date: time.Date(2024, time.September, 29, 0, 0, 0, 0, time.UTC)},
			},
		},

		{
			name: "complex-first-right",
			model: func() MonthModel {
				m := NewMonth(2024, time.October).
					Weekdays(Weekdays{
						time.Thursday: "R",
						time.Friday:   "F",
						time.Saturday: "S",
						time.Sunday:   "U",
					}).
					StartOfWeek(time.Tuesday)
				return m
			}(),
			msgs: []tea.Msg{
				tea.KeyMsg{Type: tea.KeyRight},
			},
			wantActiveDay: 3,
			wantMsgs: []tea.Msg{
				ActiveDateMsg{Date: time.Date(2024, time.October, 3, 0, 0, 0, 0, time.UTC)},
			},
		},
		{
			name: "weekend-first-down",
			model: func() MonthModel {
				m := NewMonth(2024, time.October).
					Weekdays(Weekdays{
						time.Friday:   "F",
						time.Saturday: "S",
						time.Sunday:   "U",
					}).
					StartOfWeek(time.Saturday)
				return m
			}(),
			msgs: []tea.Msg{
				tea.KeyMsg{Type: tea.KeyDown},
			},
			wantActiveDay: 4,
			wantMsgs: []tea.Msg{
				ActiveDateMsg{Date: time.Date(2024, time.October, 4, 0, 0, 0, 0, time.UTC)},
			},
		},
		{
			name: "weekend-first-up",
			model: func() MonthModel {
				m := NewMonth(2024, time.October).
					Weekdays(Weekdays{
						time.Friday:   "F",
						time.Saturday: "S",
						time.Sunday:   "U",
					}).
					StartOfWeek(time.Saturday)
				return m
			}(),
			msgs: []tea.Msg{
				tea.KeyMsg{Type: tea.KeyUp},
			},
			wantActiveDay: 25,
			wantMsgs: []tea.Msg{
				ActiveDateMsg{Date: time.Date(2024, time.October, 25, 0, 0, 0, 0, time.UTC)},
			},
		},
		{
			name: "weekend-first-left",
			model: func() MonthModel {
				m := NewMonth(2024, time.October).
					Weekdays(Weekdays{
						time.Friday:   "F",
						time.Saturday: "S",
						time.Sunday:   "U",
					}).
					StartOfWeek(time.Saturday)
				return m
			}(),
			msgs: []tea.Msg{
				tea.KeyMsg{Type: tea.KeyLeft},
			},
			wantActiveDay: 27,
			wantMsgs: []tea.Msg{
				ActiveDateMsg{Date: time.Date(2024, time.October, 27, 0, 0, 0, 0, time.UTC)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			tm := tt.model
			_ = tm.Init()
			fmt.Println(ansi.Strip(tm.View()))

			// Test
			var gotMsgs []tea.Msg
			for _, msg := range tt.msgs {
				gotModel, gotCmd := tm.Update(msg)
				tm = gotModel.(MonthModel)
				gotMsgs = append(gotMsgs, gotCmd())
			}

			// Assertions
			assert.Equal(t, tt.wantActiveDay, tm.activeDay)
			assert.Equal(t, tt.wantMsgs, gotMsgs)
		})
	}
}

func TestMonthModel_View(t *testing.T) {
	tests := []struct {
		name        string
		startOfWeek time.Weekday
		weekdays    Weekdays
	}{
		{
			name:        "seven-day-week",
			startOfWeek: time.Sunday,
			weekdays:    DefaultWeekdays(),
		},
		{
			name:        "seven-day-week-monday",
			startOfWeek: time.Monday,
			weekdays:    DefaultWeekdays(),
		},
		{
			name:        "seven-day-week-tuesday",
			startOfWeek: time.Tuesday,
			weekdays:    DefaultWeekdays(),
		},
		{
			name:        "seven-day-week-wednesday",
			startOfWeek: time.Wednesday,
			weekdays:    DefaultWeekdays(),
		},
		{
			name:        "seven-day-week-thursday",
			startOfWeek: time.Thursday,
			weekdays:    DefaultWeekdays(),
		},
		{
			name:        "seven-day-week-friday",
			startOfWeek: time.Friday,
			weekdays:    DefaultWeekdays(),
		},
		{
			name:        "seven-day-week-saturday",
			startOfWeek: time.Saturday,
			weekdays:    DefaultWeekdays(),
		},
		{
			name:        "five-day-week",
			startOfWeek: time.Sunday,
			weekdays: Weekdays{
				time.Monday:    "Mon",
				time.Tuesday:   "Tue",
				time.Wednesday: "Wed",
				time.Thursday:  "Thu",
				time.Friday:    "Fri",
			},
		},
		{
			name:        "five-day-week-monday",
			startOfWeek: time.Monday,
			weekdays: Weekdays{
				time.Monday:    "Mon",
				time.Tuesday:   "Tue",
				time.Wednesday: "Wed",
				time.Thursday:  "Thu",
				time.Friday:    "Fri",
			},
		},
		{
			name:        "five-day-week-tuesday",
			startOfWeek: time.Tuesday,
			weekdays: Weekdays{
				time.Monday:    "Mon",
				time.Tuesday:   "Tue",
				time.Wednesday: "Wed",
				time.Thursday:  "Thu",
				time.Friday:    "Fri",
			},
		},
		{
			name:        "five-day-week-wednesday",
			startOfWeek: time.Wednesday,
			weekdays: Weekdays{
				time.Monday:    "Mon",
				time.Tuesday:   "Tue",
				time.Wednesday: "Wed",
				time.Thursday:  "Thu",
				time.Friday:    "Fri",
			},
		},
		{
			name:        "five-day-week-thursday",
			startOfWeek: time.Thursday,
			weekdays: Weekdays{
				time.Monday:    "Mon",
				time.Tuesday:   "Tue",
				time.Wednesday: "Wed",
				time.Thursday:  "Thu",
				time.Friday:    "Fri",
			},
		},
		{
			name:        "five-day-week-friday",
			startOfWeek: time.Friday,
			weekdays: Weekdays{
				time.Monday:    "Mon",
				time.Tuesday:   "Tue",
				time.Wednesday: "Wed",
				time.Thursday:  "Thu",
				time.Friday:    "Fri",
			},
		},
		{
			name:        "weekend",
			startOfWeek: time.Sunday,
			weekdays: Weekdays{
				time.Friday:   "Fri",
				time.Saturday: "Sat",
				time.Sunday:   "Sun",
			},
		},
		{
			name:        "weekend-friday",
			startOfWeek: time.Friday,
			weekdays: Weekdays{
				time.Friday:   "Fri",
				time.Saturday: "Sat",
				time.Sunday:   "Sun",
			},
		},
		{
			name:        "weekend-saturday",
			startOfWeek: time.Saturday,
			weekdays: Weekdays{
				time.Friday:   "Fri",
				time.Saturday: "Sat",
				time.Sunday:   "Sun",
			},
		},
		{
			name:        "complex",
			startOfWeek: time.Tuesday,
			weekdays: Weekdays{
				time.Monday:    "Mon",
				time.Wednesday: "Wed",
				time.Thursday:  "Thu",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			tm := NewMonth(2024, time.September).
				Weekdays(tt.weekdays).
				StartOfWeek(tt.startOfWeek)
			tm.days[4] = testDayModel{}
			_ = tm.Init()
			tm.activeDay = 10

			// Test
			c := tm.View()
			got := ansi.Strip(c)

			// Assertions
			golden.RequireEqual(t, []byte(got))
		})
	}
}

func TestMonthModel_View_Comprehensive(t *testing.T) {
	tests := []int{
		2022,
		2023,
		2024,
		2025,
	}

	for _, year := range tests {
		t.Run(strconv.Itoa(year), func(t *testing.T) {
			// Setup
			var gots []string
			for i := time.January; i <= time.December; i += 1 {

				tm := NewMonth(year, i)
				_ = tm.Init()

				// Test
				c := tm.View()
				gots = append(gots, ansi.Strip(c))
			}

			// Assertions
			got := strings.Join(gots, "\n")
			golden.RequireEqual(t, []byte(got))
		})
	}
}

func TestMonthModel_Title(t *testing.T) {
	year := 2024
	month := time.September

	tests := []struct {
		name        string
		includeYear bool
		want        string
	}{
		{
			name:        "month",
			includeYear: false,
			want:        "September",
		},
		{
			name:        "month-and-year",
			includeYear: true,
			want:        "September 2024",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			tm := MonthModel{
				year:  year,
				month: month,
			}

			// Test
			got := tm.Title(tt.includeYear)

			// Assertions
			assert.Equal(t, tt.want, got)
		})
	}
}
