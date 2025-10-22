package types

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
)

// Pane represents different UI panes in the application
type Pane int

const (
	URLPane Pane = iota + 1
	MethodPane
	BodyPane
	HeaderPane
	ResponsePane
	HeadersPane
)

// HeadersMode represents the current mode of the headers pane
type HeadersMode int

const (
	HeadersViewMode HeadersMode = iota
	HeadersAddMode
	HeadersEditMode
)

// Header represents a custom HTTP header
type Header struct {
	Key   string
	Value string
}

// HeaderTemplate represents a template for creating headers
type HeaderTemplate struct {
	Name        string
	Key         string
	Placeholder string
}

// Model represents the application state
type Model struct {
	ActivePane           Pane
	SelectedMethod       int
	SelectedHeader       int
	URLInput             textinput.Model
	BodyInput            textarea.Model
	ResponseViewport     viewport.Model
	StatusCode           int
	Width                int
	Height               int
	Executing            bool
	CustomHeaders        []Header
	SelectedCustomHeader int
	HeadersMode          HeadersMode
	SelectedTemplate     int
	HeaderEditInput      textinput.Model
}

// ResponseMsg represents a message containing HTTP response data
type ResponseMsg struct {
	Body       string
	StatusCode int
	Err        error
}
