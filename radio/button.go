package radio

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Button represents a single radio button.
// It provides basic handling of radio-button-state. However, a user-defined model may be passed into the [Model]
// instead if custom behavior is needed.
type Button struct {
	// label to display
	label string

	// active is whether or not the button is the active button
	active bool

	// styles contains styles used for rendering the button label/content
	styles ButtonStyles
}

// NewButton creates a new [Button].
func NewButton(label string) Button {
	return Button{
		label:  label,
		styles: DefaultButtonStyles(),
	}
}

// Styles enables custom styling.
func (m Button) Styles(styles ButtonStyles) Button {
	m.styles = styles
	return m
}

// Init the Button.
func (m Button) Init() tea.Cmd {
	return nil
}

// Update the Button.
func (m Button) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case RadioButtonMsg:
		m.active = msg.State
	}
	return m, nil
}

// View renders the Button.
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
