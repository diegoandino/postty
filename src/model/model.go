package model

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"postty/src/types"
)

// New creates and initializes a new model
func New() types.Model {
	ti := textinput.New()
	ti.Placeholder = "https://api.example.com/endpoint"
	ti.Focus()
	ti.CharLimit = 500
	ti.Width = 40

	ta := textarea.New()
	ta.Placeholder = "Request body (JSON, XML, etc.)"
	ta.SetWidth(40)
	ta.SetHeight(8)

	vp := viewport.New(40, 10)
	vp.SetContent("")

	hei := textinput.New()
	hei.Placeholder = "Enter header value"
	hei.CharLimit = 500
	hei.Width = 30

	defaultHeaders := []types.Header{}

	return types.Model{
		ActivePane:           types.URLPane,
		SelectedMethod:       0,
		SelectedHeader:       0,
		URLInput:             ti,
		BodyInput:            ta,
		ResponseViewport:     vp,
		StatusCode:           0,
		Width:                0,
		Height:               0,
		Executing:            false,
		CustomHeaders:        defaultHeaders,
		SelectedCustomHeader: 0,
		HeadersMode:          types.HeadersViewMode,
		SelectedTemplate:     0,
		HeaderEditInput:      hei,
	}
}

// Init returns the initial command for the application
func Init() tea.Cmd {
	return textinput.Blink
}
