package handlers

import (
	tea "github.com/charmbracelet/bubbletea"

	"postty/src/types"
)

// HandleWindowSize handles window size change messages
func HandleWindowSize(m types.Model, msg tea.WindowSizeMsg) types.Model {
	m.Width = msg.Width
	m.Height = msg.Height

	rightColumnWidth := 35
	leftColumnWidth := m.Width - rightColumnWidth - 4

	if leftColumnWidth < 40 {
		leftColumnWidth = 40
		rightColumnWidth = m.Width - leftColumnWidth - 4
		if rightColumnWidth < 20 {
			rightColumnWidth = 20
		}
	}

	urlInputWidth := leftColumnWidth - 8
	if urlInputWidth < 20 {
		urlInputWidth = 20
	}
	m.URLInput.Width = urlInputWidth

	bodyInputWidth := leftColumnWidth - 8
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

	resultPaneWidth := leftColumnWidth + rightColumnWidth
	viewportWidth := resultPaneWidth - 4
	if viewportWidth < 20 {
		viewportWidth = 20
	}
	m.ResponseViewport.Width = viewportWidth
	m.ResponseViewport.Height = resultContentHeight - 1

	return m
}
