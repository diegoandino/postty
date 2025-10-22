package handlers

import (
	"postty/src/types"
)

// HandleContentTypeNavigation handles up/down navigation in content-type pane
func HandleContentTypeNavigation(m types.Model, direction string) types.Model {
	if direction == "up" {
		if m.SelectedHeader > 0 {
			m.SelectedHeader--
		}
	} else if direction == "down" {
		if m.SelectedHeader < len(types.ContentTypes)-1 {
			m.SelectedHeader++
		}
	}
	return m
}
