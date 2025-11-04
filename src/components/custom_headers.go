package components

import (
	"fmt"

	"postty/src/types"
)

// RenderCustomHeadersPane renders the custom headers management pane
func RenderCustomHeadersPane(m types.Model, styles Styles, width, height int) string {
	headersTitle := styles.PaneNumber.Render("[6] ") + styles.Title.Render("Custom Headers")
	headersContent := headersTitle + "\n"

	switch m.HeadersMode {
	case types.HeadersViewMode:
		if len(m.CustomHeaders) == 0 {
			headersContent += "  (no headers)\n"
			headersContent += "  Press 'a' to add"
		} else {
			for i, h := range m.CustomHeaders {
				prefix := "  "
				if i == m.SelectedCustomHeader {
					prefix = styles.SelectedItem.Render("▶ ")
				}
				headerLine := fmt.Sprintf("%s: %s", h.Key, h.Value)
				if h.Key == "" && h.Value == "" {
					headerLine = "(empty)"
				}
				headersContent += prefix + headerLine + "\n"
			}
			headersContent += "  a: add | d: del | e: edit"
		}

	case types.HeadersAddMode:
		headersContent += "  Select header type:\n"
		headersContent += "\n"
		for i, template := range types.HeaderTemplates {
			prefix := "  "
			if i == m.SelectedTemplate {
				prefix = styles.SelectedItem.Render("▶ ")
			}
			headersContent += prefix + template.Name + "\n"
		}
		headersContent += "\n"
		headersContent += "  Enter: select | Esc: cancel\n"

	case types.HeadersEditMode:
		if len(m.CustomHeaders) > 0 {
			header := m.CustomHeaders[m.SelectedCustomHeader]
			headersContent += fmt.Sprintf("  Editing: %s\n", header.Key)
			headersContent += "\n"
			headersContent += "  Value:\n"
			headersContent += "  " + m.HeaderEditInput.View() + "\n"
			headersContent += "\n"
			headersContent += "  Enter: save | Esc: cancel\n"
		}
	}

	style := styles.Border
	if m.ActivePane == types.HeadersPane {
		style = styles.ActiveBorder
	}

	// Subtract 2 for borders (top + bottom)
	return style.Width(width).Height(height - 2).Render(headersContent)
}
