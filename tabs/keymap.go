package tabs

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	TabLeft  key.Binding
	TabRight key.Binding

	TabNumbers []key.Binding
}

func (km KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{km.TabLeft, km.TabRight}
}

func (km KeyMap) FullHelp() [][]key.Binding {
	help := [][]key.Binding{km.ShortHelp()}

	var tnHelp []key.Binding
	for _, tn := range km.TabNumbers {
		tnHelp = append(tnHelp, tn)
	}
	if len(tnHelp) > 0 {
		help = append(help, tnHelp)
	}
	return help
}

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
