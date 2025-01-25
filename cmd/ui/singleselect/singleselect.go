package singleselect

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices   []string
	cursor    int
	selected  int
	selection *Selection
	confirmed bool
}

type Selection struct {
	Selected string
	Index    int
}

func InitialSelection() Selection {
	return Selection{Selected: "", Index: -1}
}

func (s *Selection) Reset() *Selection {
	s.Update("", -1)
	return s
}

func (s *Selection) Update(name string, index int) {
	s.Selected = name
	s.Index = index
}

func InitialModelSingleSelect(choices []string, selection *Selection) model {
	return model{
		choices:   choices,
		cursor:    0,
		selected:  selection.Index,
		selection: selection,
		confirmed: false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "y", "Y":
			m.confirmed = true
			return m, tea.Quit

		case "ctrl+c", "q":
			m.confirmed = false
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter", " ":
			if m.selected == m.cursor {
				m.selected = -1
				m.selection.Reset()
			} else {
				m.selected = m.cursor
				m.selection.Update(m.choices[m.selected], m.selected)
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.confirmed {
		if m.selected == -1 {
			s := "You need to select a resource directory to add translations.\n"
			return s
		}

		s := fmt.Sprintf("You selected to translate to languages in %s resources directory\n", m.choices[m.selected])
		return s
	}

	s := "Select one of resources directory found in the current project to add translations:\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if m.selected == i {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress q to cancel."
	s += "\nPress y to continue.\n"

	return s
}
