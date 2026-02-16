package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// PickerItem represents an item in the fuzzy picker.
type PickerItem struct {
	Name        string
	Description string
	Value       string
}

type pickerModel struct {
	items    []PickerItem
	filtered []int
	cursor   int
	filter   string
	selected *PickerItem
	done     bool
}

func newPickerModel(items []PickerItem) pickerModel {
	filtered := make([]int, len(items))
	for i := range items {
		filtered[i] = i
	}
	return pickerModel{items: items, filtered: filtered}
}

func (m pickerModel) Init() tea.Cmd {
	return nil
}

func (m pickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			if m.filter == "" || msg.String() != "q" {
				m.done = true
				return m, tea.Quit
			}
			// If filter is active, 'q' types into filter
			m.filter += "q"
			m.applyFilter()
			return m, nil
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		case "down", "j":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
			return m, nil
		case "enter":
			if len(m.filtered) > 0 {
				idx := m.filtered[m.cursor]
				m.selected = &m.items[idx]
				m.done = true
				return m, tea.Quit
			}
			return m, nil
		case "backspace":
			if len(m.filter) > 0 {
				m.filter = m.filter[:len(m.filter)-1]
				m.applyFilter()
			}
			return m, nil
		default:
			if len(msg.String()) == 1 {
				m.filter += msg.String()
				m.applyFilter()
			}
			return m, nil
		}
	}
	return m, nil
}

func (m *pickerModel) applyFilter() {
	if m.filter == "" {
		m.filtered = make([]int, len(m.items))
		for i := range m.items {
			m.filtered[i] = i
		}
	} else {
		m.filtered = nil
		lower := strings.ToLower(m.filter)
		for i, item := range m.items {
			if strings.Contains(strings.ToLower(item.Name), lower) ||
				strings.Contains(strings.ToLower(item.Description), lower) {
				m.filtered = append(m.filtered, i)
			}
		}
	}
	if m.cursor >= len(m.filtered) {
		m.cursor = max(0, len(m.filtered)-1)
	}
}

func (m pickerModel) View() string {
	if m.done {
		return ""
	}

	s := "\n"
	s += TitleStyle.Render("Select a command:") + "\n\n"

	if m.filter != "" {
		s += "  Filter: " + m.filter + "\n\n"
	}

	maxVisible := 15
	start := 0
	if m.cursor >= maxVisible {
		start = m.cursor - maxVisible + 1
	}

	for i := start; i < len(m.filtered) && i < start+maxVisible; i++ {
		item := m.items[m.filtered[i]]
		cursor := "  "
		if i == m.cursor {
			cursor = SelectedStyle.Render("> ")
			s += cursor + SelectedStyle.Render(item.Name)
			if item.Description != "" {
				s += " " + DimStyle.Render(item.Description)
			}
		} else {
			s += cursor + NormalStyle.Render(item.Name)
			if item.Description != "" {
				s += " " + DimStyle.Render(item.Description)
			}
		}
		s += "\n"
	}

	s += "\n" + HelpStyle.Render("↑/↓ navigate  type to filter  enter select  esc quit") + "\n"
	return s
}

// RunPicker displays the fuzzy picker and returns the selected item.
func RunPicker(items []PickerItem) (*PickerItem, error) {
	m := newPickerModel(items)
	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("picker: %w", err)
	}
	return result.(pickerModel).selected, nil
}
