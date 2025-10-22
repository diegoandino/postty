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
	} else {
		m.ResponseViewport.SetContent(msg.Body)
		m.StatusCode = msg.StatusCode
	}
	return m
}
