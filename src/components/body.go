package components

import (
	"postty/src/types"
)

// RenderBodyPane renders the request body input pane
func RenderBodyPane(m types.Model, styles Styles, width, height int) string {
	bodyTitle := styles.PaneNumber.Render("[3] ") + styles.Title.Render("Body")
	bodyContent := bodyTitle + "\n" + m.BodyInput.View()

	style := styles.Border
	if m.ActivePane == types.BodyPane {
		style = styles.ActiveBorder
	}

	return style.Width(width).Height(height).Render(bodyContent)
}
