package tabs

import tea "github.com/charmbracelet/bubbletea"

type Tab struct {
	title string
	child tea.Model
}

func NewTab(title string, child tea.Model) Tab {
	return Tab{
		title: title,
		child: child,
	}
}
