package tabs

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
)

// KeyMap contains relevant keys for tab navigation.
type KeyMap struct {
	TabLeft  key.Binding
	TabRight key.Binding

	TabNumbers []key.Binding
}

// DefaultKeyMap contains default key mappings for tab navigation.
func DefaultKeyMap(tabs int) KeyMap {
	km := KeyMap{
		TabLeft: key.NewBinding(
			key.WithKeys("shift+tab", "ctrl+left"),
			key.WithHelp("shift+tab/ctrl+left", "<tab"),
		),
		TabRight: key.NewBinding(
			key.WithKeys("tab", "ctrl+right"),
			key.WithHelp("tab/ctrl+right", "tab>"),
		),
	}

	for i := 1; i <= tabs; i++ {
		k := fmt.Sprintf("%d", i)
		tkb := key.NewBinding(
			key.WithKeys(k),
			key.WithHelp(k, fmt.Sprintf("tab %d", i)),
		)
		km.TabNumbers = append(km.TabNumbers, tkb)
	}

	return km
}
