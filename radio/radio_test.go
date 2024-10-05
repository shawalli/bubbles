package radio

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/golden"
	"github.com/stretchr/testify/assert"
)

func Test_New_NoButtons(t *testing.T) {
	tests := []struct {
		name     string
		vertical bool
	}{
		{
			name:     "horizontal",
			vertical: false,
		},
		{
			name:     "vertical",
			vertical: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test
			got := New(false)

			// Assertions
			assert.Len(t, got.buttons, 0)
			assert.Zero(t, got.activeButton)
			assert.False(t, got.wraparound)
		})
	}
}

func Test_New(t *testing.T) {
	tests := []struct {
		name     string
		vertical bool
	}{
		{
			name:     "horizontal",
			vertical: false,
		},
		{
			name:     "vertical",
			vertical: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			buttons := []tea.Model{
				NewButton("a"),
				NewButton("b"),
			}

			// Test
			got := New(false, buttons...)

			// Assertions
			assert.Len(t, got.buttons, 2)
			assert.Zero(t, got.activeButton)
			assert.False(t, got.wraparound)
		})
	}
}

func TestModel_Wraparound(t *testing.T) {
	// Setup
	buttons := []tea.Model{
		NewButton("a"),
		NewButton("b"),
	}

	// Test
	got := New(true, buttons...).Wraparound(true)

	// Assertions
	assert.True(t, got.wraparound)
}

func TestModel_Styles(t *testing.T) {
	// Setup
	buttons := []tea.Model{
		NewButton("a"),
		NewButton("b"),
	}

	testStyles := Styles{}

	// Test
	got := New(true, buttons...).Styles(testStyles)

	// Assertions
	assert.Equal(t, testStyles, got.styles)
}

func TestModel_SetButton(t *testing.T) {
	tests := []struct {
		name       string
		button     int
		wantButton int
	}{
		{
			name:       "first",
			button:     0,
			wantButton: 0,
		},
		{
			name:       "middle",
			button:     1,
			wantButton: 1,
		},
		{
			name:       "last",
			button:     2,
			wantButton: 2,
		},
		{
			name:       "underflow",
			button:     -1,
			wantButton: 0,
		},
		{
			name:       "overflow",
			button:     3,
			wantButton: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			buttons := []tea.Model{
				NewButton("a"),
				NewButton("b"),
				NewButton("c"),
			}

			// Test
			got := New(true, buttons...).SetButton(tt.button)

			// Assertions
			assert.Equal(t, tt.wantButton, got.activeButton)
		})
	}
}

func TestModel_NextButton(t *testing.T) {
	tests := []struct {
		name       string
		button     int
		wraparound bool
		wantButton int
	}{
		{
			name:       "first",
			button:     0,
			wantButton: 1,
		},
		{
			name:       "next",
			button:     1,
			wantButton: 2,
		},
		{
			name:       "overflow",
			button:     2,
			wantButton: 2,
		},
		{
			name:       "wraparound",
			button:     2,
			wraparound: true,
			wantButton: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			buttons := []tea.Model{
				NewButton("a"),
				NewButton("b"),
				NewButton("c"),
			}

			// Test
			got := New(true, buttons...).Wraparound(tt.wraparound).SetButton(tt.button)
			got = got.NextButton()

			// Assertions
			assert.Equal(t, tt.wantButton, got.activeButton)
		})
	}
}

func TestModel_PreviousButton(t *testing.T) {
	tests := []struct {
		name       string
		button     int
		wraparound bool
		wantButton int
	}{
		{
			name:       "last",
			button:     2,
			wantButton: 1,
		},
		{
			name:       "previous",
			button:     1,
			wantButton: 0,
		},
		{
			name:       "underflow",
			button:     0,
			wantButton: 0,
		},
		{
			name:       "wraparound",
			button:     0,
			wraparound: true,
			wantButton: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			buttons := []tea.Model{
				NewButton("a"),
				NewButton("b"),
				NewButton("c"),
			}

			// Test
			got := New(true, buttons...).Wraparound(tt.wraparound).SetButton(tt.button)
			got = got.PreviousButton()

			// Assertions
			assert.Equal(t, tt.wantButton, got.activeButton)
		})
	}
}

func TestModel_Update(t *testing.T) {
	tests := []struct {
		name         string
		vertical     bool
		activeButton int
		wraparound   bool
		testMsg      tea.Msg
	}{
		{
			name:         "horizontal-previous-button",
			activeButton: 1,
			testMsg:      tea.KeyMsg{Type: tea.KeyLeft},
		},
		{
			name:    "horizontal-previous-button-lower-bound",
			testMsg: tea.KeyMsg{Type: tea.KeyLeft},
		},
		{
			name:       "horizontal-previous-button-lower-bound-wraparound",
			wraparound: true,
			testMsg:    tea.KeyMsg{Type: tea.KeyLeft},
		},
		{
			name:    "horizontal-next-button",
			testMsg: tea.KeyMsg{Type: tea.KeyRight},
		},
		{
			name:         "horizontal-next-button-upper-bound",
			activeButton: 2,
			testMsg:      tea.KeyMsg{Type: tea.KeyRight},
		},
		{
			name:         "horizontal-next-button-upper-bound-wraparound",
			activeButton: 2,
			wraparound:   true,
			testMsg:      tea.KeyMsg{Type: tea.KeyRight},
		},

		{
			name:         "vertical-previous-button",
			activeButton: 1,
			vertical:     true,
			testMsg:      tea.KeyMsg{Type: tea.KeyUp},
		},
		{
			name:     "vertical-previous-button-lower-bound",
			vertical: true,
			testMsg:  tea.KeyMsg{Type: tea.KeyUp},
		},
		{
			name:       "vertical-previous-button-lower-bound-wraparound",
			wraparound: true,
			vertical:   true,
			testMsg:    tea.KeyMsg{Type: tea.KeyUp},
		},
		{
			name:     "vertical-next-button",
			vertical: true,
			testMsg:  tea.KeyMsg{Type: tea.KeyDown},
		},
		{
			name:         "vertical-next-button-upper-bound",
			activeButton: 2,
			vertical:     true,
			testMsg:      tea.KeyMsg{Type: tea.KeyDown},
		},
		{
			name:         "vertical-next-button-upper-bound-wraparound",
			activeButton: 2,
			wraparound:   true,
			vertical:     true,
			testMsg:      tea.KeyMsg{Type: tea.KeyDown},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			testButtons := []tea.Model{
				NewButton("a"),
				NewButton("b"),
				NewButton("c"),
			}

			tm := New(tt.vertical, testButtons...).
				Wraparound(tt.wraparound).
				Styles(DefaultGroupedStyles(tt.vertical)).
				SetButton(tt.activeButton)
			// Because we are calling SetButton directly, without use key messages, manually update the active button
			b, _ := tm.buttons[0].Update(RadioButtonMsg{State: false})
			tm.buttons[0] = b
			b, _ = tm.buttons[tt.activeButton].Update(RadioButtonMsg{State: true})
			tm.buttons[tt.activeButton] = b
			_ = tm.Init()

			// Test
			v, _ := tm.Update(tt.testMsg)
			tm = v.(Model)
			got := ansi.Strip(tm.View())

			// Assertions
			golden.RequireEqual(t, []byte(got))
		})
	}
}
