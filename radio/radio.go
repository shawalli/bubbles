package radio

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type RadioButtonMsg struct {
	State bool
}

type ButtonStyles struct {
	LeftIndicatorCharacter        string
	RightIndicatorCharacter       string
	ActiveLeftIndicatorCharacter  string
	ActiveRightIndicatorCharacter string

	LeftIndicator        gloss.Style
	RightIndicator       gloss.Style
	ActiveLeftIndicator  gloss.Style
	ActiveRightIndicator gloss.Style

	Label       gloss.Style
	ActiveLabel gloss.Style
}

func DefaultButtonStyles() ButtonStyles {
	return ButtonStyles{
		LeftIndicatorCharacter:       "○",
		ActiveLeftIndicatorCharacter: "●",

		LeftIndicator:       gloss.NewStyle().Foreground(DefaultUnfocusedColor),
		ActiveLeftIndicator: gloss.NewStyle().Foreground(DefaultActiveTabIndicatorColor),
	}
}

type Button struct {
	label string

	active bool

	styles ButtonStyles
}

func NewButton(label string) Button {
	return Button{
		label:  label,
		styles: DefaultButtonStyles(),
	}
}

func (m Button) Styles(styles ButtonStyles) Button {
	m.styles = styles
	return m
}

func (m Button) Init() tea.Cmd {
	return nil
}

func (m Button) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case RadioButtonMsg:
		m.active = msg.State
	}
	return m, nil
}

func (m Button) View() string {
	lic := m.styles.LeftIndicatorCharacter
	ric := m.styles.RightIndicatorCharacter
	li := m.styles.LeftIndicator
	ri := m.styles.RightIndicator
	l := m.styles.Label

	if m.active {
		lic = m.styles.ActiveLeftIndicatorCharacter
		ric = m.styles.ActiveRightIndicatorCharacter
		li = m.styles.ActiveLeftIndicator
		ri = m.styles.ActiveRightIndicator
		l = m.styles.ActiveLabel
	}

	var left, right, label string
	if lic != "" {
		left = fmt.Sprintf("%s ", li.Render(lic))
	}
	if m.label != "" {
		label = l.Render(m.label)
	}
	if ric != "" {
		right = fmt.Sprintf(" %s", ri.Render(ric))
	}

	s := fmt.Sprintf("%s%s%s", left, label, right)

	return s
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

	// styles for the tab headers and bodies
	styles Styles
}

// New creates a new Model.
func New(vertical bool, buttons ...tea.Model) Model {
	// Activate first button
	b, _ := buttons[0].Update(RadioButtonMsg{true})
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

// Styles enables custom styling.
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
		n, c := m.buttons[i].Update(RadioButtonMsg{state})
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
	for _, button := range m.buttons {
		buttons = append(buttons, m.styles.Button.Render(button.View()))
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
