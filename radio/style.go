package radio

import gloss "github.com/charmbracelet/lipgloss"

// Styles for tab rendering.
type Styles struct {
	FirstButton gloss.Style
	Button      gloss.Style
	LastButton  gloss.Style
}

func DefaultStyles(vertical bool) Styles {
	b := gloss.NewStyle().PaddingRight(2)
	s := Styles{
		FirstButton: b,
		Button:      b,
		LastButton:  b,
	}

	if vertical {
		b = s.Button.UnsetPadding()
		s = Styles{
			FirstButton: b,
			Button:      b,
			LastButton:  b,
		}
	}

	return s
}

func DefaultPillboxStyles(vertical bool) Styles {
	b := gloss.NewStyle().
		Border(DefaultPillboxBorder, true).
		Padding(0, 2, 0, 2)
	s := Styles{
		FirstButton: b,
		Button:      b,
		LastButton:  b,
	}

	if vertical {
		b = s.Button.UnsetPadding()
		s = Styles{
			FirstButton: b,
			Button:      b,
			LastButton:  b,
		}
	}

	return s
}

func DefaultGroupedStyles(vertical bool) Styles {
	s := Styles{
		FirstButton: gloss.NewStyle().
			Border(DefaultGroupedHorizontalFirstBorder, true),
		Button: gloss.NewStyle().
			Border(DefaultGroupedHorizontalBorder, true, true, true, false),
		LastButton: gloss.NewStyle().
			Border(DefaultGroupedHorizontalLastBorder, true, true, true, false),
	}

	if vertical {
		s = Styles{
			FirstButton: gloss.NewStyle().
				Border(DefaultGroupedVerticalFirstBorder, true),
			Button: gloss.NewStyle().
				Border(DefaultGroupedVerticalBorder, false, true, true, true),
			LastButton: gloss.NewStyle().
				Border(DefaultGroupedVerticalLastBorder, false, true, true, true),
		}
	}

	return s
}

var (
	DefaultUnfocusedColor          = gloss.AdaptiveColor{Light: "#3a3a3a", Dark: "#b0b0b0"}
	DefaultActiveTabIndicatorColor = gloss.AdaptiveColor{Light: "#bb99fe", Dark: "#997bf6"}

	// ╭───╮
	// │foo│
	// ╰───╯
	DefaultPillboxBorder = gloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}

	// ╭───┬
	// │foo│
	// ╰───┴
	DefaultGroupedHorizontalFirstBorder = gloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "┬",
		BottomLeft:  "╰",
		BottomRight: "┴",
	}
	//  ───┬
	//  foo│
	//  ───┴
	DefaultGroupedHorizontalBorder = gloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "",
		Right:       "│",
		TopLeft:     "",
		TopRight:    "┬",
		BottomLeft:  "",
		BottomRight: "┴",
	}
	//  ───╮
	//  foo│
	//  ───╯
	DefaultGroupedHorizontalLastBorder = gloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "",
		Right:       "│",
		TopLeft:     "",
		TopRight:    "╮",
		BottomLeft:  "",
		BottomRight: "╯",
	}

	// ┌───┐
	// │foo│
	// ├───┤
	DefaultGroupedVerticalFirstBorder = gloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "├",
		BottomRight: "┤",
	}
	// │foo│
	// ├───┤
	DefaultGroupedVerticalBorder = gloss.Border{
		Top:         "",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "",
		TopRight:    "",
		BottomLeft:  "├",
		BottomRight: "┤",
	}
	// │foo│
	// ╰───╯
	DefaultGroupedVerticalLastBorder = gloss.Border{
		Top:         "",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "",
		TopRight:    "",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}
)
