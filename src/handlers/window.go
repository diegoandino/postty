package handlers

import (
	tea "github.com/charmbracelet/bubbletea"

	"postty/src/components"
	"postty/src/types"
)

// HandleWindowSize handles window size change messages
func HandleWindowSize(m types.Model, msg tea.WindowSizeMsg) types.Model {
	m.Width = msg.Width
	m.Height = msg.Height

	// Use centralized dimension calculations
	dims := components.CalculateDimensions(m.Width, m.Height)

	// Update URL input width (account for borders + padding)
	urlInputWidth := dims.MiddleColumnWidth - 6
	if urlInputWidth < 20 {
		urlInputWidth = 20
	}
	m.URLInput.Width = urlInputWidth

	// Update Body input width and height
	bodyInputWidth := dims.MiddleColumnWidth - 6
	if bodyInputWidth < 20 {
		bodyInputWidth = 20
	}
	m.BodyInput.SetWidth(bodyInputWidth)

	// Body height: pane height minus border (2) and title (1) and padding (1)
	bodyContentHeight := dims.BodyHeight - 4
	if bodyContentHeight < 3 {
		bodyContentHeight = 3
	}
	m.BodyInput.SetHeight(bodyContentHeight)

	// Update Response viewport
	viewportWidth := dims.MiddleColumnWidth - 4
	if viewportWidth < 20 {
		viewportWidth = 20
	}
	m.ResponseViewport.Width = viewportWidth

	// Response height: pane height minus border (2) and title line (1) and padding (1)
	responseViewportHeight := dims.ResultHeight - 4
	if responseViewportHeight < 5 {
		responseViewportHeight = 5
	}
	m.ResponseViewport.Height = responseViewportHeight

	// Update History viewport
	historyViewportWidth := dims.HistoryColumnWidth - 4
	if historyViewportWidth < 20 {
		historyViewportWidth = 20
	}
	m.HistoryViewport.Width = historyViewportWidth

	// History height: pane height minus border (2) and title (1) and help text (1) and padding (2)
	historyViewportHeight := dims.HistoryHeight - 6
	if historyViewportHeight < 5 {
		historyViewportHeight = 5
	}
	m.HistoryViewport.Height = historyViewportHeight

	// Update Method viewport
	methodViewportWidth := dims.RightColumnWidth - 4
	if methodViewportWidth < 15 {
		methodViewportWidth = 15
	}
	m.MethodViewport.Width = methodViewportWidth

	// Method height: pane height minus border (2) and title (1) and padding (1)
	methodViewportHeight := dims.MethodHeight - 4
	if methodViewportHeight < 3 {
		methodViewportHeight = 3
	}
	m.MethodViewport.Height = methodViewportHeight

	// Update ContentType viewport
	contentTypeViewportWidth := dims.RightColumnWidth - 4
	if contentTypeViewportWidth < 15 {
		contentTypeViewportWidth = 15
	}
	m.ContentTypeViewport.Width = contentTypeViewportWidth

	// ContentType height: pane height minus border (2) and title (1) and padding (1)
	contentTypeViewportHeight := dims.HeaderHeight - 4
	if contentTypeViewportHeight < 3 {
		contentTypeViewportHeight = 3
	}
	m.ContentTypeViewport.Height = contentTypeViewportHeight

	return m
}
