package handlers

import (
	"postty/src/types"
)

// HandleResponseScroll handles scrolling in the response pane
func HandleResponseScroll(m types.Model, action string) types.Model {
	switch action {
	case "up":
		m.ResponseViewport.ScrollUp(1)
	case "down":
		m.ResponseViewport.ScrollDown(1)
	case "pgup":
		m.ResponseViewport.HalfPageUp()
	case "pgdown":
		m.ResponseViewport.HalfPageDown()
	case "top":
		m.ResponseViewport.GotoTop()
	case "bottom":
		m.ResponseViewport.GotoBottom()
	}
	return m
}
