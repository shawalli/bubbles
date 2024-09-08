package tabs

import gloss "github.com/charmbracelet/lipgloss"

// Styles for tab rendering.
type Styles struct {
	// Inactive tab header
	Tab gloss.Style

	// Active tab header
	ActiveTab gloss.Style

	// Gap between rightmost tab header and right side of screen
	TabSpacer gloss.Style

	// Area to left and right of active tab header title
	TabIndicator gloss.Style

	// Character(s) to left of active tab header title
	TabIndicatorLeft string

	// Character(s) to right of active tab header title
	TabIndicatorRight string

	// Tab content
	TabWindow gloss.Style
}

// DefaultStyles provides default tab styles.
func DefaultStyles() Styles {
	return Styles{
		Tab: gloss.NewStyle().
			Foreground(DefaultUnfocusedColor).
			Border(DefaultTabBorder, true).
			BorderForeground(DefaultForegroundColor).
			Padding(0, 1),
		ActiveTab: gloss.NewStyle().
			Foreground(DefaultUnfocusedColor).
			Border(DefaultActiveTabBorder, true).
			BorderForeground(DefaultForegroundColor).
			Padding(0, 1).
			Bold(true),
		TabSpacer: gloss.NewStyle().
			Border(DefaultTabSpacerBorder, false, true, true, false).
			BorderForeground(DefaultForegroundColor).
			Padding(0, 1),
		TabIndicator: gloss.NewStyle().
			Foreground(DefaultActiveTabIndicatorColor).
			Bold(true),
		TabIndicatorLeft:  "=",
		TabIndicatorRight: "=",
		TabWindow: gloss.NewStyle().
			Border(DefaultWindowBorder, true).
			BorderForeground(DefaultForegroundColor).
			Padding(0, 1),
	}
}

var (
	DefaultForegroundColor         = gloss.AdaptiveColor{Light: "#874Bfd", Dark: "#7d56f4"}
	DefaultUnfocusedColor          = gloss.AdaptiveColor{Light: "#3a3a3a", Dark: "#b0b0b0"}
	DefaultActiveTabIndicatorColor = gloss.AdaptiveColor{Light: "#bb99fe", Dark: "#997bf6"}

	DefaultWindowBorder = gloss.Border{
		Top:         " ",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "│",
		TopRight:    "│",
		BottomLeft:  "└",
		BottomRight: "┘",
	}

	DefaultTabBorder = gloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭'",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}
	DefaultActiveTabBorder = gloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭'",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	DefaultTabSpacerBorder = gloss.Border{
		Bottom:      "─",
		BottomRight: "┐",
	}
)
