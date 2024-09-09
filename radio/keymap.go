package radio

import (
	"github.com/charmbracelet/bubbles/key"
)

// KeyMap contains relevant keys for tab navigation.
type KeyMap struct {
	Select key.Binding

	Previous key.Binding
	Next     key.Binding
}

// defaultKeyMap contains default key mappings for radio group navigation.
func DefaultKeyMap(vertical bool) KeyMap {
	km := KeyMap{
		Select: key.NewBinding(
			key.WithKeys("enter", "space"),
		),
		Previous: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("left", "←"),
		),
		Next: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("right", "→"),
		),
	}

	if vertical {
		km.Previous = key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("up", "↑"),
		)
		km.Next = key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("down", "↓"),
		)
	}

	return km
}
