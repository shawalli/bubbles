package radio

import gloss "github.com/charmbracelet/lipgloss"

// Styles for button rendering.
//
// [DefaultStyles] and [DefaultPillStyles] only really use Button. [DefaultGroupedStyles] uses all
// three button styles for grouped rendering.
type Styles struct {
	FirstButton gloss.Style
	Button      gloss.Style
	LastButton  gloss.Style
}

// DefaultStyles provides default button styles.
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

// DefaultPillStyles provides default button pill styles.
func DefaultPillStyles(vertical bool) Styles {
	b := gloss.NewStyle().
		Border(DefaultPillBorder, true).
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

// DefaultGroupedStyles provides default button grouped styles.
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

// Styles for rendering the interior of the button itself.
type ButtonStyles struct {
	// Character(s) to left of the button label/content
	LeftIndicatorCharacter string

	// Character(s) to left of the action button label/content
	ActiveLeftIndicatorCharacter string

	// Character(s) to right of the button label/content
	RightIndicatorCharacter string

	// Character(s) to right of the active button label/content
	ActiveRightIndicatorCharacter string

	// Area to the left of the button label/content
	LeftIndicator gloss.Style

	// Area to the left of the active button label/content
	ActiveLeftIndicator gloss.Style

	// Area to the right of the button label/content
	RightIndicator gloss.Style

	// Area to the right of the active button label/content
	ActiveRightIndicator gloss.Style

	// Label/content for the button
	Label gloss.Style

	// Label/content for the active button
	ActiveLabel gloss.Style
}

// DefaultButtonStyles provides default styles for a button.
func DefaultButtonStyles() ButtonStyles {
	return ButtonStyles{
		LeftIndicatorCharacter:       "○",
		ActiveLeftIndicatorCharacter: "●",

		LeftIndicator:       gloss.NewStyle().Foreground(DefaultUnfocusedColor),
		ActiveLeftIndicator: gloss.NewStyle().Foreground(DefaultActiveButtonIndicatorColor),
	}
}

var (
	DefaultUnfocusedColor             = gloss.AdaptiveColor{Light: "#3a3a3a", Dark: "#b0b0b0"}
	DefaultActiveButtonIndicatorColor = gloss.AdaptiveColor{Light: "#bb99fe", Dark: "#997bf6"}

	// ╭───╮
	// │foo│
	// ╰───╯
	DefaultPillBorder = gloss.Border{
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
