package calendar

import (
	gloss "github.com/charmbracelet/lipgloss"
)

type CellStyles struct {
	// Width of the date block
	Width int

	// Height of the date block
	Height int

	// Contents style
	BodyStyle gloss.Style
}

// Styles for date block rendering.
type DateStyles struct {
	// Width of the date block
	Width int

	// Height of the date block
	Height int

	// Date number style
	NumberStyle gloss.Style

	ActiveNumberStyle gloss.Style

	// Contents style
	BodyStyle gloss.Style
}

// DefaultStyles provides default styles for the date block.
func DefaultDateStyles() DateStyles {
	// Default days-of-the-week labels are 3 characters, so this is 3-characters and 1-character
	// of left/right padding.
	defaultWidth := 5
	defaultHeight := 2

	return DateStyles{
		Width:  defaultWidth,
		Height: defaultHeight,

		NumberStyle: gloss.NewStyle().
			Width(defaultWidth).
			Align(gloss.Left),
		ActiveNumberStyle: gloss.NewStyle().
			Width(defaultWidth).
			Align(gloss.Left).
			Bold(true).
			Foreground(DefaultActiveColor),
		BodyStyle: gloss.NewStyle().
			Width(defaultWidth).
			Height(defaultHeight - 1).
			Align(gloss.Center),
	}
}

// Styles for rendering the calendar.
type MonthStyles struct {
	// Slice of labels to use for the days-of-the-week header
	// NOTE: Order of this slice should correspond to the start of the week
	DaysOfWeek []string

	// Days-of-the-week header
	LeftDaysOfWeekStyle   gloss.Style
	MiddleDaysOfWeekStyle gloss.Style
	RightDaysOfWeekStyle  gloss.Style

	// Month block interior
	MiddleLeftDayStyle  gloss.Style
	MiddleDayStyle      gloss.Style
	MiddleRightDayStyle gloss.Style

	// Month block bottom row
	BottomLeftDayStyle  gloss.Style
	BottomDayStyle      gloss.Style
	BottomRightDayStyle gloss.Style

	// Date interior
	DateStyles DateStyles
}

// DefaultMonthStyles provides default month styles.
func DefaultMonthStyles() MonthStyles {
	return MonthStyles{
		DaysOfWeek: DefaultDaysOfWeek(),

		DateStyles: DefaultDateStyles(),

		LeftDaysOfWeekStyle: gloss.NewStyle().
			Border(DefaultWeekdayFirstBorder, true).
			Padding(0, 1),
		MiddleDaysOfWeekStyle: gloss.NewStyle().
			Border(DefaultWeekdayBorder, true, true, true, false).
			Padding(0, 1),
		RightDaysOfWeekStyle: gloss.NewStyle().
			Border(DefaultWeekdayLastBorder, true, true, true, false).
			Padding(0, 1),

		MiddleLeftDayStyle: gloss.NewStyle().
			Border(DefaultMiddleLeftDayBorder, false, true, true, true),
		MiddleDayStyle: gloss.NewStyle().
			Border(DefaultMiddleDayBorder, false, true, true, false),
		MiddleRightDayStyle: gloss.NewStyle().
			Border(DefaultMiddleRightDayBorder, false, true, true, false),

		BottomLeftDayStyle: gloss.NewStyle().
			Border(DefaultBottomLeftDayBorder, false, true, true, true),
		BottomDayStyle: gloss.NewStyle().
			Border(DefaultBottomDayBorder, false, true, true, false),
		BottomRightDayStyle: gloss.NewStyle().
			Border(DefaultBottomRightDayBorder, false, true, true, false),
	}
}

// DefaultDaysOfWeek provides default days-of-the-week shortened values.
func DefaultDaysOfWeek() []string {
	return []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
}

// DefaultDaysOfWeek provides default days-of-the-week abbreviated values.
func DefaultDaysOfWeekShort() []string {
	return []string{"U", "M", "T", "W", "R", "F", "S"}
}

// Styles for rendering the calendar.
type WeekStyles struct {
	// Days-of-the-week header
	LeftHeaderStyle   gloss.Style
	MiddleHeaderStyle gloss.Style
	RightHeaderStyle  gloss.Style

	// Style to indicate a day is the active day from the header
	ActiveHeaderStyle gloss.Style

	// Month block bottom row
	LeftCellStyle   gloss.Style
	MiddleCellStyle gloss.Style
	RightCellStyle  gloss.Style

	// Date interior
	CellStyles CellStyles
	DateFormat string
}

func DefaultWeekStyles() WeekStyles {
	// Default days-of-the-week labels are 3 characters, so this is 3-characters and 1-character
	// of left/right padding.
	defaultWidth := 15
	defaultHeight := 5

	return WeekStyles{
		LeftHeaderStyle: gloss.NewStyle().
			Border(DefaultLeftHeaderBorder, true).
			Align(gloss.Center).
			Padding(1, 0),
		MiddleHeaderStyle: gloss.NewStyle().
			Border(DefaultMiddleHeaderBorder, true, true, true, false).
			Align(gloss.Center).
			Padding(1, 0),
		RightHeaderStyle: gloss.NewStyle().
			Border(DefaultRightHeaderBorder, true, true, true, false).
			Align(gloss.Center).
			Padding(1, 0),

		ActiveHeaderStyle: gloss.NewStyle().
			Align(gloss.Center).
			Bold(true).
			Foreground(DefaultActiveColor),

		LeftCellStyle: gloss.NewStyle().
			Border(DefaultBottomLeftDayBorder, false, true, true, true).
			Align(gloss.Center, gloss.Top).
			Padding(1, 0),
		MiddleCellStyle: gloss.NewStyle().
			Border(DefaultBottomDayBorder, false, true, true, false).
			Align(gloss.Center, gloss.Top).
			Padding(1, 0),
		RightCellStyle: gloss.NewStyle().
			Border(DefaultBottomRightDayBorder, false, true, true, false).
			Align(gloss.Center, gloss.Top).
			Padding(1, 0),

		CellStyles: CellStyles{
			Width:  defaultWidth,
			Height: defaultHeight,

			BodyStyle: gloss.NewStyle().
				Width(defaultWidth).
				Height(defaultHeight - 1).
				Align(gloss.Center),
		},
		DateFormat: "1/02",
	}
}

var (
	DefaultActiveColor = gloss.AdaptiveColor{Light: "#3E5AFA", Dark: "#7DD6FA"}

	// ╭───┬
	// │Sun│
	// ├───┼
	DefaultWeekdayFirstBorder = gloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "┬",
		BottomLeft:  "├",
		BottomRight: "┼",
	}
	DefaultLeftHeaderBorder = DefaultWeekdayFirstBorder
	//  ┬───┬
	//  │Mon│
	//  ┼───┼
	DefaultWeekdayBorder = gloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "┬",
		TopRight:    "┬",
		BottomLeft:  "┼",
		BottomRight: "┼",
	}
	DefaultMiddleHeaderBorder = DefaultWeekdayBorder
	//  ┬───╮
	//  │Sat│
	//  ┼───┤
	DefaultWeekdayLastBorder = gloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "┬",
		TopRight:    "╮",
		BottomLeft:  "┼",
		BottomRight: "┤",
	}
	DefaultRightHeaderBorder = DefaultWeekdayLastBorder

	//  ├───┼
	//  │12 │
	//  ├───┼
	DefaultMiddleLeftDayBorder = gloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "├",
		TopRight:    "┼",
		BottomLeft:  "├",
		BottomRight: "┼",
	}
	//  ┼───┼
	//  │12 │
	//  ┼───┼
	DefaultMiddleDayBorder = gloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "┼",
		TopRight:    "┼",
		BottomLeft:  "┼",
		BottomRight: "┼",
	}

	//  ┼───┤
	//  │12 │
	//  ┼───┤
	DefaultMiddleRightDayBorder = gloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "┼",
		TopRight:    "┤",
		BottomLeft:  "┼",
		BottomRight: "┤",
	}

	//  ├───┼
	//  │12 │
	//  ╰───┴
	DefaultBottomLeftDayBorder = gloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "├",
		TopRight:    "┼",
		BottomLeft:  "╰",
		BottomRight: "┴",
	}
	//  ┼───┼
	//  │12 │
	//  ┴───┴
	DefaultBottomDayBorder = gloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "┼",
		TopRight:    "┼",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	//  ┼───┤
	//  │12 │
	//  ┴───╯
	DefaultBottomRightDayBorder = gloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "┼",
		TopRight:    "┤",
		BottomLeft:  "┴",
		BottomRight: "╯",
	}
)
