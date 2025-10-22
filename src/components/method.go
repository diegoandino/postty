package components

import (
	"postty/src/types"
)

// RenderMethodPane renders the HTTP method selection pane
func RenderMethodPane(m types.Model, styles Styles, width, height int) string {
	methodTitle := styles.PaneNumber.Render("[2] ") + styles.Title.Render("Method")
	methodContent := methodTitle + "\n"

	for i, method := range types.HTTPMethods {
		if i == m.SelectedMethod {
			methodContent += styles.SelectedItem.Render("▶ "+method) + "\n"
		} else {
			methodContent += "  " + method + "\n"
		}
	}

	style := styles.Border
	if m.ActivePane == types.MethodPane {
		style = styles.ActiveBorder
	}

	return style.Width(width).Height(height).Render(methodContent)
}
