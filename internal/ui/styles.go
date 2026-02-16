package ui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))
	SelectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	NormalStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	DimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	SuccessStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	ErrorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	HelpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)
