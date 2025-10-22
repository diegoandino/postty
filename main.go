package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"postty/src/components"
	"postty/src/handlers"
	"postty/src/model"
	"postty/src/types"
)

// app wraps the model to implement tea.Model interface
type app struct {
	model types.Model
}

// Update is the main update function for the Bubbletea application
func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	newModel, cmd := handlers.Update(msg, a.model)
	a.model = newModel
	return a, cmd
}

// View renders the application UI
func (a app) View() string {
	return components.RenderLayout(a.model)
}

// Init initializes the application
func (a app) Init() tea.Cmd {
	return model.Init()
}

func main() {
	a := app{model: model.New()}
	p := tea.NewProgram(a, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
