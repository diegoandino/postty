package handlers

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"postty/src/types"
)

// HandleCustomHeadersNavigation handles up/down navigation in headers pane
func HandleCustomHeadersNavigation(m types.Model, direction string) types.Model {
	if direction == "up" {
		if m.SelectedCustomHeader > 0 {
			m.SelectedCustomHeader--
		}
	} else if direction == "down" {
		if m.SelectedCustomHeader < len(m.CustomHeaders)-1 {
			m.SelectedCustomHeader++
		}
	}
	return m
}

// HandleCustomHeadersAdd initiates adding a new header
func HandleCustomHeadersAdd(m types.Model) types.Model {
	m.HeadersMode = types.HeadersAddMode
	m.SelectedTemplate = 0
	return m
}

// HandleCustomHeadersDelete deletes the selected header
func HandleCustomHeadersDelete(m types.Model) types.Model {
	if len(m.CustomHeaders) > 0 {
		m.CustomHeaders = append(m.CustomHeaders[:m.SelectedCustomHeader], m.CustomHeaders[m.SelectedCustomHeader+1:]...)
		if m.SelectedCustomHeader >= len(m.CustomHeaders) && len(m.CustomHeaders) > 0 {
			m.SelectedCustomHeader = len(m.CustomHeaders) - 1
		}
		if len(m.CustomHeaders) == 0 {
			m.SelectedCustomHeader = 0
		}
	}
	return m
}

// HandleCustomHeadersEdit initiates editing the selected header
func HandleCustomHeadersEdit(m types.Model) (types.Model, tea.Cmd) {
	if len(m.CustomHeaders) > 0 {
		m.HeadersMode = types.HeadersEditMode
		m.HeaderEditInput.SetValue(m.CustomHeaders[m.SelectedCustomHeader].Value)
		m.HeaderEditInput.Focus()
		return m, textinput.Blink
	}
	return m, nil
}

// HandleTemplateNavigation handles up/down navigation in template selection
func HandleTemplateNavigation(m types.Model, direction string) types.Model {
	if direction == "up" {
		if m.SelectedTemplate > 0 {
			m.SelectedTemplate--
		}
	} else if direction == "down" {
		if m.SelectedTemplate < len(types.HeaderTemplates)-1 {
			m.SelectedTemplate++
		}
	}
	return m
}

// HandleTemplateSelect selects a template and adds it as a header
func HandleTemplateSelect(m types.Model) (types.Model, tea.Cmd) {
	template := types.HeaderTemplates[m.SelectedTemplate]

	// Add the new header
	if template.Key == "" {
		// Custom header - user will define both key and value
		m.CustomHeaders = append(m.CustomHeaders, types.Header{Key: "Custom-Header", Value: ""})
	} else {
		// Template header - use template key and placeholder as initial value
		m.CustomHeaders = append(m.CustomHeaders, types.Header{
			Key:   template.Key,
			Value: template.Placeholder,
		})
	}

	// Always go to edit mode to let user edit the value immediately
	m.SelectedCustomHeader = len(m.CustomHeaders) - 1
	m.HeadersMode = types.HeadersEditMode
	m.HeaderEditInput.SetValue(m.CustomHeaders[m.SelectedCustomHeader].Value)
	m.HeaderEditInput.Focus()

	return m, textinput.Blink
}

// HandleHeaderEditSave saves the edited header value
func HandleHeaderEditSave(m types.Model) types.Model {
	if len(m.CustomHeaders) > 0 {
		m.CustomHeaders[m.SelectedCustomHeader].Value = m.HeaderEditInput.Value()
		m.HeadersMode = types.HeadersViewMode
		m.HeaderEditInput.Blur()
	}
	return m
}

// HandleHeaderEditCancel cancels header editing
func HandleHeaderEditCancel(m types.Model) types.Model {
	m.HeadersMode = types.HeadersViewMode
	m.HeaderEditInput.Blur()
	return m
}

// HandleAddModeCancel cancels header template selection
func HandleAddModeCancel(m types.Model) types.Model {
	m.HeadersMode = types.HeadersViewMode
	return m
}
