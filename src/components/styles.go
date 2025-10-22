package components

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles holds all the styling definitions for the application
type Styles struct {
	Border         lipgloss.Style
	ActiveBorder   lipgloss.Style
	Title          lipgloss.Style
	PaneNumber     lipgloss.Style
	StatusGreen    lipgloss.Style
	StatusRed      lipgloss.Style
	StatusYellow   lipgloss.Style
	SelectedItem   lipgloss.Style
	Help           lipgloss.Style
	Key            lipgloss.Style
}

// NewStyles creates and returns a new Styles instance
func NewStyles() Styles {
	return Styles{
		Border: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1),

		ActiveBorder: lipgloss.NewStyle().
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("213")).
			Padding(0, 1).
			Bold(true),

		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("213")).
			Underline(true),

		PaneNumber: lipgloss.NewStyle().
			Foreground(lipgloss.Color("51")).
			Bold(true),

		StatusGreen: lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Background(lipgloss.Color("22")).
			Bold(true).
			Padding(0, 1),

		StatusRed: lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Background(lipgloss.Color("52")).
			Bold(true).
			Padding(0, 1),

		StatusYellow: lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Background(lipgloss.Color("58")).
			Bold(true).
			Padding(0, 1),

		SelectedItem: lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Bold(true),

		Help: lipgloss.NewStyle().
			Foreground(lipgloss.Color("51")).
			Background(lipgloss.Color("235")).
			Padding(0, 1).
			Bold(true),

		Key: lipgloss.NewStyle().
			Foreground(lipgloss.Color("213")).
			Bold(true),
	}
}
