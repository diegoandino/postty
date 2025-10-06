package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Pane represents which pane is currently focused
type Pane int

const (
	URLPane      Pane = iota + 1 // 1
	MethodPane                   // 2
	BodyPane                     // 3
	HeaderPane                   // 4
	ResponsePane                 // 5
)

var httpMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
var contentTypes = []string{
	"application/json",
	"application/xml",
	"text/plain",
	"application/x-www-form-urlencoded",
	"multipart/form-data",
}

type model struct {
	activePane       Pane
	selectedMethod   int
	selectedHeader   int
	urlInput         textinput.Model
	bodyInput        textarea.Model
	responseViewport viewport.Model
	statusCode       int
	width            int
	height           int
	executing        bool
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "https://api.example.com/endpoint"
	ti.Focus()
	ti.CharLimit = 500
	ti.Width = 40 // will be updated on WindowSizeMsg

	ta := textarea.New()
	ta.Placeholder = "Request body (JSON, XML, etc.)"
	ta.SetWidth(40) // will be updated on WindowSizeMsg
	ta.SetHeight(8) // will be updated on WindowSizeMsg

	vp := viewport.New(40, 10)
	vp.SetContent("")

	return model{
		activePane:       URLPane,
		selectedMethod:   0,
		selectedHeader:   0,
		urlInput:         ti,
		bodyInput:        ta,
		responseViewport: vp,
		statusCode:       0,
		width:            0,
		height:           0,
		executing:        false,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

type responseMsg struct {
	body       string
	statusCode int
	err        error
}

func executeRequest(method, url, body, contentType string) tea.Cmd {
	return func() tea.Msg {
		var req *http.Request
		var err error

		if body != "" && (method == "POST" || method == "PUT" || method == "PATCH") {
			req, err = http.NewRequest(method, url, strings.NewReader(body))
		} else {
			req, err = http.NewRequest(method, url, nil)
		}

		if err != nil {
			return responseMsg{err: err}
		}

		req.Header.Set("Content-Type", contentType)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return responseMsg{err: err}
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return responseMsg{err: err}
		}

		// Try to format JSON
		var prettyJSON bytes.Buffer
		if strings.Contains(resp.Header.Get("Content-Type"), "json") {
			if err := json.Indent(&prettyJSON, bodyBytes, "", "  "); err == nil {
				bodyBytes = prettyJSON.Bytes()
			}
		}

		return responseMsg{
			body:       string(bodyBytes),
			statusCode: resp.StatusCode,
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		rightColumnWidth := 35
		leftColumnWidth := m.width - rightColumnWidth - 4

		// ensure min widths
		if leftColumnWidth < 40 {
			leftColumnWidth = 40
			rightColumnWidth = m.width - leftColumnWidth - 4
			if rightColumnWidth < 20 {
				rightColumnWidth = 20
			}
		}

		// URL input fills left column - subtract border (2) + padding (2) + extra (4) = 8
		urlInputWidth := leftColumnWidth - 8
		if urlInputWidth < 20 {
			urlInputWidth = 20
		}
		m.urlInput.Width = urlInputWidth

		// body input fills left column
		bodyInputWidth := leftColumnWidth - 8 // account for borders and padding
		if bodyInputWidth < 20 {
			bodyInputWidth = 20
		}
		m.bodyInput.SetWidth(bodyInputWidth)

		// calculate and set body height - new layout with result at bottom
		// Body height should match: rightColumnTotal - urlHeight
		availableHeight := m.height - 2
		urlHeight := 5
		methodHeight := 12
		headerHeight := 9
		resultHeight := 18

		// Adjust if terminal is too small
		if availableHeight < 32 {
			resultHeight = 8
			if availableHeight < 28 {
				methodHeight = 10
				headerHeight = 7
			}
		}

		rightColumnTotalHeight := methodHeight + headerHeight
		bodyHeight := rightColumnTotalHeight - urlHeight

		if bodyHeight < 5 {
			bodyHeight = 5
		}
		m.bodyInput.SetHeight(bodyHeight)

		resultPaneWidth := leftColumnWidth + rightColumnWidth
		viewportWidth := resultPaneWidth - 4 // Subtract borders and padding
		if viewportWidth < 20 {
			viewportWidth = 20
		}
		m.responseViewport.Width = viewportWidth
		m.responseViewport.Height = resultHeight - 2 // subtract title line and border

		return m, nil

	case responseMsg:
		m.executing = false
		if msg.err != nil {
			m.responseViewport.SetContent(fmt.Sprintf("Error: %v", msg.err))
			m.statusCode = 0
		} else {
			m.responseViewport.SetContent(msg.body)
			m.statusCode = msg.statusCode
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}

		switch msg.String() {
		case "tab":
			m.activePane++
			if m.activePane > ResponsePane {
				m.activePane = URLPane
			}

			m.urlInput.Blur()
			m.bodyInput.Blur()
			if m.activePane == URLPane {
				m.urlInput.Focus()
				return m, textinput.Blink
			} else if m.activePane == BodyPane {
				m.bodyInput.Focus()
				return m, textarea.Blink
			}
			return m, nil
		case "shift+tab":
			// cycle to previous pane
			m.activePane--
			if m.activePane < URLPane {
				m.activePane = ResponsePane
			}
			// update focus
			m.urlInput.Blur()
			m.bodyInput.Blur()
			if m.activePane == URLPane {
				m.urlInput.Focus()
				return m, textinput.Blink
			} else if m.activePane == BodyPane {
				m.bodyInput.Focus()
				return m, textarea.Blink
			}
			return m, nil
		}

		// In these panes, q and 1-5 should be typed, not trigger actions
		if m.activePane == URLPane || m.activePane == BodyPane {
			// Check for Alt+Enter first (using key type for modifier keys)
			if msg.Type == tea.KeyEnter && msg.Alt {
				// Alt+Enter executes request from Body pane
				if m.activePane == BodyPane && m.urlInput.Value() != "" && !m.executing {
					m.executing = true
					m.responseViewport.SetContent("Executing request...")
					return m, executeRequest(
						httpMethods[m.selectedMethod],
						m.urlInput.Value(),
						m.bodyInput.Value(),
						contentTypes[m.selectedHeader],
					)
				}
				return m, nil
			}

			switch msg.String() {
			case "esc":
				return m, tea.Quit
			case "enter":
				if m.activePane == URLPane {
					// Execute request from URL pane
					if m.urlInput.Value() != "" && !m.executing {
						m.executing = true
						m.responseViewport.SetContent("Executing request...")
						return m, executeRequest(
							httpMethods[m.selectedMethod],
							m.urlInput.Value(),
							m.bodyInput.Value(),
							contentTypes[m.selectedHeader],
						)
					}
					return m, nil
				}
				// For BodyPane, enter creates newline - fall through to textarea update
			}
		} else {
			switch msg.String() {
			case "q", "esc":
				return m, tea.Quit

			case "1":
				m.activePane = URLPane
				m.urlInput.Focus()
				m.bodyInput.Blur()
				return m, textinput.Blink

			case "2":
				m.activePane = MethodPane
				m.urlInput.Blur()
				m.bodyInput.Blur()
				return m, nil

			case "3":
				m.activePane = BodyPane
				m.urlInput.Blur()
				m.bodyInput.Focus()
				return m, textarea.Blink

			case "4":
				m.activePane = HeaderPane
				m.urlInput.Blur()
				m.bodyInput.Blur()
				return m, nil

			case "5":
				m.activePane = ResponsePane
				m.urlInput.Blur()
				m.bodyInput.Blur()
				return m, nil
			}

			switch m.activePane {
			case MethodPane, HeaderPane:
				switch msg.String() {
				case "up", "k":
					if m.activePane == MethodPane {
						if m.selectedMethod > 0 {
							m.selectedMethod--
						}
					} else if m.activePane == HeaderPane {
						if m.selectedHeader > 0 {
							m.selectedHeader--
						}
					}
					return m, nil

				case "down", "j":
					if m.activePane == MethodPane {
						if m.selectedMethod < len(httpMethods)-1 {
							m.selectedMethod++
						}
					} else if m.activePane == HeaderPane {
						if m.selectedHeader < len(contentTypes)-1 {
							m.selectedHeader++
						}
					}
					return m, nil

				case "enter":
					if m.urlInput.Value() != "" && !m.executing {
						m.executing = true
						m.responseViewport.SetContent("Executing request...")
						return m, executeRequest(
							httpMethods[m.selectedMethod],
							m.urlInput.Value(),
							m.bodyInput.Value(),
							contentTypes[m.selectedHeader],
						)
					}
					return m, nil
				}

			case ResponsePane:
				switch msg.String() {
				case "enter":
					if m.urlInput.Value() != "" && !m.executing {
						m.executing = true
						m.responseViewport.SetContent("Executing request...")
						return m, executeRequest(
							httpMethods[m.selectedMethod],
							m.urlInput.Value(),
							m.bodyInput.Value(),
							contentTypes[m.selectedHeader],
						)
					}
					return m, nil
				case "up", "k":
					m.responseViewport.ScrollUp(1)
					return m, nil
				case "down", "j":
					m.responseViewport.ScrollDown(1)
					return m, nil
				case "pgup":
					m.responseViewport.HalfPageUp()
					return m, nil
				case "pgdown":
					m.responseViewport.HalfPageDown()
					return m, nil
				case "home", "g":
					m.responseViewport.GotoTop()
					return m, nil
				case "end", "G":
					m.responseViewport.GotoBottom()
					return m, nil
				}
			}
		}
	}

	switch m.activePane {
	case URLPane:
		m.urlInput, cmd = m.urlInput.Update(msg)
		cmds = append(cmds, cmd)
	case BodyPane:
		m.bodyInput, cmd = m.bodyInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// ✨ Enhanced fancy styles ✨
	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	activeBorderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("213")). // Bright pink/magenta
		Padding(0, 1).
		Bold(true)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("213")). // Bright pink
		Underline(true)

	paneNumberStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("51")). // Cyan
		Bold(true)

	statusGreenStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")). // White text
		Background(lipgloss.Color("22")).
		Bold(true).
		Padding(0, 1)

	statusRedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Background(lipgloss.Color("52")).
		Bold(true).
		Padding(0, 1)

	statusYellowStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")). // White text
		Background(lipgloss.Color("58")).
		Bold(true).
		Padding(0, 1)

	selectedItemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")). // Pink
		Bold(true)

	// layout: URL + Body on left, Method + Content-Type on right, Result full width at bottom
	rightColumnWidth := 40
	leftColumnWidth := m.width - rightColumnWidth - 4 // Leave 4-char buffer to prevent overflow

	// Ensure we have minimum widths
	if leftColumnWidth < 40 {
		leftColumnWidth = 40
		rightColumnWidth = m.width - leftColumnWidth - 4
		if rightColumnWidth < 20 {
			rightColumnWidth = 20
		}
	}

	// Calculate heights based on available terminal space
	// Reserve space for help bar (2 lines)
	availableHeight := m.height - 2

	// Fixed heights for panes
	urlHeight := 5
	methodHeight := 12 // Height for Method pane (7 methods + title + padding)
	headerHeight := 10 // Height for Content-Type pane (5 types + title + padding)
	resultHeight := 18 // Height for result pane

	// Adjust if terminal is too small
	if availableHeight < 32 {
		resultHeight = 8
		if availableHeight < 28 {
			methodHeight = 10
			headerHeight = 7
		}
	}

	// Calculate body height to match the right column total height
	// Right column = Method + Content-Type
	// Left column = URL + Body
	// Body height = (Method + Content-Type) - URL
	rightColumnTotalHeight := methodHeight + headerHeight
	bodyHeight := rightColumnTotalHeight - urlHeight

	// Ensure minimum heights
	if bodyHeight < 5 {
		bodyHeight = 5
	}
	if resultHeight < 5 {
		resultHeight = 5
	}

	// Build URL pane (1 - top left, wide)
	urlTitle := paneNumberStyle.Render("[1] ") + titleStyle.Render("URL")
	urlContent := urlTitle + "\n" + m.urlInput.View()
	urlStyle := borderStyle
	if m.activePane == URLPane {
		urlStyle = activeBorderStyle
	}
	urlPane := urlStyle.Width(leftColumnWidth).Height(urlHeight).Render(urlContent)

	// Build Method pane (2 - top right, narrow)
	methodTitle := paneNumberStyle.Render("[2] ") + titleStyle.Render("Method")
	methodContent := methodTitle + "\n"
	for i, method := range httpMethods {
		if i == m.selectedMethod {
			methodContent += selectedItemStyle.Render("▶ "+method) + "\n"
		} else {
			methodContent += "  " + method + "\n"
		}
	}
	methodStyle := borderStyle
	if m.activePane == MethodPane {
		methodStyle = activeBorderStyle
	}
	methodPane := methodStyle.Width(rightColumnWidth).Height(methodHeight).Render(methodContent)

	// Build Body pane (3 - middle left)
	bodyTitle := paneNumberStyle.Render("[3] ") + titleStyle.Render("Body")
	bodyContent := bodyTitle + "\n" + m.bodyInput.View()
	bodyStyle := borderStyle
	if m.activePane == BodyPane {
		bodyStyle = activeBorderStyle
	}
	bodyPane := bodyStyle.Width(leftColumnWidth).Height(bodyHeight).Render(bodyContent)

	// Build Content-Type pane (4 - right side, below methods)
	headerTitle := paneNumberStyle.Render("[4] ") + titleStyle.Render("Content-Type")
	headerContent := headerTitle + "\n"
	for i, ct := range contentTypes {
		if i == m.selectedHeader {
			headerContent += selectedItemStyle.Render("▶ "+ct) + "\n"
		} else {
			headerContent += "  " + ct + "\n"
		}
	}
	headerStyle := borderStyle
	if m.activePane == HeaderPane {
		headerStyle = activeBorderStyle
	}
	headerPane := headerStyle.Width(rightColumnWidth).Height(headerHeight).Render(headerContent)

	// Build Result pane (5 - bottom, full width)
	resultTitle := paneNumberStyle.Render("[5] ") + titleStyle.Render("Result")
	if m.statusCode > 0 {
		statusStyle := statusGreenStyle
		if m.statusCode >= 400 {
			statusStyle = statusRedStyle
		} else if m.statusCode >= 300 {
			statusStyle = statusYellowStyle
		}
		resultTitle += " " + statusStyle.Render(fmt.Sprintf("[%d]", m.statusCode))
	}

	resultContent := resultTitle + "\n" + m.responseViewport.View()
	resultStyle := borderStyle

	if m.activePane == ResponsePane {
		resultStyle = activeBorderStyle
	}

	// Result pane width should match the combined width of left + right columns
	resultPaneWidth := (leftColumnWidth + rightColumnWidth) + 2
	resultPane := resultStyle.Width(resultPaneWidth).Height(resultHeight).Render(resultContent)

	// Build left column: URL + Body
	leftColumn := lipgloss.JoinVertical(lipgloss.Left, urlPane, bodyPane)

	// Build right column: Method + Content-Type
	rightColumn := lipgloss.JoinVertical(lipgloss.Left, methodPane, headerPane)

	// Join left and right columns for top section
	topSection := lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, rightColumn)

	// Join top section with result pane at bottom
	mainView := lipgloss.JoinVertical(lipgloss.Left, topSection, resultPane)

	// Build help bar with enhanced styling
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("51")). // Cyan
		Background(lipgloss.Color("235")).
		Padding(0, 1).
		Bold(true)

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("213")). // Pink
		Bold(true)

	help := helpStyle.Render(
		keyStyle.Render("Tab") + " Next Pane │ " +
			keyStyle.Render("1-5") + " Jump │ " +
			keyStyle.Render("↑↓jk") + " Scroll │ " +
			keyStyle.Render("Enter") + "/" + keyStyle.Render("Alt+Enter") + " Send │ " +
			keyStyle.Render("esc") + "/" + keyStyle.Render("q") + " Quit",
	)

	return mainView + "\n" + help
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
