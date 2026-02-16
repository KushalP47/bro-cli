package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	CommitAccept     = "accept"
	CommitEdit       = "edit"
	CommitRegenerate = "regenerate"
	CommitQuit       = "quit"
)

type commitModel struct {
	message string
	choice  string
	done    bool
}

func (m commitModel) Init() tea.Cmd {
	return nil
}

func (m commitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "y":
			m.choice = CommitAccept
			m.done = true
			return m, tea.Quit
		case "e":
			m.choice = CommitEdit
			m.done = true
			return m, tea.Quit
		case "r":
			m.choice = CommitRegenerate
			m.done = true
			return m, tea.Quit
		case "q", "esc", "ctrl+c":
			m.choice = CommitQuit
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m commitModel) View() string {
	if m.done {
		return ""
	}

	s := "\n"
	s += TitleStyle.Render("Generated commit message:") + "\n\n"
	s += "  " + m.message + "\n\n"
	s += HelpStyle.Render("[enter] accept  [e] edit  [r] regenerate  [q] quit") + "\n"
	return s
}

// RunCommitFlow displays the commit message and returns the user's choice.
func RunCommitFlow(message string) (string, error) {
	m := commitModel{message: message}
	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("commit flow: %w", err)
	}
	return result.(commitModel).choice, nil
}
