package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shawalli/bubbles/tabs"
)

func tabSizeCmd(width, height int) tea.Cmd {
	return func() tea.Msg { return tabs.TabSizeMsg{Width: width, Height: height} }
}

type PageModel struct {
	content string

	width  int
	height int
}

func (pm PageModel) Init() tea.Cmd { return nil }

func (pm PageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left":
			pm.width = pm.width - 1
			cmd = tabSizeCmd(pm.width, pm.height)
		case "right":
			pm.width = pm.width + 1
			cmd = tabSizeCmd(pm.width, pm.height)
		}
	case tea.WindowSizeMsg:
		pm.width = msg.Width
		pm.height = msg.Height
	}
	return pm, cmd

}

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
			PageModel{content: "This is the content 1.\n\nYou can adjust my width with the left and right arrows!"},
		),
		tabs.NewTab(
			"Tab 2",
			PageModel{content: "This is the content 2.\n\nWhen you're looking at me,\nany change in my width affects my sibling tabs as well!"},
		),
		tabs.NewTab(
			"Tab 3",
			PageModel{content: "This is the content 3.\n\nResizing the terminal window resets all tabs to the terminal width."},
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
