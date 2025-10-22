package components

import (
	"postty/src/types"
)

// RenderContentTypePane renders the Content-Type selection pane
func RenderContentTypePane(m types.Model, styles Styles, width, height int) string {
	headerTitle := styles.PaneNumber.Render("[4] ") + styles.Title.Render("Content-Type")
	headerContent := headerTitle + "\n"

	for i, ct := range types.ContentTypes {
		if i == m.SelectedHeader {
			headerContent += styles.SelectedItem.Render("▶ "+ct) + "\n"
		} else {
			headerContent += "  " + ct + "\n"
		}
	}

	style := styles.Border
	if m.ActivePane == types.HeaderPane {
		style = styles.ActiveBorder
	}

	return style.Width(width).Height(height).Render(headerContent)
}
