package components

import (
	"postty/src/types"
)

// RenderURLPane renders the URL input pane
func RenderURLPane(m types.Model, styles Styles, width, height int) string {
	urlTitle := styles.PaneNumber.Render("[1] ") + styles.Title.Render("URL")
	urlContent := urlTitle + "\n" + m.URLInput.View()

	style := styles.Border
	if m.ActivePane == types.URLPane {
		style = styles.ActiveBorder
	}

	// Subtract 2 for borders (top + bottom)
	return style.Width(width).Height(height - 2).Render(urlContent)
}
