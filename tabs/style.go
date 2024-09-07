package tabs

import gloss "github.com/charmbracelet/lipgloss"

type Styles struct {
	Tab          gloss.Style
	ActiveTab    gloss.Style
	TabSpacer    gloss.Style
	TabIndicator gloss.Style

	TabWindow gloss.Style
}

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
