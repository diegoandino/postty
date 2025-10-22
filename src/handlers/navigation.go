package handlers

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"postty/src/types"
)

// HandleTab handles Tab key navigation
func HandleTab(m types.Model) (types.Model, tea.Cmd) {
	m.ActivePane++
	if m.ActivePane > types.HistoryPane {
		m.ActivePane = types.URLPane
	}

	m.URLInput.Blur()
	m.BodyInput.Blur()
	m.HeaderEditInput.Blur()
	m.HeadersMode = types.HeadersViewMode

	if m.ActivePane == types.URLPane {
		m.URLInput.Focus()
		return m, textinput.Blink
	} else if m.ActivePane == types.BodyPane {
		m.BodyInput.Focus()
		return m, textarea.Blink
	}
	return m, nil
}

// HandleShiftTab handles Shift+Tab key navigation
func HandleShiftTab(m types.Model) (types.Model, tea.Cmd) {
	m.ActivePane--
	if m.ActivePane < types.URLPane {
		m.ActivePane = types.HistoryPane
	}

	m.URLInput.Blur()
	m.BodyInput.Blur()
	m.HeaderEditInput.Blur()
	m.HeadersMode = types.HeadersViewMode

	if m.ActivePane == types.URLPane {
		m.URLInput.Focus()
		return m, textinput.Blink
	} else if m.ActivePane == types.BodyPane {
		m.BodyInput.Focus()
		return m, textarea.Blink
	}
	return m, nil
}

// HandleJumpToPane handles jumping to a specific pane
func HandleJumpToPane(m types.Model, pane types.Pane) (types.Model, tea.Cmd) {
	m.ActivePane = pane
	m.URLInput.Blur()
	m.BodyInput.Blur()
	m.HeaderEditInput.Blur()
	m.HeadersMode = types.HeadersViewMode

	if m.ActivePane == types.URLPane {
		m.URLInput.Focus()
		return m, textinput.Blink
	} else if m.ActivePane == types.BodyPane {
		m.BodyInput.Focus()
		return m, textarea.Blink
	}
	return m, nil
}
