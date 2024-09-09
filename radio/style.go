package radio

import gloss "github.com/charmbracelet/lipgloss"

// Styles for tab rendering.
type Styles struct {
	Button gloss.Style
}

func DefaultStyles(vertical bool) Styles {
	s := Styles{
		Button: gloss.NewStyle().PaddingRight(2),
	}

	if vertical {
		s.Button = s.Button.UnsetPadding()
	}

	return s
}

func DefaultPillboxStyles(vertical bool) Styles {
	s := Styles{
		Button: gloss.NewStyle().
			Border(DefaultPillboxBorder, true).
			Padding(0, 2, 0, 2),
	}

	if vertical {
		s.Button = s.Button.UnsetPadding()
	}

	return s
}

var (
	DefaultUnfocusedColor          = gloss.AdaptiveColor{Light: "#3a3a3a", Dark: "#b0b0b0"}
	DefaultActiveTabIndicatorColor = gloss.AdaptiveColor{Light: "#bb99fe", Dark: "#997bf6"}

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
)
