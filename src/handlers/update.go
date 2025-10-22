package handlers

import (
	tea "github.com/charmbracelet/bubbletea"

	"postty/src/types"
)

// Update handles all state updates for the application
func Update(msg tea.Msg, m types.Model) (types.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m = HandleWindowSize(m, msg)
		return m, nil

	case types.ResponseMsg:
		m = HandleResponse(m, msg)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}

		// Tab navigation (works from any pane)
		switch msg.String() {
		case "tab":
			return HandleTab(m)
		case "shift+tab":
			return HandleShiftTab(m)
		}

		// Handle keys based on active pane
		if m.ActivePane == types.URLPane || m.ActivePane == types.BodyPane || (m.ActivePane == types.HeadersPane && m.HeadersMode == types.HeadersEditMode) {
			// Text input panes - handle Alt+Enter for body pane execution
			if msg.Type == tea.KeyEnter && msg.Alt {
				if m.ActivePane == types.BodyPane {
					return ExecuteRequestWithHistory(m)
				}
				return m, nil
			}

			switch msg.String() {
			case "esc":
				if m.ActivePane == types.HeadersPane && m.HeadersMode == types.HeadersEditMode {
					m = HandleHeaderEditCancel(m)
					return m, nil
				}
				return m, tea.Quit
			case "enter":
				if m.ActivePane == types.HeadersPane && m.HeadersMode == types.HeadersEditMode {
					m = HandleHeaderEditSave(m)
					return m, nil
				}

				if m.ActivePane == types.URLPane {
					return ExecuteRequestWithHistory(m)
				}
			}
		} else {
			// Non-text-input panes
			switch msg.String() {
			case "q", "esc":
				if m.ActivePane == types.HeadersPane && (m.HeadersMode == types.HeadersAddMode || m.HeadersMode == types.HeadersEditMode) {
					break
				}
				return m, tea.Quit

			case "1":
				return HandleJumpToPane(m, types.URLPane)
			case "2":
				return HandleJumpToPane(m, types.MethodPane)
			case "3":
				return HandleJumpToPane(m, types.BodyPane)
			case "4":
				return HandleJumpToPane(m, types.HeaderPane)
			case "5":
				return HandleJumpToPane(m, types.ResponsePane)
			case "6":
				return HandleJumpToPane(m, types.HeadersPane)
			case "7":
				return HandleJumpToPane(m, types.HistoryPane)
			}

			// Pane-specific handlers
			switch m.ActivePane {
			case types.MethodPane, types.HeaderPane:
				switch msg.String() {
				case "up", "k":
					if m.ActivePane == types.MethodPane {
						m = HandleMethodNavigation(m, "up")
					} else if m.ActivePane == types.HeaderPane {
						m = HandleContentTypeNavigation(m, "up")
					}
					return m, nil

				case "down", "j":
					if m.ActivePane == types.MethodPane {
						m = HandleMethodNavigation(m, "down")
					} else if m.ActivePane == types.HeaderPane {
						m = HandleContentTypeNavigation(m, "down")
					}
					return m, nil

				case "enter":
					return HandleMethodExecute(m)
				}

			case types.ResponsePane:
				switch msg.String() {
				case "enter":
					return HandleMethodExecute(m)
				case "up", "k":
					m = HandleResponseScroll(m, "up")
					return m, nil
				case "down", "j":
					m = HandleResponseScroll(m, "down")
					return m, nil
				case "pgup":
					m = HandleResponseScroll(m, "pgup")
					return m, nil
				case "pgdown":
					m = HandleResponseScroll(m, "pgdown")
					return m, nil
				case "home", "g":
					m = HandleResponseScroll(m, "top")
					return m, nil
				case "end", "G":
					m = HandleResponseScroll(m, "bottom")
					return m, nil
				}

			case types.HeadersPane:
				switch m.HeadersMode {
				case types.HeadersViewMode:
					switch msg.String() {
					case "up", "k":
						m = HandleCustomHeadersNavigation(m, "up")
						return m, nil
					case "down", "j":
						m = HandleCustomHeadersNavigation(m, "down")
						return m, nil
					case "a", "n":
						m = HandleCustomHeadersAdd(m)
						return m, nil
					case "d", "x":
						m = HandleCustomHeadersDelete(m)
						return m, nil
					case "e", "enter":
						return HandleCustomHeadersEdit(m)
					case "esc":
						return m, tea.Quit
					}

				case types.HeadersAddMode:
					switch msg.String() {
					case "up", "k":
						m = HandleTemplateNavigation(m, "up")
						return m, nil
					case "down", "j":
						m = HandleTemplateNavigation(m, "down")
						return m, nil
					case "enter":
						return HandleTemplateSelect(m)
					case "esc":
						m = HandleAddModeCancel(m)
						return m, nil
					}

				case types.HeadersEditMode:
					switch msg.String() {
					case "enter":
						m = HandleHeaderEditSave(m)
						return m, nil
					case "esc":
						m = HandleHeaderEditCancel(m)
						return m, nil
					}
				}

			case types.HistoryPane:
				switch msg.String() {
				case "up", "k":
					m = HandleHistoryNavigation(m, "up")
					return m, nil
				case "down", "j":
					m = HandleHistoryNavigation(m, "down")
					return m, nil
				case "pgup":
					m.HistoryViewport.HalfPageUp()
					return m, nil
				case "pgdown":
					m.HistoryViewport.HalfPageDown()
					return m, nil
				case "home", "g":
					m.SelectedHistory = 0
					m.HistoryViewport.SetYOffset(0)
					return m, nil
				case "end", "G":
					if len(m.History) > 0 {
						m.SelectedHistory = len(m.History) - 1
						// Scroll to bottom (will be adjusted by render if needed)
						m.HistoryViewport.GotoBottom()
					}
					return m, nil
				case "enter":
					return HandleHistoryLoad(m)
				case "d", "x":
					m = HandleHistoryDelete(m)
					return m, nil
				case "esc":
					return m, tea.Quit
				}
			}
		}
	}

	// Update the active input component
	switch m.ActivePane {
	case types.URLPane:
		m.URLInput, cmd = m.URLInput.Update(msg)
		cmds = append(cmds, cmd)
	case types.BodyPane:
		m.BodyInput, cmd = m.BodyInput.Update(msg)
		cmds = append(cmds, cmd)
	case types.HeadersPane:
		if m.HeadersMode == types.HeadersEditMode {
			m.HeaderEditInput, cmd = m.HeaderEditInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}
