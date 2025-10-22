package handlers

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"postty/src/types"
)

// HandleHistoryNavigation handles up/down navigation in history pane
func HandleHistoryNavigation(m types.Model, direction string) types.Model {
	if direction == "up" {
		if m.SelectedHistory > 0 {
			m.SelectedHistory--
		}
	} else if direction == "down" {
		if m.SelectedHistory < len(m.History)-1 {
			m.SelectedHistory++
		}
	}
	// Scrolling is automatically handled in RenderHistoryPane
	return m
}

// HandleHistoryLoad loads a history item into the form
func HandleHistoryLoad(m types.Model) (types.Model, tea.Cmd) {
	if len(m.History) == 0 {
		return m, nil
	}

	item := m.History[m.SelectedHistory]

	// Set URL
	m.URLInput.SetValue(item.URL)

	// Set method
	for i, method := range types.HTTPMethods {
		if method == item.Method {
			m.SelectedMethod = i
			break
		}
	}

	// Set body
	m.BodyInput.SetValue(item.Body)

	// Set content type
	for i, ct := range types.ContentTypes {
		if ct == item.ContentType {
			m.SelectedHeader = i
			break
		}
	}

	// Set custom headers
	m.CustomHeaders = make([]types.Header, len(item.Headers))
	copy(m.CustomHeaders, item.Headers)

	// Set response
	if item.ResponseBody != "" {
		m.ResponseViewport.SetContent(item.ResponseBody)
		m.StatusCode = item.StatusCode
	}

	// Switch to URL pane
	m.ActivePane = types.URLPane
	m.URLInput.Focus()

	return m, nil
}

// HandleHistoryDelete deletes the selected history item
func HandleHistoryDelete(m types.Model) types.Model {
	if len(m.History) == 0 {
		return m
	}

	// Remove the selected item
	m.History = append(m.History[:m.SelectedHistory], m.History[m.SelectedHistory+1:]...)

	// Adjust selection
	if m.SelectedHistory >= len(m.History) && len(m.History) > 0 {
		m.SelectedHistory = len(m.History) - 1
	}
	if len(m.History) == 0 {
		m.SelectedHistory = 0
	}

	return m
}

// AddToHistory adds a request to the history
func AddToHistory(m types.Model, method, url, body, contentType string, headers []types.Header, statusCode int, responseBody string) types.Model {
	// Create timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// Copy headers
	headersCopy := make([]types.Header, len(headers))
	copy(headersCopy, headers)

	// Create history item
	item := types.HistoryItem{
		Method:       method,
		URL:          url,
		Body:         body,
		ContentType:  contentType,
		Headers:      headersCopy,
		StatusCode:   statusCode,
		Timestamp:    timestamp,
		ResponseBody: responseBody,
	}

	// Add to beginning of history (most recent first)
	m.History = append([]types.HistoryItem{item}, m.History...)

	// Limit history to 50 items
	if len(m.History) > 50 {
		m.History = m.History[:50]
	}

	// Reset selection to the newest item
	m.SelectedHistory = 0

	return m
}

// HandleHistoryScroll handles scrolling in the history viewport
func HandleHistoryScroll(m types.Model, action string) types.Model {
	switch action {
	case "up":
		m.HistoryViewport.LineUp(1)
	case "down":
		m.HistoryViewport.LineDown(1)
	case "pgup":
		m.HistoryViewport.HalfPageUp()
	case "pgdown":
		m.HistoryViewport.HalfPageDown()
	}
	return m
}

// GetHistorySummary returns a summary string for the history pane title
func GetHistorySummary(m types.Model) string {
	if len(m.History) == 0 {
		return ""
	}
	return fmt.Sprintf(" (%d requests)", len(m.History))
}
