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
	// Days-of-the-week header
	LeftHeaderStyle   gloss.Style
	MiddleHeaderStyle gloss.Style
	RightHeaderStyle  gloss.Style

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
		DateStyles: DefaultDateStyles(),

		LeftHeaderStyle: gloss.NewStyle().
			Border(DefaultLeftHeaderBorder, true).
			Padding(0, 1),
		MiddleHeaderStyle: gloss.NewStyle().
			Border(DefaultMiddleHeaderBorder, true, true, true, false).
			Padding(0, 1),
		RightHeaderStyle: gloss.NewStyle().
			Border(DefaultRightHeaderBorder, true, true, true, false).
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

// Styles for rendering the calendar.
type WeekStyles struct {
	// Days-of-the-week header
	LeftHeaderStyle   gloss.Style
	MiddleHeaderStyle gloss.Style
	RightHeaderStyle  gloss.Style

	// Style to indicate a day is the active day from the header
	ActiveHeaderStyle gloss.Style

	// Month block bottom row
	LeftDayStyle   gloss.Style
	MiddleDayStyle gloss.Style
	RightDayStyle  gloss.Style

	// Date interior
	//
	// Note: NumberStyle and ActiveNumber styles are ignored for WeekModel.
	DateStyles DateStyles
	DateFormat string
}

func DefaultWeekStyles() WeekStyles {
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

		LeftDayStyle: gloss.NewStyle().
			Border(DefaultBottomLeftDayBorder, false, true, true, true).
			Align(gloss.Center, gloss.Top).
			Padding(1, 0),
		MiddleDayStyle: gloss.NewStyle().
			Border(DefaultBottomDayBorder, false, true, true, false).
			Align(gloss.Center, gloss.Top).
			Padding(1, 0),
		RightDayStyle: gloss.NewStyle().
			Border(DefaultBottomRightDayBorder, false, true, true, false).
			Align(gloss.Center, gloss.Top).
			Padding(1, 0),

		DateStyles: DateStyles{
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
	DefaultLeftHeaderBorder = gloss.Border{
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
	DefaultMiddleHeaderBorder = gloss.Border{
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
	//  ┼───┤
	DefaultRightHeaderBorder = gloss.Border{
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
