package components

import (
	"strings"

	"postty/src/types"
)

// RenderMethodPane renders the HTTP method selection pane
func RenderMethodPane(m types.Model, styles Styles, width, height int) string {
	methodTitle := styles.PaneNumber.Render("[2] ") + styles.Title.Render("Method")

	// Build method list content for viewport
	var methodLines []string
	for i, method := range types.HTTPMethods {
		if i == m.SelectedMethod {
			methodLines = append(methodLines, styles.SelectedItem.Render("▶ "+method))
		} else {
			methodLines = append(methodLines, "  "+method)
		}
	}

	// Set viewport content
	m.MethodViewport.SetContent(strings.Join(methodLines, "\n"))

	// Auto-scroll to keep selected item visible
	if m.SelectedMethod < m.MethodViewport.YOffset {
		m.MethodViewport.SetYOffset(m.SelectedMethod)
	} else if m.SelectedMethod >= m.MethodViewport.YOffset+m.MethodViewport.Height {
		m.MethodViewport.SetYOffset(m.SelectedMethod - m.MethodViewport.Height + 1)
	}

	methodContent := methodTitle + "\n" + m.MethodViewport.View()

	style := styles.Border
	if m.ActivePane == types.MethodPane {
		style = styles.ActiveBorder
	}

	// Subtract 2 for borders (top + bottom)
	return style.Width(width).Height(height - 2).Render(methodContent)
}
