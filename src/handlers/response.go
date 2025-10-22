package handlers

import (
	"fmt"

	"postty/src/types"
)

// HandleResponse handles HTTP response messages
func HandleResponse(m types.Model, msg types.ResponseMsg) types.Model {
	m.Executing = false
	if msg.Err != nil {
		m.ResponseViewport.SetContent(fmt.Sprintf("Error: %v", msg.Err))
		m.StatusCode = 0

		// Still add to history even if there was an error
		if m.PendingRequest != nil {
			m.PendingRequest.StatusCode = 0
			m.PendingRequest.ResponseBody = fmt.Sprintf("Error: %v", msg.Err)
			m = AddToHistory(m,
				m.PendingRequest.Method,
				m.PendingRequest.URL,
				m.PendingRequest.Body,
				m.PendingRequest.ContentType,
				m.PendingRequest.Headers,
				0,
				fmt.Sprintf("Error: %v", msg.Err),
			)
		}
	} else {
		m.ResponseViewport.SetContent(msg.Body)
		m.StatusCode = msg.StatusCode

		// Add successful request to history
		if m.PendingRequest != nil {
			m = AddToHistory(m,
				m.PendingRequest.Method,
				m.PendingRequest.URL,
				m.PendingRequest.Body,
				m.PendingRequest.ContentType,
				m.PendingRequest.Headers,
				msg.StatusCode,
				msg.Body,
			)
		}
	}

	// Clear pending request
	m.PendingRequest = nil

	return m
}
