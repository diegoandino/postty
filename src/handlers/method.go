package handlers

import (
	tea "github.com/charmbracelet/bubbletea"

	"postty/src/services"
	"postty/src/types"
)

// HandleMethodNavigation handles up/down navigation in method pane
func HandleMethodNavigation(m types.Model, direction string) types.Model {
	if direction == "up" {
		if m.SelectedMethod > 0 {
			m.SelectedMethod--
		}
	} else if direction == "down" {
		if m.SelectedMethod < len(types.HTTPMethods)-1 {
			m.SelectedMethod++
		}
	}
	return m
}

// HandleMethodExecute handles request execution from method pane
func HandleMethodExecute(m types.Model) (types.Model, tea.Cmd) {
	if m.URLInput.Value() != "" && !m.Executing {
		m.Executing = true
		m.ResponseViewport.SetContent("Executing request...")
		return m, services.ExecuteRequest(
			types.HTTPMethods[m.SelectedMethod],
			m.URLInput.Value(),
			m.BodyInput.Value(),
			types.ContentTypes[m.SelectedHeader],
			m.CustomHeaders,
		)
	}
	return m, nil
}
