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

	// Calculate total available height for panes
	// Layout format: "\n" + mainView + "\n" + help
	// Account for:
	// - Top padding (\n before mainView): 1 line
	// - Space before help (\n after mainView): 1 line
	// - Help bar: 1 line
	// Total overhead: 3 lines
	totalAvailable := height - 3
	if totalAvailable < 20 {
		totalAvailable = 20 // Minimum height to fit all content
	}

	// ALL THREE COLUMNS MUST BE THIS HEIGHT
	columnHeight := totalAvailable
	urlHeight := 4

	// Body: flexible but reasonable minimum
	bodyHeight := columnHeight / 4
	if bodyHeight < 6 {
		bodyHeight = 6
	}

	// Result: remainder (ensures middle column = columnHeight exactly)
	resultHeight := columnHeight - urlHeight - bodyHeight

	// Right column breakdown
	methodHeight := (columnHeight * 40) / 100
	if methodHeight < 10 {
		methodHeight = 10
	}

	// Content-Type: needs to fit 5 types + title + borders = ~8 lines minimum
	headerHeight := (columnHeight * 30) / 100
	if headerHeight < 8 {
		headerHeight = 8
	}

	// Headers: remainder (ensures right column = columnHeight exactly)
	headersHeight := columnHeight - methodHeight - headerHeight
	if headersHeight < 5 {
		// If not enough space, shrink the other panes proportionally
		headersHeight = 5
		methodHeight = (columnHeight - 5 - headerHeight)
		if methodHeight < 10 {
			methodHeight = 10
			headerHeight = columnHeight - 15
		}
	}

	// History: matches column height exactly
	historyHeight := columnHeight

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

	// Add top padding and combine with help bar
	return "\n" + mainView + "\n" + help
}
