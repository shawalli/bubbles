package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shawalli/bubbles/tabs"
)

type PageModel struct {
	content string
}

func (pm PageModel) Init() tea.Cmd { return nil }

func (pm PageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return pm, nil }

func (pm PageModel) View() string { return pm.content }

type Model struct {
	tabs tabs.Model
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

	t, cmd := m.tabs.Update(msg)
	m.tabs = t.(tabs.Model)

	return m, cmd
}

func (m Model) View() string {
	return m.tabs.View()
}

func main() {
	t := []tabs.Tab{
		tabs.NewTab(
			"Tab 1",
			PageModel{content: "This is the content 1.\n\nYou can get back to me with crtl+1, shift+tab, or tabbing past tab 4!"},
		),
		tabs.NewTab(
			"Tab 2",
			PageModel{content: "This is the content 2.\n\nHappy to be here!"},
		),
		tabs.NewTab(
			"Tab 3",
			PageModel{content: "This is the content 3.\n\nHappy to be here!"},
		),
		tabs.NewTab(
			"Tab 4",
			PageModel{content: "This is the content 4.\n\nYou can get back to me with ctrl+4, tab, or shift+tabbing from tab 1!"},
		),
	}

	m := Model{tabs: tabs.New(t...).Wraparound()}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("could not run program: %v", err)
		os.Exit(1)
	}
}
