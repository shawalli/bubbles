package tabs

import tea "github.com/charmbracelet/bubbletea"

// Tab is the basic model for a single tab.
type Tab struct {
	// title is the tab title
	title string

	// child is the tab window body.
	child tea.Model
}

// NewTab creates a new tab.
func NewTab(title string, child tea.Model) Tab {
	return Tab{
		title: title,
		child: child,
	}
}
