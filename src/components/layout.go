package components

import (
	"github.com/charmbracelet/lipgloss"

	"postty/src/types"
)

// Dimensions holds the calculated dimensions for each pane
type Dimensions struct {
	LeftColumnWidth  int
	RightColumnWidth int
	URLHeight        int
	MethodHeight     int
	BodyHeight       int
	HeaderHeight     int
	HeadersHeight    int
	ResultHeight     int
}

// CalculateDimensions calculates all pane dimensions based on terminal size
func CalculateDimensions(width, height int) Dimensions {
	rightColumnWidth := 40
	leftColumnWidth := width - rightColumnWidth - 4

	if leftColumnWidth < 40 {
		leftColumnWidth = 40
		rightColumnWidth = width - leftColumnWidth - 4
		if rightColumnWidth < 20 {
			rightColumnWidth = 20
		}
	}

	availableHeight := height - 2
	borderOverhead := 6
	contentHeight := availableHeight - borderOverhead

	if contentHeight < 24 {
		contentHeight = 24
	}

	urlHeight := 3
	sectionHeight := contentHeight / 3
	if sectionHeight < 8 {
		sectionHeight = 8
	}

	methodHeight := sectionHeight
	if methodHeight < 8 {
		methodHeight = 8
	}

	headerHeight := sectionHeight - 3
	if headerHeight < 5 {
		headerHeight = 5
	}

	headersHeight := contentHeight - methodHeight - headerHeight
	if headersHeight < 6 {
		headersHeight = 6
	}

	bodyHeight := methodHeight - urlHeight
	if bodyHeight < 5 {
		bodyHeight = 5
	}

	resultHeight := (contentHeight - urlHeight - bodyHeight) - 10
	if resultHeight < 8 {
		resultHeight = 8
	}

	return Dimensions{
		LeftColumnWidth:  leftColumnWidth,
		RightColumnWidth: rightColumnWidth,
		URLHeight:        urlHeight,
		MethodHeight:     methodHeight,
		BodyHeight:       bodyHeight,
		HeaderHeight:     headerHeight,
		HeadersHeight:    headersHeight,
		ResultHeight:     resultHeight,
	}
}

// RenderLayout renders the complete application layout
func RenderLayout(m types.Model) string {
	if m.Width == 0 {
		return "Loading..."
	}

	styles := NewStyles()
	dims := CalculateDimensions(m.Width, m.Height)

	// Render all panes
	urlPane := RenderURLPane(m, styles, dims.LeftColumnWidth, dims.URLHeight)
	methodPane := RenderMethodPane(m, styles, dims.RightColumnWidth, dims.MethodHeight)
	bodyPane := RenderBodyPane(m, styles, dims.LeftColumnWidth, dims.BodyHeight)
	headerPane := RenderContentTypePane(m, styles, dims.RightColumnWidth, dims.HeaderHeight)
	resultPane := RenderResponsePane(m, styles, dims.LeftColumnWidth, dims.ResultHeight)
	headersPane := RenderCustomHeadersPane(m, styles, dims.RightColumnWidth, dims.HeadersHeight)

	// Compose layout
	leftColumn := lipgloss.JoinVertical(lipgloss.Left, urlPane, bodyPane, resultPane)
	rightColumn := lipgloss.JoinVertical(lipgloss.Left, methodPane, headerPane, headersPane)
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, rightColumn)

	// Render help bar
	help := styles.Help.Render(
		styles.Key.Render("Tab") + " Next Pane │ " +
			styles.Key.Render("1-6") + " Jump │ " +
			styles.Key.Render("↑↓jk") + " Scroll │ " +
			styles.Key.Render("Enter") + "/" + styles.Key.Render("Alt+Enter") + " Send │ " +
			styles.Key.Render("esc") + "/" + styles.Key.Render("q") + " Quit",
	)

	return mainView + "\n" + help
}
