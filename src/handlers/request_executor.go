package handlers

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"postty/src/services"
	"postty/src/types"
)

// ExecuteRequestWithHistory executes an HTTP request and stores it for history tracking
func ExecuteRequestWithHistory(m types.Model) (types.Model, tea.Cmd) {
	if m.URLInput.Value() == "" || m.Executing {
		return m, nil
	}

	// Get request details
	method := types.HTTPMethods[m.SelectedMethod]
	url := m.URLInput.Value()
	body := m.BodyInput.Value()
	contentType := types.ContentTypes[m.SelectedHeader]
	headers := make([]types.Header, len(m.CustomHeaders))
	copy(headers, m.CustomHeaders)

	// Store pending request for history
	m.PendingRequest = &types.HistoryItem{
		Method:      method,
		URL:         url,
		Body:        body,
		ContentType: contentType,
		Headers:     headers,
		Timestamp:   time.Now().Format("2006-01-02 15:04:05"),
	}

	// Mark as executing
	m.Executing = true
	m.ResponseViewport.SetContent("Executing request...")

	// Execute the request
	return m, services.ExecuteRequest(method, url, body, contentType, m.CustomHeaders)
}
