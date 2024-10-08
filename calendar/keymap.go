package calendar

import (
	"github.com/charmbracelet/bubbles/key"
)

// KeyMap contains relevant keys for tab navigation.
type KeyMap struct {
	Left  key.Binding
	Right key.Binding
	Up    key.Binding
	Down  key.Binding
}

// DefaultMonthKeyMap contains default key mappings for monthly navigation.
func DefaultMonthKeyMap() KeyMap {
	return KeyMap{
		Left:  key.NewBinding(key.WithKeys("left"), key.WithHelp("left", "←")),
		Right: key.NewBinding(key.WithKeys("right"), key.WithHelp("right", "→")),
		Up:    key.NewBinding(key.WithKeys("up"), key.WithHelp("up", "↑")),
		Down:  key.NewBinding(key.WithKeys("down"), key.WithHelp("down", "↓")),
	}
}

// DefaultWeekKeyMap contains default key mappings for weekly navigation.
func DefaultWeekKeyMap() KeyMap {
	return KeyMap{
		Left:  key.NewBinding(key.WithKeys("left"), key.WithHelp("left", "←")),
		Right: key.NewBinding(key.WithKeys("right"), key.WithHelp("right", "→")),
	}
}
