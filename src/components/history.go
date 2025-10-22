package components

import (
	"fmt"
	"strings"

	"postty/src/types"
)

// wrapText wraps text to fit within the given width
func wrapText(text string, width int) []string {
	if len(text) <= width {
		return []string{text}
	}

	var lines []string
	for len(text) > 0 {
		if len(text) <= width {
			lines = append(lines, text)
			break
		}
		lines = append(lines, text[:width])
		text = text[width:]
	}
	return lines
}

// RenderHistoryPane renders the request history pane
func RenderHistoryPane(m types.Model, styles Styles, width, height int) string {
	// Title with count
	countText := ""
	if len(m.History) > 0 {
		countText = fmt.Sprintf(" (%d)", len(m.History))
	}
	historyTitle := styles.PaneNumber.Render("[7] ") + styles.Title.Render("History"+countText)
	historyContent := historyTitle + "\n\n"

	if len(m.History) == 0 {
		historyContent += "  No history yet.\n\n"
		historyContent += "  Make a request to\n"
		historyContent += "  see it here!\n"
	} else {
		// Build history list
		var historyLines []string
		for i, item := range m.History {
			// Add separator between items for clarity
			if i > 0 {
				historyLines = append(historyLines, "")
			}

			// Request number (1-indexed)
			requestNum := fmt.Sprintf("[%d]", i+1)

			// Format: METHOD [STATUS]
			statusText := ""
			if item.StatusCode > 0 {
				statusStyle := styles.StatusGreen
				if item.StatusCode >= 400 {
					statusStyle = styles.StatusRed
				} else if item.StatusCode >= 300 {
					statusStyle = styles.StatusYellow
				}
				statusText = " " + statusStyle.Render(fmt.Sprintf("[%d]", item.StatusCode))
			}

			methodLine := requestNum + " " + item.Method + statusText

			// Wrap URL to show more context (show up to 100 chars, wrapped to fit width)
			displayURL := item.URL
			maxURLChars := 100 // Show up to 100 characters
			if len(displayURL) > maxURLChars {
				displayURL = displayURL[:maxURLChars] + "..."
			}

			// Wrap URL to fit viewport width (account for indentation)
			urlWidth := width - 6 // Account for borders and indentation
			wrappedURLLines := wrapText(displayURL, urlWidth)

			// Format timestamp (just time, not date)
			timePart := ""
			if len(item.Timestamp) >= 16 {
				timePart = item.Timestamp[11:16] // HH:MM
			}

			if i == m.SelectedHistory {
				historyLines = append(historyLines, styles.SelectedItem.Render("▶ "+methodLine))
				for _, urlLine := range wrappedURLLines {
					historyLines = append(historyLines, styles.SelectedItem.Render("  "+urlLine))
				}
				if timePart != "" {
					historyLines = append(historyLines, styles.SelectedItem.Render("  "+timePart))
				}
			} else {
				historyLines = append(historyLines, "  "+methodLine)
				for _, urlLine := range wrappedURLLines {
					historyLines = append(historyLines, "  "+urlLine)
				}
				if timePart != "" {
					historyLines = append(historyLines, "  "+timePart)
				}
			}
		}

		// Set viewport content
		m.HistoryViewport.SetContent(strings.Join(historyLines, "\n"))

		// Calculate and set scroll position to keep selected item visible
		// Estimate 5 lines per item on average (separator + method + url lines + time)
		estimatedLinePosition := m.SelectedHistory * 5
		viewportHeight := m.HistoryViewport.Height

		// Keep selected item in view
		if estimatedLinePosition < m.HistoryViewport.YOffset {
			// Selected item is above visible area, scroll up
			m.HistoryViewport.SetYOffset(estimatedLinePosition)
		} else if estimatedLinePosition >= m.HistoryViewport.YOffset + viewportHeight - 5 {
			// Selected item is below visible area, scroll down
			m.HistoryViewport.SetYOffset(estimatedLinePosition - viewportHeight + 5)
		}

		historyContent += m.HistoryViewport.View()
		historyContent += "\n"
		historyContent += "  Enter: load | d: del\n"
	}

	style := styles.Border
	if m.ActivePane == types.HistoryPane {
		style = styles.ActiveBorder
	}

	return style.Width(width).Height(height).Render(historyContent)
}
