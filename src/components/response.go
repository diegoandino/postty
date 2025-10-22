package components

import (
	"fmt"

	"postty/src/types"
)

// RenderResponsePane renders the HTTP response pane
func RenderResponsePane(m types.Model, styles Styles, width, height int) string {
	resultTitle := styles.PaneNumber.Render("[5] ") + styles.Title.Render("Result")

	if m.StatusCode > 0 {
		statusStyle := styles.StatusGreen
		if m.StatusCode >= 400 {
			statusStyle = styles.StatusRed
		} else if m.StatusCode >= 300 {
			statusStyle = styles.StatusYellow
		}
		resultTitle += " " + statusStyle.Render(fmt.Sprintf("[%d]", m.StatusCode))
	}

	resultContent := resultTitle + "\n" + m.ResponseViewport.View()

	style := styles.Border
	if m.ActivePane == types.ResponsePane {
		style = styles.ActiveBorder
	}

	return style.Width(width).Height(height).Render(resultContent)
}
