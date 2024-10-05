package tabs

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/golden"
	"github.com/stretchr/testify/assert"
)

const (
	testTermWidth  = 80
	testTermHeight = 20
)

func mockTerm(t *testing.T) {
	t.Helper()

	origGetTermSize := getTermSize
	t.Cleanup(func() {
		getTermSize = origGetTermSize
	})
	getTermSize = func(fd uintptr) (width int, height int, err error) {
		return testTermWidth, testTermHeight, nil
	}
}

func Test_New(t *testing.T) {
	// Setup
	mockTerm(t)

	// Test
	got := New()

	// Assertions
	assert.Len(t, got.tabs, 0)
	assert.Zero(t, got.activeTab)
	assert.Equal(t, got.width, testTermWidth-2)
	assert.Equal(t, got.height, testTermHeight)
	assert.False(t, got.wraparound)
}

func TestModel_Wraparound(t *testing.T) {
	// Setup
	mockTerm(t)

	// Test
	got := New().Wraparound(true)

	// Assertions
	assert.True(t, got.wraparound)
}

func TestModel_Styles(t *testing.T) {
	// Setup
	mockTerm(t)

	testStyles := Styles{}

	// Test
	got := New().Styles(testStyles)

	// Assertions
	assert.Equal(t, testStyles, got.styles)
}

func TestModel_Width(t *testing.T) {
	// Setup
	mockTerm(t)

	// Test
	got := New().Width(10)

	// Assertions
	assert.Equal(t, 10, got.width)
	assert.Equal(t, testTermHeight, got.height)
}

func TestModel_Height(t *testing.T) {
	// Setup
	mockTerm(t)

	// Test
	got := New().Height(10)

	// Assertions
	assert.Equal(t, testTermWidth-2, got.width)
	assert.Equal(t, 10, got.height)
}

func TestModel_SetTab(t *testing.T) {
	tests := []struct {
		name         string
		testTabIndex int
		wantTabIndex int
	}{
		{
			name:         "standard",
			testTabIndex: 1,
			wantTabIndex: 1,
		},
		{
			name:         "same-value",
			testTabIndex: 2,
			wantTabIndex: 2,
		},
		{
			name:         "out-of-bounds-lower",
			testTabIndex: -1,
			wantTabIndex: 0,
		},
		{
			name:         "out-of-bounds-upper",
			testTabIndex: 4,
			wantTabIndex: 3,
		},
		{
			name:         "bound-lower",
			testTabIndex: 0,
			wantTabIndex: 0,
		},
		{
			name:         "bound-upper",
			testTabIndex: 3,
			wantTabIndex: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			testTabs := []Tab{
				NewTab("0", nil),
				NewTab("1", nil),
				NewTab("2", nil),
				NewTab("3", nil),
			}
			testModel := New(testTabs...)
			testModel.activeTab = 2

			// Test
			got := testModel.SetTab(tt.testTabIndex)

			// Assertions
			assert.Equal(t, tt.wantTabIndex, got.activeTab)
		})
	}
}

func TestModel_NextTab(t *testing.T) {
	tests := []struct {
		name         string
		testTabIndex int
		wantTabIndex int
		wraparound   bool
	}{
		{
			name:         "start",
			testTabIndex: 0,
			wantTabIndex: 1,
		},
		{
			name:         "iterate-1",
			testTabIndex: 1,
			wantTabIndex: 2,
		},
		{
			name:         "iterate-2",
			testTabIndex: 2,
			wantTabIndex: 3,
		},
		{
			name:         "bound-upper",
			testTabIndex: 3,
			wantTabIndex: 3,
		},
		{
			name:         "wraparound",
			testTabIndex: 3,
			wantTabIndex: 0,
			wraparound:   true,
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			testTabs := []Tab{
				NewTab("0", nil),
				NewTab("1", nil),
				NewTab("2", nil),
				NewTab("3", nil),
			}
			testModel := New(testTabs...).Wraparound(tt.wraparound)
			testModel.activeTab = tt.testTabIndex

			// Test
			got := testModel.NextTab()

			// Assertions
			assert.Equal(t, tt.wantTabIndex, got.activeTab)
		})
	}
}

func TestModel_PreviousTab(t *testing.T) {
	tests := []struct {
		name         string
		testTabIndex int
		wantTabIndex int
		wraparound   bool
	}{
		{
			name:         "end",
			testTabIndex: 3,
			wantTabIndex: 2,
		},
		{
			name:         "iterate-1",
			testTabIndex: 2,
			wantTabIndex: 1,
		},
		{
			name:         "iterate-2",
			testTabIndex: 1,
			wantTabIndex: 0,
		},
		{
			name:         "bound-lower",
			testTabIndex: 0,
			wantTabIndex: 0,
		},
		{
			name:         "wraparound",
			testTabIndex: 0,
			wantTabIndex: 3,
			wraparound:   true,
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			testTabs := []Tab{
				NewTab("0", nil),
				NewTab("1", nil),
				NewTab("2", nil),
				NewTab("3", nil),
			}
			testModel := New(testTabs...).Wraparound(tt.wraparound)
			testModel.activeTab = tt.testTabIndex

			// Test
			got := testModel.PreviousTab()

			// Assertions
			assert.Equal(t, tt.wantTabIndex, got.activeTab)
		})
	}
}

type testModel struct {
	content string
}

func (tm testModel) Init() tea.Cmd                           { return nil }
func (tm testModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return tm, nil }
func (tm testModel) View() string                            { return tm.content }

func TestModel_Update(t *testing.T) {
	tests := []struct {
		name          string
		testActiveTab int
		testMsg       tea.Msg
		wraparound    bool
	}{
		{
			name:          "next-tab",
			testActiveTab: 0,
			testMsg:       tea.KeyMsg{Type: tea.KeyTab},
		},
		{
			name:          "previous-tab-lower-bound",
			testActiveTab: 0,
			testMsg:       tea.KeyMsg{Type: tea.KeyShiftTab},
		},
		{
			name:          "tab-num-4",
			testActiveTab: 0,
			testMsg: tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune{'4'},
			},
		},
		{
			name:          "next-tab-upper-bound",
			testActiveTab: 3,
			testMsg:       tea.KeyMsg{Type: tea.KeyTab},
		},
		{
			name:          "next-tab-upper-bound-wraparound",
			testActiveTab: 3,
			testMsg:       tea.KeyMsg{Type: tea.KeyTab},
			wraparound:    true,
		},
		{
			name:          "previous-tab-lower-bound-wraparound",
			testActiveTab: 0,
			testMsg:       tea.KeyMsg{Type: tea.KeyShiftTab},
			wraparound:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			testTabs := []Tab{
				NewTab("1", testModel{"I am tab 1"}),
				NewTab("2", testModel{"I am tab 2"}),
				NewTab("3", testModel{"I am tab 3"}),
				NewTab("4", testModel{"I am tab 4"}),
			}
			tm := New(testTabs...).
				Wraparound(tt.wraparound).
				SetTab(tt.testActiveTab)

			// Test
			v, _ := tm.Update(tt.testMsg)
			tm = v.(Model)
			got := ansi.Strip(tm.View())

			// Assertions
			golden.RequireEqual(t, []byte(got))
		})
	}
}
