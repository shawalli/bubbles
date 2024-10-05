package calendar

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/golden"
	"github.com/stretchr/testify/assert"
)

type testDayModel struct{}

func (m testDayModel) Init() tea.Cmd { return nil }
func (m testDayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, func() tea.Msg { return nil }
}
func (m testDayModel) View() string { return "" }

func newDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func TestWeekdays_Get(t *testing.T) {
	// Setup
	weekdays := Weekdays{
		time.Sunday: "Sun",
	}

	// Test
	got, gotOK := weekdays.Get(time.Sunday)

	// Assertions
	assert.True(t, gotOK)
	assert.Equal(t, "Sun", got)

	// Test
	got, gotOK = weekdays.Get(time.Monday)

	// Assertions
	assert.False(t, gotOK)
	assert.Empty(t, got)
}

func TestWeekdays_IsVisible(t *testing.T) {
	// Setup
	weekdays := Weekdays{
		time.Sunday: "Sun",
	}

	// Test
	got := weekdays.IsVisible(time.Sunday)

	// Assertions
	assert.True(t, got)

	// Test
	got = weekdays.IsVisible(time.Monday)

	// Assertions
	assert.False(t, got)
}

func TestWeekdays_First_NotFound(t *testing.T) {
	// Setup
	weekdays := Weekdays{}

	startDate := newDate(2024, time.September, 22)

	// Test
	got := weekdays.First(startDate)

	// Assertions
	assert.Equal(t, time.Weekday(-1), got)
}

func TestWeekdays_First(t *testing.T) {
	sundayStartDate := newDate(2024, time.September, 22)
	mondayStartDate := newDate(2024, time.September, 23)
	thursdayStartDate := newDate(2024, time.September, 26)
	saturdayStartDate := newDate(2024, time.September, 28)

	tests := []struct {
		name      string
		weekdays  Weekdays
		startDate time.Time
		want      time.Weekday
	}{
		{
			name:      "sd-sunday-wd-sunday",
			weekdays:  Weekdays{time.Sunday: ""},
			startDate: sundayStartDate,
			want:      time.Sunday,
		},
		{
			name:      "sd-sunday-wd-monday",
			weekdays:  Weekdays{time.Monday: ""},
			startDate: sundayStartDate,
			want:      time.Monday,
		},
		{
			name:      "sd-sunday-wd-friday",
			weekdays:  Weekdays{time.Friday: ""},
			startDate: sundayStartDate,
			want:      time.Friday,
		},
		{
			name:      "sd-sunday-wd-saturday",
			weekdays:  Weekdays{time.Saturday: ""},
			startDate: sundayStartDate,
			want:      time.Saturday,
		},
		{
			name:      "sd-monday-wd-monday",
			weekdays:  Weekdays{time.Sunday: "", time.Monday: ""},
			startDate: mondayStartDate,
			want:      time.Monday,
		},
		{
			name:      "sd-monday-wd-tuesday",
			weekdays:  Weekdays{time.Sunday: "", time.Tuesday: ""},
			startDate: mondayStartDate,
			want:      time.Tuesday,
		},
		{
			name:      "sd-monday-wd-saturday",
			weekdays:  Weekdays{time.Sunday: "", time.Saturday: ""},
			startDate: mondayStartDate,
			want:      time.Saturday,
		},
		{
			name:      "sd-monday-wd-sunday",
			weekdays:  Weekdays{time.Sunday: ""},
			startDate: mondayStartDate,
			want:      time.Sunday,
		},
		{
			name:      "sd-thursday-wd-thursday",
			weekdays:  Weekdays{time.Sunday: "", time.Wednesday: "", time.Thursday: ""},
			startDate: thursdayStartDate,
			want:      time.Thursday,
		},
		{
			name:      "sd-thursday-wd-saturday",
			weekdays:  Weekdays{time.Sunday: "", time.Wednesday: "", time.Saturday: ""},
			startDate: thursdayStartDate,
			want:      time.Saturday,
		},
		{
			name:      "sd-thursday-wd-sunday",
			weekdays:  Weekdays{time.Sunday: "", time.Wednesday: ""},
			startDate: thursdayStartDate,
			want:      time.Sunday,
		},
		{
			name:      "sd-thursday-wd-wednesday",
			weekdays:  Weekdays{time.Wednesday: ""},
			startDate: thursdayStartDate,
			want:      time.Wednesday,
		},
		{
			name:      "sd-saturday-wd-saturday",
			weekdays:  Weekdays{time.Sunday: "", time.Saturday: ""},
			startDate: saturdayStartDate,
			want:      time.Saturday,
		},
		{
			name:      "sd-saturday-wd-sunday",
			weekdays:  Weekdays{time.Sunday: ""},
			startDate: saturdayStartDate,
			want:      time.Sunday,
		},
		{
			name:      "sd-saturday-wd-tuesday",
			weekdays:  Weekdays{time.Tuesday: ""},
			startDate: saturdayStartDate,
			want:      time.Tuesday,
		},
		{
			name:      "sd-saturday-wd-friday",
			weekdays:  Weekdays{time.Friday: ""},
			startDate: saturdayStartDate,
			want:      time.Friday,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test
			got := tt.weekdays.First(tt.startDate)

			// Assertions
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWeekdays_Last_NotFound(t *testing.T) {
	// Setup
	weekdays := Weekdays{}

	startDate := newDate(2024, time.September, 22)

	// Test
	got := weekdays.Last(startDate)

	// Assertions
	assert.Equal(t, time.Weekday(-1), got)
}

func TestWeekdays_Last(t *testing.T) {
	sundayStartDate := newDate(2024, time.September, 22)
	mondayStartDate := newDate(2024, time.September, 23)
	thursdayStartDate := newDate(2024, time.September, 26)
	saturdayStartDate := newDate(2024, time.September, 28)

	tests := []struct {
		name      string
		weekdays  Weekdays
		startDate time.Time
		want      time.Weekday
	}{
		{
			name:      "sd-sunday-wd-saturday",
			weekdays:  Weekdays{time.Sunday: "", time.Wednesday: "", time.Saturday: ""},
			startDate: sundayStartDate,
			want:      time.Saturday,
		},
		{
			name:      "sd-sunday-wd-friday",
			weekdays:  Weekdays{time.Sunday: "", time.Wednesday: "", time.Friday: ""},
			startDate: sundayStartDate,
			want:      time.Friday,
		},
		{
			name:      "sd-sunday-wd-monday",
			weekdays:  Weekdays{time.Sunday: "", time.Monday: ""},
			startDate: sundayStartDate,
			want:      time.Monday,
		},
		{
			name:      "sd-monday-wd-sunday",
			weekdays:  Weekdays{time.Monday: "", time.Wednesday: "", time.Sunday: ""},
			startDate: mondayStartDate,
			want:      time.Sunday,
		},
		{
			name:      "sd-monday-wd-saturday",
			weekdays:  Weekdays{time.Monday: "", time.Wednesday: "", time.Saturday: ""},
			startDate: mondayStartDate,
			want:      time.Saturday,
		},
		{
			name:      "sd-monday-wd-tuesday",
			weekdays:  Weekdays{time.Monday: "", time.Tuesday: ""},
			startDate: mondayStartDate,
			want:      time.Tuesday,
		},
		{
			name:      "sd-thursday-wd-wednesday",
			weekdays:  Weekdays{time.Thursday: "", time.Friday: "", time.Wednesday: ""},
			startDate: thursdayStartDate,
			want:      time.Wednesday,
		},
		{
			name:      "sd-thursday-wd-saturday",
			weekdays:  Weekdays{time.Thursday: "", time.Friday: "", time.Saturday: ""},
			startDate: thursdayStartDate,
			want:      time.Saturday,
		},
		{
			name:      "sd-thursday-wd-sunday",
			weekdays:  Weekdays{time.Thursday: "", time.Sunday: ""},
			startDate: thursdayStartDate,
			want:      time.Sunday,
		},
		{
			name:      "sd-saturday-wd-friday",
			weekdays:  Weekdays{time.Saturday: "", time.Sunday: "", time.Friday: ""},
			startDate: saturdayStartDate,
			want:      time.Friday,
		},
		{
			name:      "sd-saturday-wd-monday",
			weekdays:  Weekdays{time.Saturday: "", time.Sunday: "", time.Monday: ""},
			startDate: saturdayStartDate,
			want:      time.Monday,
		},
		{
			name:      "sd-saturday-wd-tuesday",
			weekdays:  Weekdays{time.Saturday: "", time.Tuesday: ""},
			startDate: saturdayStartDate,
			want:      time.Tuesday,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test
			got := tt.weekdays.Last(tt.startDate)

			// Assertions
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWeekdays_Next_NotFound(t *testing.T) {
	// Setup
	weekdays := Weekdays{}

	startDate := newDate(2024, time.September, 22)

	// Test
	got := weekdays.Next(startDate)

	// Assertions
	assert.Equal(t, time.Weekday(-1), got)
}

func TestWeekdays_Next(t *testing.T) {
	sundayStartDate := newDate(2024, time.September, 22)
	mondayStartDate := newDate(2024, time.September, 23)
	thursdayStartDate := newDate(2024, time.September, 26)
	saturdayStartDate := newDate(2024, time.September, 28)

	tests := []struct {
		name      string
		weekdays  Weekdays
		startDate time.Time
		want      time.Weekday
	}{
		{
			name:      "sd-sunday-wd-monday",
			weekdays:  Weekdays{time.Sunday: "", time.Monday: "", time.Tuesday: "", time.Wednesday: ""},
			startDate: sundayStartDate,
			want:      time.Monday,
		},
		{
			name:      "sd-sunday-wd-tuesday",
			weekdays:  Weekdays{time.Sunday: "", time.Tuesday: "", time.Wednesday: "", time.Thursday: ""},
			startDate: sundayStartDate,
			want:      time.Tuesday,
		},
		{
			name:      "sd-sunday-wd-friday",
			weekdays:  Weekdays{time.Sunday: "", time.Friday: "", time.Saturday: ""},
			startDate: sundayStartDate,
			want:      time.Friday,
		},
		{
			name:      "sd-sunday-wd-saturday",
			weekdays:  Weekdays{time.Sunday: "", time.Saturday: ""},
			startDate: sundayStartDate,
			want:      time.Saturday,
		},
		{
			name:      "sd-monday-wd-tuesday",
			weekdays:  Weekdays{time.Monday: "", time.Tuesday: "", time.Wednesday: ""},
			startDate: mondayStartDate,
			want:      time.Tuesday,
		},
		{
			name:      "sd-monday-wd-wednesday",
			weekdays:  Weekdays{time.Monday: "", time.Wednesday: "", time.Thursday: ""},
			startDate: mondayStartDate,
			want:      time.Wednesday,
		},
		{
			name:      "sd-monday-wd-saturday",
			weekdays:  Weekdays{time.Monday: "", time.Saturday: "", time.Sunday: ""},
			startDate: mondayStartDate,
			want:      time.Saturday,
		},
		{
			name:      "sd-monday-wd-sunday",
			weekdays:  Weekdays{time.Monday: "", time.Sunday: ""},
			startDate: mondayStartDate,
			want:      time.Sunday,
		},
		{
			name:      "sd-thursday-wd-friday",
			weekdays:  Weekdays{time.Thursday: "", time.Friday: "", time.Saturday: ""},
			startDate: thursdayStartDate,
			want:      time.Friday,
		},
		{
			name:      "sd-thursday-wd-saturday",
			weekdays:  Weekdays{time.Thursday: "", time.Saturday: "", time.Sunday: ""},
			startDate: thursdayStartDate,
			want:      time.Saturday,
		},
		{
			name:      "sd-thursday-wd-tuesday",
			weekdays:  Weekdays{time.Thursday: "", time.Tuesday: "", time.Wednesday: ""},
			startDate: thursdayStartDate,
			want:      time.Tuesday,
		},
		{
			name:      "sd-thursday-wd-wednesday",
			weekdays:  Weekdays{time.Thursday: "", time.Wednesday: ""},
			startDate: thursdayStartDate,
			want:      time.Wednesday,
		},
		{
			name:      "sd-saturday-wd-sunday",
			weekdays:  Weekdays{time.Saturday: "", time.Sunday: "", time.Monday: "", time.Tuesday: ""},
			startDate: saturdayStartDate,
			want:      time.Sunday,
		},
		{
			name:      "sd-saturday-wd-monday",
			weekdays:  Weekdays{time.Saturday: "", time.Monday: "", time.Tuesday: "", time.Wednesday: ""},
			startDate: saturdayStartDate,
			want:      time.Monday,
		},
		{
			name:      "sd-saturday-wd-thursday",
			weekdays:  Weekdays{time.Saturday: "", time.Thursday: "", time.Friday: ""},
			startDate: saturdayStartDate,
			want:      time.Thursday,
		},
		{
			name:      "sd-saturday-wd-friday",
			weekdays:  Weekdays{time.Saturday: "", time.Friday: ""},
			startDate: saturdayStartDate,
			want:      time.Friday,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test
			got := tt.weekdays.Next(tt.startDate)

			// Assertions
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWeekdays_Previous_NotFound(t *testing.T) {
	// Setup
	weekdays := Weekdays{}

	startDate := newDate(2024, time.September, 22)

	// Test
	got := weekdays.Previous(startDate)

	// Assertions
	assert.Equal(t, time.Weekday(-1), got)
}

func TestWeekdays_Previous(t *testing.T) {
	sundayStartDate := newDate(2024, time.September, 22)
	mondayStartDate := newDate(2024, time.September, 23)
	thursdayStartDate := newDate(2024, time.September, 26)
	saturdayStartDate := newDate(2024, time.September, 28)

	tests := []struct {
		name      string
		weekdays  Weekdays
		startDate time.Time
		want      time.Weekday
	}{
		{
			name:      "sd-sunday-wd-saturday",
			weekdays:  Weekdays{time.Sunday: "", time.Wednesday: "", time.Saturday: ""},
			startDate: sundayStartDate,
			want:      time.Saturday,
		},
		{
			name:      "sd-sunday-wd-friday",
			weekdays:  Weekdays{time.Sunday: "", time.Wednesday: "", time.Friday: ""},
			startDate: sundayStartDate,
			want:      time.Friday,
		},
		{
			name:      "sd-sunday-wd-monday",
			weekdays:  Weekdays{time.Sunday: "", time.Monday: ""},
			startDate: sundayStartDate,
			want:      time.Monday,
		},
		{
			name:      "sd-monday-wd-sunday",
			weekdays:  Weekdays{time.Monday: "", time.Wednesday: "", time.Sunday: ""},
			startDate: mondayStartDate,
			want:      time.Sunday,
		},
		{
			name:      "sd-monday-wd-saturday",
			weekdays:  Weekdays{time.Monday: "", time.Wednesday: "", time.Saturday: ""},
			startDate: mondayStartDate,
			want:      time.Saturday,
		},
		{
			name:      "sd-monday-wd-tuesday",
			weekdays:  Weekdays{time.Monday: "", time.Tuesday: ""},
			startDate: mondayStartDate,
			want:      time.Tuesday,
		},
		{
			name:      "sd-thursday-wd-wednesday",
			weekdays:  Weekdays{time.Thursday: "", time.Friday: "", time.Wednesday: ""},
			startDate: thursdayStartDate,
			want:      time.Wednesday,
		},
		{
			name:      "sd-thursday-wd-saturday",
			weekdays:  Weekdays{time.Thursday: "", time.Friday: "", time.Saturday: ""},
			startDate: thursdayStartDate,
			want:      time.Saturday,
		},
		{
			name:      "sd-thursday-wd-sunday",
			weekdays:  Weekdays{time.Thursday: "", time.Sunday: ""},
			startDate: thursdayStartDate,
			want:      time.Sunday,
		},
		{
			name:      "sd-saturday-wd-friday",
			weekdays:  Weekdays{time.Saturday: "", time.Sunday: "", time.Friday: ""},
			startDate: saturdayStartDate,
			want:      time.Friday,
		},
		{
			name:      "sd-saturday-wd-monday",
			weekdays:  Weekdays{time.Saturday: "", time.Sunday: "", time.Monday: ""},
			startDate: saturdayStartDate,
			want:      time.Monday,
		},
		{
			name:      "sd-saturday-wd-tuesday",
			weekdays:  Weekdays{time.Saturday: "", time.Tuesday: ""},
			startDate: saturdayStartDate,
			want:      time.Tuesday,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test
			got := tt.weekdays.Previous(tt.startDate)

			// Assertions
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_NewWeek(t *testing.T) {
	sampleDate := newDate(2024, time.September, 22)

	tests := []struct {
		name       string
		sampleDate time.Time
	}{
		{
			name:       "sunday",
			sampleDate: sampleDate,
		},
		{
			name:       "monday",
			sampleDate: sampleDate.AddDate(0, 0, int(time.Monday)),
		},
		{
			name:       "tuesday",
			sampleDate: sampleDate.AddDate(0, 0, int(time.Tuesday)),
		},
		{
			name:       "wednesday",
			sampleDate: sampleDate.AddDate(0, 0, int(time.Wednesday)),
		},
		{
			name:       "thursday",
			sampleDate: sampleDate.AddDate(0, 0, int(time.Thursday)),
		},
		{
			name:       "friday",
			sampleDate: sampleDate.AddDate(0, 0, int(time.Friday)),
		},
		{
			name:       "saturday",
			sampleDate: sampleDate.AddDate(0, 0, int(time.Saturday)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test
			got := NewWeek(tt.sampleDate)

			// Assertions
			assert.Equal(t, time.Sunday, got.startOfWeek)
			wantStartDate := newDate(2024, time.September, 22)
			assert.Equal(t, wantStartDate, got.startDate)
			assert.Empty(t, got.days)
			assert.Equal(t, time.Time{}, got.activeDate)
		})
	}
}

func TestWeekModel_StartOfWeek_Forward(t *testing.T) {
	// Setup
	sampleDate := newDate(2024, time.September, 26)

	tm := NewWeek(sampleDate)

	// Test
	got := tm.StartOfWeek(time.Tuesday)

	// Assertions
	wantStartDate := newDate(2024, time.September, 24)
	assert.Equal(t, wantStartDate, got.startDate)
	assert.Equal(t, time.Tuesday, got.startOfWeek)
}

func TestWeekModel_StartOfWeek_Backward(t *testing.T) {
	// Setup
	sampleDate := newDate(2024, time.September, 26)

	tm := NewWeek(sampleDate)
	tm.startOfWeek = time.Thursday
	tm.startDate = sampleDate

	// Test
	got := tm.StartOfWeek(time.Tuesday)

	// Assertions
	wantStartDate := newDate(2024, time.September, 24)
	assert.Equal(t, wantStartDate, got.startDate)
	assert.Equal(t, time.Tuesday, got.startOfWeek)
}

func TestWeekModel_StartOfWeek_Nop(t *testing.T) {
	// Setup
	sampleDate := newDate(2024, time.September, 26)

	tm := NewWeek(sampleDate)
	tm.startOfWeek = time.Thursday
	tm.startDate = sampleDate

	// Test
	got := tm.StartOfWeek(time.Thursday)

	// Assertions
	wantStartDate := newDate(2024, time.September, 26)
	assert.Equal(t, wantStartDate, got.startDate)
	assert.Equal(t, time.Thursday, got.startOfWeek)
}

func TestWeekModel_Weekdays(t *testing.T) {
	// Setup
	sampleDate := newDate(2024, time.September, 24)
	tm := NewWeek(sampleDate)

	weekdays := DefaultWeekdaysShort()

	// Test
	got := tm.Weekdays(weekdays)

	// Assertions
	assert.Equal(t, weekdays, got.weekdays)
}

func TestWeekModel_Styles(t *testing.T) {
	// Setup
	sampleDate := newDate(2024, time.September, 24)
	tm := NewWeek(sampleDate)

	styles := WeekStyles{}

	// Test
	got := tm.Styles(styles)

	// Assertions
	assert.Equal(t, styles, got.styles)
}

func TestWeekModel_PreviousDate(t *testing.T) {
	tests := []struct {
		name           string
		startOfWeek    time.Weekday
		weekdays       Weekdays
		activeDate     time.Time
		wantActiveDate time.Time
	}{
		{
			name:           "inital-sunday-prev-saturday",
			startOfWeek:    time.Sunday,
			weekdays:       Weekdays{time.Monday: "", time.Tuesday: "", time.Saturday: ""},
			activeDate:     time.Time{},
			wantActiveDate: newDate(2024, time.September, 28),
		},
		{
			name:           "inital-sunday-prev-wednesday",
			startOfWeek:    time.Sunday,
			weekdays:       Weekdays{time.Monday: "", time.Tuesday: "", time.Wednesday: ""},
			activeDate:     time.Time{},
			wantActiveDate: newDate(2024, time.September, 25),
		},
		{
			name:           "inital-tuesday-prev-monday",
			startOfWeek:    time.Tuesday,
			weekdays:       Weekdays{time.Monday: "", time.Tuesday: "", time.Wednesday: ""},
			activeDate:     time.Time{},
			wantActiveDate: newDate(2024, time.September, 23),
		},
		{
			name:           "inital-tuesday-prev-sunday",
			startOfWeek:    time.Tuesday,
			weekdays:       Weekdays{time.Sunday: "", time.Tuesday: "", time.Wednesday: ""},
			activeDate:     time.Time{},
			wantActiveDate: newDate(2024, time.September, 22),
		},
		{
			name:           "inital-tuesday-prev-wednesday",
			startOfWeek:    time.Tuesday,
			weekdays:       Weekdays{time.Tuesday: "", time.Wednesday: ""},
			activeDate:     time.Time{},
			wantActiveDate: newDate(2024, time.September, 25),
		},
		{
			name:           "tuesday-prev-wednesday",
			startOfWeek:    time.Sunday,
			weekdays:       Weekdays{time.Tuesday: "", time.Wednesday: ""},
			activeDate:     newDate(2024, time.September, 24),
			wantActiveDate: newDate(2024, time.September, 25),
		},
		{
			name:           "sunday-prev-friday",
			startOfWeek:    time.Sunday,
			weekdays:       Weekdays{time.Tuesday: "", time.Wednesday: "", time.Friday: ""},
			activeDate:     newDate(2024, time.September, 22),
			wantActiveDate: newDate(2024, time.September, 27),
		},
		{
			name:           "friday-prev-sunday",
			startOfWeek:    time.Sunday,
			weekdays:       Weekdays{time.Sunday: "", time.Friday: "", time.Saturday: ""},
			activeDate:     newDate(2024, time.September, 27),
			wantActiveDate: newDate(2024, time.September, 22),
		},
		{
			name:           "saturday-prev-thursday",
			startOfWeek:    time.Sunday,
			weekdays:       Weekdays{time.Sunday: "", time.Thursday: "", time.Saturday: ""},
			activeDate:     newDate(2024, time.September, 28),
			wantActiveDate: newDate(2024, time.September, 26),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			sampleDate := newDate(2024, time.September, 24)
			tm := NewWeek(sampleDate).
				Weekdays(tt.weekdays).
				StartOfWeek(tt.startOfWeek)
			tm.activeDate = tt.activeDate

			// Test
			got := tm.PreviousDate()

			// Assertions
			assert.Equal(t, tt.wantActiveDate, got.activeDate)
		})
	}
}

func TestWeekModel_NextDate(t *testing.T) {
	tests := []struct {
		name           string
		startOfWeek    time.Weekday
		weekdays       Weekdays
		activeDate     time.Time
		wantActiveDate time.Time
	}{
		{
			name:           "inital-sunday-next-monday",
			startOfWeek:    time.Sunday,
			weekdays:       Weekdays{time.Sunday: "", time.Monday: "", time.Wednesday: ""},
			activeDate:     time.Time{},
			wantActiveDate: newDate(2024, time.September, 23),
		},
		{
			name:           "inital-sunday-next-wednesday",
			startOfWeek:    time.Sunday,
			weekdays:       Weekdays{time.Sunday: "", time.Wednesday: "", time.Friday: ""},
			activeDate:     time.Time{},
			wantActiveDate: newDate(2024, time.September, 25),
		},
		{
			name:           "inital-tuesday-next-wednesday",
			startOfWeek:    time.Tuesday,
			weekdays:       Weekdays{time.Monday: "", time.Tuesday: "", time.Wednesday: ""},
			activeDate:     time.Time{},
			wantActiveDate: newDate(2024, time.September, 25),
		},
		{
			name:           "inital-tuesday-next-saturday",
			startOfWeek:    time.Tuesday,
			weekdays:       Weekdays{time.Sunday: "", time.Tuesday: "", time.Saturday: ""},
			activeDate:     time.Time{},
			wantActiveDate: newDate(2024, time.September, 28),
		},
		{
			name:           "inital-saturday-next-sunday",
			startOfWeek:    time.Saturday,
			weekdays:       Weekdays{time.Sunday: "", time.Wednesday: ""},
			activeDate:     time.Time{},
			wantActiveDate: newDate(2024, time.September, 22),
		},
		{
			name:           "tuesday-next-wednesday",
			startOfWeek:    time.Sunday,
			weekdays:       Weekdays{time.Tuesday: "", time.Wednesday: "", time.Friday: ""},
			activeDate:     newDate(2024, time.September, 24),
			wantActiveDate: newDate(2024, time.September, 25),
		},
		{
			name:           "wednesday-next-tuesday",
			startOfWeek:    time.Sunday,
			weekdays:       Weekdays{time.Tuesday: "", time.Wednesday: ""},
			activeDate:     newDate(2024, time.September, 25),
			wantActiveDate: newDate(2024, time.September, 24),
		},
		{
			name:           "saturday-next-sunday",
			startOfWeek:    time.Sunday,
			weekdays:       Weekdays{time.Sunday: "", time.Friday: "", time.Saturday: ""},
			activeDate:     newDate(2024, time.September, 28),
			wantActiveDate: newDate(2024, time.September, 22),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			sampleDate := newDate(2024, time.September, 24)
			tm := NewWeek(sampleDate).
				Weekdays(tt.weekdays).
				StartOfWeek(tt.startOfWeek)
			tm.activeDate = tt.activeDate

			// Test
			got := tm.NextDate()

			// Assertions
			assert.Equal(t, tt.wantActiveDate, got.activeDate)
		})
	}
}

func TestWeekModel_View(t *testing.T) {
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
			name:        "five-day-work-week",
			startOfWeek: time.Sunday,
			weekdays: Weekdays{
				time.Monday:    "",
				time.Tuesday:   "",
				time.Wednesday: "",
				time.Thursday:  "",
				time.Friday:    "",
			},
		},
		{
			name:        "weekend",
			startOfWeek: time.Sunday,
			weekdays: Weekdays{
				time.Friday:   "",
				time.Saturday: "",
				time.Sunday:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			sampleDate := newDate(2024, time.September, 24)

			tm := NewWeek(sampleDate).
				Weekdays(tt.weekdays).
				StartOfWeek(tt.startOfWeek)
			tm.days[newDate(2024, time.September, 24)] = testDayModel{}
			_ = tm.Init()

			// Test
			got := ansi.Strip(tm.View())

			// Assertions
			golden.RequireEqual(t, []byte(got))
		})
	}
}

func TestWeekModel_Update(t *testing.T) {
	tests := []struct {
		name           string
		msg            tea.Msg
		wantActiveDate time.Time
		wantCmd        tea.Cmd
	}{
		{
			name:           "keymsg-right",
			msg:            tea.KeyMsg{Type: tea.KeyRight},
			wantActiveDate: newDate(2024, time.September, 23),
			wantCmd:        func() tea.Msg { return ActiveDateMsg{Date: newDate(2024, time.September, 23)} },
		},
		{
			name:           "keymsg-left",
			msg:            tea.KeyMsg{Type: tea.KeyLeft},
			wantActiveDate: newDate(2024, time.September, 28),
			wantCmd:        func() tea.Msg { return ActiveDateMsg{Date: newDate(2024, time.September, 28)} },
		},
		{
			name:           "keymsg-ignore",
			msg:            tea.KeyMsg{Type: tea.KeyBackspace},
			wantActiveDate: time.Time{},
			wantCmd:        nil,
		},
		{
			name: "daycontentmsg",
			msg: DayContentMsg{
				Date:    newDate(2024, time.September, 23),
				Content: testDayModel{},
			},
			wantActiveDate: time.Time{},
			wantCmd:        nil,
		},
		{
			name: "daycontentmsg-wrong-week-early",
			msg: DayContentMsg{
				Date:    newDate(2024, time.September, 13),
				Content: testDayModel{},
			},
			wantActiveDate: time.Time{},
			wantCmd:        nil,
		},
		{
			name: "daycontentmsg-wrong-week-late",
			msg: DayContentMsg{
				Date:    newDate(2024, time.October, 13),
				Content: testDayModel{},
			},
			wantActiveDate: time.Time{},
			wantCmd:        nil,
		},
		{
			name:           "othermsg",
			msg:            tea.MouseMsg{},
			wantActiveDate: time.Time{},
			wantCmd:        nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			sampleDate := newDate(2024, time.September, 24)

			tm := NewWeek(sampleDate)
			_ = tm.Init()
			tm.days[newDate(2024, time.September, 24)] = testDayModel{}

			// Test
			got, gotCmd := tm.Update(tt.msg)

			// Assertions
			assert.Equal(t, tt.wantActiveDate, got.(WeekModel).activeDate)
			var wantMsg tea.Msg
			if tt.wantCmd != nil {
				wantMsg = tt.wantCmd()
			}
			var gotMsg tea.Msg
			if gotCmd != nil {
				gotMsg = gotCmd()
			}
			assert.Equal(t, wantMsg, gotMsg)
		})
	}
}
