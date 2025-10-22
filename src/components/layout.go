package components

import (
	"github.com/charmbracelet/lipgloss"

	"postty/src/types"
)

// Dimensions holds the calculated dimensions for each pane
type Dimensions struct {
	HistoryColumnWidth int
	MiddleColumnWidth  int
	RightColumnWidth   int
	URLHeight          int
	MethodHeight       int
	BodyHeight         int
	HeaderHeight       int
	HeadersHeight      int
	HistoryHeight      int
	ResultHeight       int
}

// CalculateDimensions calculates all pane dimensions based on terminal size
func CalculateDimensions(width, height int) Dimensions {
	// Three-column layout: History | Middle | Right
	historyColumnWidth := 35
	rightColumnWidth := 40
	middleColumnWidth := width - historyColumnWidth - rightColumnWidth - 6

	// Ensure minimum widths
	if middleColumnWidth < 40 {
		middleColumnWidth = 40
		// Adjust other columns if needed
		remainingWidth := width - middleColumnWidth - 6
		historyColumnWidth = remainingWidth / 2
		rightColumnWidth = remainingWidth - historyColumnWidth

		if historyColumnWidth < 25 {
			historyColumnWidth = 25
		}
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

	// History takes full height of left column
	historyHeight := contentHeight

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

	// Result height without history
	resultHeight := (contentHeight - urlHeight - bodyHeight) - 10
	if resultHeight < 8 {
		resultHeight = 8
	}

	return Dimensions{
		HistoryColumnWidth: historyColumnWidth,
		MiddleColumnWidth:  middleColumnWidth,
		RightColumnWidth:   rightColumnWidth,
		URLHeight:          urlHeight,
		MethodHeight:       methodHeight,
		BodyHeight:         bodyHeight,
		HeaderHeight:       headerHeight,
		HeadersHeight:      headersHeight,
		HistoryHeight:      historyHeight,
		ResultHeight:       resultHeight,
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
	historyPane := RenderHistoryPane(m, styles, dims.HistoryColumnWidth, dims.HistoryHeight)

	urlPane := RenderURLPane(m, styles, dims.MiddleColumnWidth, dims.URLHeight)
	bodyPane := RenderBodyPane(m, styles, dims.MiddleColumnWidth, dims.BodyHeight)
	resultPane := RenderResponsePane(m, styles, dims.MiddleColumnWidth, dims.ResultHeight)

	methodPane := RenderMethodPane(m, styles, dims.RightColumnWidth, dims.MethodHeight)
	headerPane := RenderContentTypePane(m, styles, dims.RightColumnWidth, dims.HeaderHeight)
	headersPane := RenderCustomHeadersPane(m, styles, dims.RightColumnWidth, dims.HeadersHeight)

	// Compose layout: History | Middle | Right
	historyColumn := historyPane
	middleColumn := lipgloss.JoinVertical(lipgloss.Left, urlPane, bodyPane, resultPane)
	rightColumn := lipgloss.JoinVertical(lipgloss.Left, methodPane, headerPane, headersPane)
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, historyColumn, middleColumn, rightColumn)

	// Render help bar
	help := styles.Help.Render(
		styles.Key.Render("Tab") + " Next Pane │ " +
			styles.Key.Render("1-7") + " Jump │ " +
			styles.Key.Render("↑↓jk") + " Scroll │ " +
			styles.Key.Render("Enter") + "/" + styles.Key.Render("Alt+Enter") + " Send │ " +
			styles.Key.Render("esc") + "/" + styles.Key.Render("q") + " Quit",
	)

	return mainView + "\n" + help
}
