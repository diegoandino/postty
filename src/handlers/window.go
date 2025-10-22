package handlers

import (
	tea "github.com/charmbracelet/bubbletea"

	"postty/src/types"
)

// HandleWindowSize handles window size change messages
func HandleWindowSize(m types.Model, msg tea.WindowSizeMsg) types.Model {
	m.Width = msg.Width
	m.Height = msg.Height

	// Three-column layout: History | Middle | Right
	historyColumnWidth := 35
	rightColumnWidth := 40
	middleColumnWidth := m.Width - historyColumnWidth - rightColumnWidth - 6

	// Ensure minimum widths
	if middleColumnWidth < 40 {
		middleColumnWidth = 40
		remainingWidth := m.Width - middleColumnWidth - 6
		historyColumnWidth = remainingWidth / 2
		rightColumnWidth = remainingWidth - historyColumnWidth

		if historyColumnWidth < 25 {
			historyColumnWidth = 25
		}
		if rightColumnWidth < 20 {
			rightColumnWidth = 20
		}
	}

	// URL and Body inputs are in the middle column
	urlInputWidth := middleColumnWidth - 8
	if urlInputWidth < 20 {
		urlInputWidth = 20
	}
	m.URLInput.Width = urlInputWidth

	bodyInputWidth := middleColumnWidth - 8
	if bodyInputWidth < 20 {
		bodyInputWidth = 20
	}
	m.BodyInput.SetWidth(bodyInputWidth)

	availableHeight := m.Height - 2
	borderOverhead := 6
	contentHeight := availableHeight - borderOverhead

	if contentHeight < 24 {
		contentHeight = 24
	}

	urlContentHeight := 3
	topSectionContentHeight := contentHeight - int(float64(contentHeight)*0.4)
	if topSectionContentHeight < 14 {
		topSectionContentHeight = 14
	}

	resultContentHeight := contentHeight - topSectionContentHeight
	if resultContentHeight < 6 {
		resultContentHeight = 6
	}

	methodContentHeight := int(float64(topSectionContentHeight) * 0.55)
	if methodContentHeight < 8 {
		methodContentHeight = 8
	}
	headerContentHeight := (topSectionContentHeight - methodContentHeight) - 2
	if headerContentHeight < 5 {
		headerContentHeight = 5
	}

	rightColumnContentHeight := methodContentHeight + headerContentHeight
	bodyContentHeight := rightColumnContentHeight - urlContentHeight
	if bodyContentHeight < 5 {
		bodyContentHeight = 5
	}
	m.BodyInput.SetHeight(bodyContentHeight)

	// Response viewport is in the middle column
	viewportWidth := middleColumnWidth - 4
	if viewportWidth < 20 {
		viewportWidth = 20
	}
	m.ResponseViewport.Width = viewportWidth
	m.ResponseViewport.Height = resultContentHeight - 1

	// History viewport takes full height of left column
	historyViewportWidth := historyColumnWidth - 4
	if historyViewportWidth < 20 {
		historyViewportWidth = 20
	}
	m.HistoryViewport.Width = historyViewportWidth
	m.HistoryViewport.Height = contentHeight - 4 // Account for title and help text

	return m
}
