package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shawalli/bubbles/tabs"
)

type PageModel struct {
	content string

	width  int
	height int
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
			PageModel{content: "This is the content 1.\n\nYou can get to the next tab by typing the tab key."},
		),
		tabs.NewTab(
			"Tab 2",
			PageModel{content: "This is the content 2.\n\nYou can get back to me with shift+tab."},
		),
		tabs.NewTab(
			"Tab 3",
			PageModel{content: "This is the content 3.\n\nYou can get to me by typing the number 3 (or tabbing)."},
		),
		tabs.NewTab(
			"Tab 4",
			PageModel{content: "This is the content 4.\n\nHappy to be here!"},
		),
	}

	m := Model{tabs: tabs.New(t...)}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("could not run program: %v", err)
		os.Exit(1)
	}
}
