package radio

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

// RadioButtonMsg enables communication down into radio buttons
type RadioButtonMsg struct {
	// State is whether or not the radio button is active
	State bool

	// User-defined data to pass down into the button
	Data any
}

// Model represents a group of radio buttons.
type Model struct {
	// keyMap is key bindings for button navigation
	keyMap KeyMap

	// button labels
	buttons []tea.Model

	// activeButton is the index of the button currently being selected
	activeButton int

	vertical bool

	// wraparound enables wrapping around from the last button to the first button or from the first button
	// back to the last button
	wraparound bool

	// styles for the radio buttons
	styles Styles
}

// New creates a new [Model].
func New(vertical bool, buttons ...tea.Model) Model {
	// Activate first button
	b, _ := buttons[0].Update(RadioButtonMsg{State: true})
	buttons[0] = b

	m := Model{
		keyMap:   DefaultKeyMap(vertical),
		buttons:  buttons,
		vertical: vertical,
		styles:   DefaultStyles(vertical),
	}

	return m
}

// Wraparound enables or disables wraparound navigation from the last button to the first button and
// from the first button back to the last button.
func (m Model) Wraparound(w bool) Model {
	m.wraparound = w
	return m
}

// Styles sets custom styling.
func (m Model) Styles(styles Styles) Model {
	m.styles = styles
	return m
}

// SetButton sets the active button.
func (m Model) SetButton(i int) Model {
	if i < 0 {
		i = 0
	} else if i >= len(m.buttons) {
		i = len(m.buttons) - 1
	}
	m.activeButton = i

	return m
}

// NextButton moves the active button forward.
func (m Model) NextButton() Model {
	i := m.activeButton + 1

	if i == len(m.buttons) {
		if m.wraparound {
			i = 0
		} else {
			i--
		}
	}

	return m.SetButton(i)
}

// PreviousButton moves the active button backward.
func (m Model) PreviousButton() Model {
	i := m.activeButton - 1

	if i < 0 {
		if m.wraparound {
			i = len(m.buttons) - 1
		} else {
			i = 0
		}
	}

	return m.SetButton(i)
}

// Init the Model.
func (m Model) Init() tea.Cmd { return nil }

// Update the Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	signalButton := func(i int, state bool) tea.Cmd {
		n, c := m.buttons[i].Update(RadioButtonMsg{State: state})
		m.buttons[i] = n

		return c
	}

	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Previous):
			cmds = append(cmds, signalButton(m.activeButton, false))
			m = m.PreviousButton()
			cmds = append(cmds, signalButton(m.activeButton, true))
		case key.Matches(msg, m.keyMap.Next):
			cmds = append(cmds, signalButton(m.activeButton, false))
			m = m.NextButton()
			cmds = append(cmds, signalButton(m.activeButton, true))
		}
	}

	return m, tea.Sequence(cmds...)
}

// View renders the Model.
func (m Model) View() string {
	doc := strings.Builder{}

	var buttons []string
	for i, label := range m.buttons {
		var buttonStyle gloss.Style
		switch i {
		case 0:
			buttonStyle = m.styles.FirstButton
		case len(m.buttons) - 1:
			buttonStyle = m.styles.LastButton
		default:
			buttonStyle = m.styles.Button
		}

		buttons = append(buttons, buttonStyle.Render(label.View()))
	}

	var row string
	if m.vertical {
		row = gloss.JoinVertical(gloss.Top, buttons...)
	} else {
		row = gloss.JoinHorizontal(gloss.Top, buttons...)
	}

	doc.WriteString(row + "\n")

	return doc.String()
}
