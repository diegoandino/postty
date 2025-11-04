package components

import (
	"strings"

	"postty/src/types"
)

// RenderContentTypePane renders the Content-Type selection pane
func RenderContentTypePane(m types.Model, styles Styles, width, height int) string {
	headerTitle := styles.PaneNumber.Render("[4] ") + styles.Title.Render("Content-Type")

	// Build content type list for viewport
	var ctLines []string
	for i, ct := range types.ContentTypes {
		if i == m.SelectedHeader {
			ctLines = append(ctLines, styles.SelectedItem.Render("▶ "+ct))
		} else {
			ctLines = append(ctLines, "  "+ct)
		}
	}

	// Set viewport content
	m.ContentTypeViewport.SetContent(strings.Join(ctLines, "\n"))

	// Auto-scroll to keep selected item visible
	if m.SelectedHeader < m.ContentTypeViewport.YOffset {
		m.ContentTypeViewport.SetYOffset(m.SelectedHeader)
	} else if m.SelectedHeader >= m.ContentTypeViewport.YOffset+m.ContentTypeViewport.Height {
		m.ContentTypeViewport.SetYOffset(m.SelectedHeader - m.ContentTypeViewport.Height + 1)
	}

	headerContent := headerTitle + "\n" + m.ContentTypeViewport.View()

	style := styles.Border
	if m.ActivePane == types.HeaderPane {
		style = styles.ActiveBorder
	}

	// Subtract 2 for borders (top + bottom)
	return style.Width(width).Height(height - 2).Render(headerContent)
}
