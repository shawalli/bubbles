package calendar

import (
	"testing"
	"time"

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

		//dim:30 fwdom:Sunday wim:5
		//i:1 sow:Monday fvwkd:Monday wd:Sunday l(wks):0 l(wk):0

		//dim:30 fwdom:Sunday wim:6
		//i:1 sow:Monday fvwkd:Monday wd:Monday l(wks):0 l(wk):0

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
		/*
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
						}, {
							name:        "weekend-saturday",
							startOfWeek: time.Saturday,
							weekdays: Weekdays{
								time.Friday:   "Fri",
								time.Saturday: "Sat",
								time.Sunday:   "Sun",
							},
						},
		*/
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			tm := NewMonth(2024, time.September).
				Weekdays(tt.weekdays).
				StartOfWeek(tt.startOfWeek)
			tm.days[4] = testDayModel{}
			_ = tm.Init()

			// Test
			c := tm.View()
			got := ansi.Strip(c)

			// Assertions
			golden.RequireEqual(t, []byte(got))
		})
	}
}
