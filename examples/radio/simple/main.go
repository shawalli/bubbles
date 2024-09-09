package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/shawalli/bubbles/radio"
)

type Model struct {
	vRadio radio.Model
	hRadio radio.Model
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	t, _ := m.hRadio.Update(msg)
	m.hRadio = t.(radio.Model)

	t, _ = m.vRadio.Update(msg)
	m.vRadio = t.(radio.Model)

	return m, nil
}

func (m Model) View() string {
	return gloss.JoinVertical(
		gloss.Top,
		m.vRadio.View(),
		m.hRadio.View(),
	)
}

func main() {
	vButtons := []tea.Model{
		radio.NewButton("10%"),
		radio.NewButton("15%"),
		radio.NewButton("20%"),
		radio.NewButton("22%"),
		radio.NewButton("25%"),
		radio.NewButton("None"),
	}

	hButtons := []tea.Model{
		radio.NewButton("Pay Cash"),
		radio.NewButton("Pay Credit"),
		radio.NewButton("PayBuddy"),
		radio.NewButton("BNPL"),
	}

	m := Model{
		vRadio: radio.New(true, vButtons...),
		hRadio: radio.New(false, hButtons...),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("could not run program: %v", err)
		os.Exit(1)
	}
}
