package calendar

import (
	gloss "github.com/charmbracelet/lipgloss"
)

// Styles for date block rendering.
type DateStyles struct {
	// Width of the date block
	Width int

	// Height of the date block
	Height int

	// Date number style
	NumberStyle gloss.Style

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

var (
	DefaultUnfocusedColor             = gloss.AdaptiveColor{Light: "#3a3a3a", Dark: "#b0b0b0"}
	DefaultActiveButtonIndicatorColor = gloss.AdaptiveColor{Light: "#bb99fe", Dark: "#997bf6"}

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
	//  ┬───╮
	//  │Sat│
	//  ┼───┼
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
