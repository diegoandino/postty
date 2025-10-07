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

type Pane int

const (
	URLPane Pane = iota + 1
	MethodPane
	BodyPane
	HeaderPane
	ResponsePane
	HeadersPane
)

var httpMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
var contentTypes = []string{
	"application/json",
	"application/xml",
	"text/plain",
	"application/x-www-form-urlencoded",
	"multipart/form-data",
}

type Header struct {
	Key   string
	Value string
}

type HeaderTemplate struct {
	Name        string
	Key         string
	Placeholder string
}

var headerTemplates = []HeaderTemplate{
	{Name: "Authorization (Bearer)", Key: "Authorization", Placeholder: "Bearer <your-token>"},
	{Name: "Authorization (Basic)", Key: "Authorization", Placeholder: "Basic <base64-credentials>"},
	{Name: "API Key", Key: "X-API-Key", Placeholder: "<your-api-key>"},
	{Name: "Cookie", Key: "Cookie", Placeholder: "session_id=<value>"},
	{Name: "User Agent", Key: "User-Agent", Placeholder: "MyApp/1.0"},
	{Name: "Accept", Key: "Accept", Placeholder: "application/json"},
	{Name: "Custom Header", Key: "", Placeholder: ""},
}

type HeadersMode int

const (
	HeadersViewMode HeadersMode = iota
	HeadersAddMode
	HeadersEditMode
)

type model struct {
	activePane           Pane
	selectedMethod       int
	selectedHeader       int
	urlInput             textinput.Model
	bodyInput            textarea.Model
	responseViewport     viewport.Model
	statusCode           int
	width                int
	height               int
	executing            bool
	customHeaders        []Header
	selectedCustomHeader int
	headersMode          HeadersMode
	selectedTemplate     int
	headerEditInput      textinput.Model
}

func initialModel() model {
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

	defaultHeaders := []Header{}

	return model{
		activePane:           URLPane,
		selectedMethod:       0,
		selectedHeader:       0,
		urlInput:             ti,
		bodyInput:            ta,
		responseViewport:     vp,
		statusCode:           0,
		width:                0,
		height:               0,
		executing:            false,
		customHeaders:        defaultHeaders,
		selectedCustomHeader: 0,
		headersMode:          HeadersViewMode,
		selectedTemplate:     0,
		headerEditInput:      hei,
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

func executeRequest(method, url, body, contentType string, customHeaders []Header) tea.Cmd {
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

		for _, header := range customHeaders {
			if header.Key != "" && header.Value != "" {
				req.Header.Set(header.Key, header.Value)
			}
		}

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

		if leftColumnWidth < 40 {
			leftColumnWidth = 40
			rightColumnWidth = m.width - leftColumnWidth - 4
			if rightColumnWidth < 20 {
				rightColumnWidth = 20
			}
		}

		urlInputWidth := leftColumnWidth - 8
		if urlInputWidth < 20 {
			urlInputWidth = 20
		}
		m.urlInput.Width = urlInputWidth

		bodyInputWidth := leftColumnWidth - 8
		if bodyInputWidth < 20 {
			bodyInputWidth = 20
		}
		m.bodyInput.SetWidth(bodyInputWidth)

		availableHeight := m.height - 2
		borderOverhead := 6
		contentHeight := availableHeight - borderOverhead

		if contentHeight < 24 {
			contentHeight = 24
		}

		urlContentHeight := 3
		topSectionContentHeight := contentHeight - int(float64(contentHeight)*0.4)
		if topSectionContentHeight < 14 {
			topSectionContentHeight = 14
		}

		resultContentHeight := contentHeight - topSectionContentHeight
		if resultContentHeight < 6 {
			resultContentHeight = 6
		}

		methodContentHeight := int(float64(topSectionContentHeight) * 0.55)
		if methodContentHeight < 8 {
			methodContentHeight = 8
		}
		headerContentHeight := (topSectionContentHeight - methodContentHeight) - 2
		if headerContentHeight < 5 {
			headerContentHeight = 5
		}

		rightColumnContentHeight := methodContentHeight + headerContentHeight
		bodyContentHeight := rightColumnContentHeight - urlContentHeight
		if bodyContentHeight < 5 {
			bodyContentHeight = 5
		}
		m.bodyInput.SetHeight(bodyContentHeight)

		resultPaneWidth := leftColumnWidth + rightColumnWidth
		viewportWidth := resultPaneWidth - 4
		if viewportWidth < 20 {
			viewportWidth = 20
		}
		m.responseViewport.Width = viewportWidth
		m.responseViewport.Height = resultContentHeight - 1

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
			if m.activePane > HeadersPane {
				m.activePane = URLPane
			}

			m.urlInput.Blur()
			m.bodyInput.Blur()
			m.headerEditInput.Blur()
			m.headersMode = HeadersViewMode
			if m.activePane == URLPane {
				m.urlInput.Focus()
				return m, textinput.Blink
			} else if m.activePane == BodyPane {
				m.bodyInput.Focus()
				return m, textarea.Blink
			}
			return m, nil
		case "shift+tab":
			m.activePane--
			if m.activePane < URLPane {
				m.activePane = HeadersPane
			}
			m.urlInput.Blur()
			m.bodyInput.Blur()
			m.headerEditInput.Blur()
			m.headersMode = HeadersViewMode
			if m.activePane == URLPane {
				m.urlInput.Focus()
				return m, textinput.Blink
			} else if m.activePane == BodyPane {
				m.bodyInput.Focus()
				return m, textarea.Blink
			}
			return m, nil
		}

		if m.activePane == URLPane || m.activePane == BodyPane || (m.activePane == HeadersPane && m.headersMode == HeadersEditMode) {
			if msg.Type == tea.KeyEnter && msg.Alt {
				if m.activePane == BodyPane && m.urlInput.Value() != "" && !m.executing {
					m.executing = true
					m.responseViewport.SetContent("Executing request...")
					return m, executeRequest(
						httpMethods[m.selectedMethod],
						m.urlInput.Value(),
						m.bodyInput.Value(),
						contentTypes[m.selectedHeader],
						m.customHeaders,
					)
				}
				return m, nil
			}

			switch msg.String() {
			case "esc":
				if m.activePane == HeadersPane && m.headersMode == HeadersEditMode {
					m.headersMode = HeadersViewMode
					m.headerEditInput.Blur()
					return m, nil
				}
				return m, tea.Quit
			case "enter":
				if m.activePane == HeadersPane && m.headersMode == HeadersEditMode {
					if len(m.customHeaders) > 0 {
						m.customHeaders[m.selectedCustomHeader].Value = m.headerEditInput.Value()
						m.headersMode = HeadersViewMode
						m.headerEditInput.Blur()
					}
					return m, nil
				}

				if m.activePane == URLPane {
					if m.urlInput.Value() != "" && !m.executing {
						m.executing = true
						m.responseViewport.SetContent("Executing request...")
						return m, executeRequest(
							httpMethods[m.selectedMethod],
							m.urlInput.Value(),
							m.bodyInput.Value(),
							contentTypes[m.selectedHeader],
							m.customHeaders,
						)
					}
					return m, nil
				}
			}
		} else {
			switch msg.String() {
			case "q", "esc":
				if m.activePane == HeadersPane && (m.headersMode == HeadersAddMode || m.headersMode == HeadersEditMode) {
					break
				}
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
				m.headerEditInput.Blur()
				m.headersMode = HeadersViewMode
				return m, nil

			case "6":
				m.activePane = HeadersPane
				m.urlInput.Blur()
				m.bodyInput.Blur()
				m.headerEditInput.Blur()
				m.headersMode = HeadersViewMode
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
							m.customHeaders,
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
							m.customHeaders,
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

			case HeadersPane:
				switch m.headersMode {
				case HeadersViewMode:
					switch msg.String() {
					case "up", "k":
						if m.selectedCustomHeader > 0 {
							m.selectedCustomHeader--
						}
						return m, nil
					case "down", "j":
						if m.selectedCustomHeader < len(m.customHeaders)-1 {
							m.selectedCustomHeader++
						}
						return m, nil
					case "a", "n":
						m.headersMode = HeadersAddMode
						m.selectedTemplate = 0
						return m, nil
					case "d", "x":
						if len(m.customHeaders) > 0 {
							m.customHeaders = append(m.customHeaders[:m.selectedCustomHeader], m.customHeaders[m.selectedCustomHeader+1:]...)
							if m.selectedCustomHeader >= len(m.customHeaders) && len(m.customHeaders) > 0 {
								m.selectedCustomHeader = len(m.customHeaders) - 1
							}
							if len(m.customHeaders) == 0 {
								m.selectedCustomHeader = 0
							}
						}
						return m, nil
					case "e", "enter":
						if len(m.customHeaders) > 0 {
							m.headersMode = HeadersEditMode
							m.headerEditInput.SetValue(m.customHeaders[m.selectedCustomHeader].Value)
							m.headerEditInput.Focus()
							return m, textinput.Blink
						}
						return m, nil
					case "esc":
						return m, tea.Quit
					}

				case HeadersAddMode:
					switch msg.String() {
					case "up", "k":
						if m.selectedTemplate > 0 {
							m.selectedTemplate--
						}
						return m, nil
					case "down", "j":
						if m.selectedTemplate < len(headerTemplates)-1 {
							m.selectedTemplate++
						}
						return m, nil
					case "enter":
						template := headerTemplates[m.selectedTemplate]
						if template.Key == "" {
							m.customHeaders = append(m.customHeaders, Header{Key: "Custom-Header", Value: ""})
							m.selectedCustomHeader = len(m.customHeaders) - 1
							m.headersMode = HeadersEditMode
							m.headerEditInput.SetValue("")
							m.headerEditInput.Focus()
							return m, textinput.Blink
						} else {
							m.customHeaders = append(m.customHeaders, Header{
								Key:   template.Key,
								Value: template.Placeholder,
							})
							m.selectedCustomHeader = len(m.customHeaders) - 1
							m.headersMode = HeadersViewMode
						}
						return m, nil
					case "esc":
						m.headersMode = HeadersViewMode
						return m, nil
					}

				case HeadersEditMode:
					switch msg.String() {
					case "enter":
						if len(m.customHeaders) > 0 {
							m.customHeaders[m.selectedCustomHeader].Value = m.headerEditInput.Value()
							m.headersMode = HeadersViewMode
							m.headerEditInput.Blur()
						}
						return m, nil
					case "esc":
						m.headersMode = HeadersViewMode
						m.headerEditInput.Blur()
						return m, nil
					}
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
	case HeadersPane:
		if m.headersMode == HeadersEditMode {
			m.headerEditInput, cmd = m.headerEditInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	activeBorderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("213")).
		Padding(0, 1).
		Bold(true)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("213")).
		Underline(true)

	paneNumberStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("51")).
		Bold(true)

	statusGreenStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Background(lipgloss.Color("22")).
		Bold(true).
		Padding(0, 1)

	statusRedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Background(lipgloss.Color("52")).
		Bold(true).
		Padding(0, 1)

	statusYellowStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Background(lipgloss.Color("58")).
		Bold(true).
		Padding(0, 1)

	selectedItemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		Bold(true)

	rightColumnWidth := 40
	leftColumnWidth := m.width - rightColumnWidth - 4

	if leftColumnWidth < 40 {
		leftColumnWidth = 40
		rightColumnWidth = m.width - leftColumnWidth - 4
		if rightColumnWidth < 20 {
			rightColumnWidth = 20
		}
	}

	availableHeight := m.height - 2
	borderOverhead := 6
	contentHeight := availableHeight - borderOverhead

	if contentHeight < 24 {
		contentHeight = 24
	}

	urlHeight := 3
	sectionHeight := contentHeight / 3
	if sectionHeight < 8 {
		sectionHeight = 8
	}

	methodHeight := sectionHeight
	if methodHeight < 8 {
		methodHeight = 8
	}

	headerHeight := sectionHeight - 3
	if headerHeight < 5 {
		headerHeight = 5
	}

	headersHeight := contentHeight - methodHeight - headerHeight
	if headersHeight < 6 {
		headersHeight = 6
	}

	bodyHeight := methodHeight - urlHeight
	if bodyHeight < 5 {
		bodyHeight = 5
	}

	resultHeight := (contentHeight - urlHeight - bodyHeight) - 10
	if resultHeight < 8 {
		resultHeight = 8
	}

	urlTitle := paneNumberStyle.Render("[1] ") + titleStyle.Render("URL")
	urlContent := urlTitle + "\n" + m.urlInput.View()
	urlStyle := borderStyle
	if m.activePane == URLPane {
		urlStyle = activeBorderStyle
	}
	urlPane := urlStyle.Width(leftColumnWidth).Height(urlHeight).Render(urlContent)

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

	bodyTitle := paneNumberStyle.Render("[3] ") + titleStyle.Render("Body")
	bodyContent := bodyTitle + "\n" + m.bodyInput.View()
	bodyStyle := borderStyle
	if m.activePane == BodyPane {
		bodyStyle = activeBorderStyle
	}
	bodyPane := bodyStyle.Width(leftColumnWidth).Height(bodyHeight).Render(bodyContent)

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

	resultPaneWidth := leftColumnWidth
	resultPane := resultStyle.Width(resultPaneWidth).Height(resultHeight).Render(resultContent)

	headersTitle := paneNumberStyle.Render("[6] ") + titleStyle.Render("Custom Headers")
	headersContent := headersTitle + "\n"

	switch m.headersMode {
	case HeadersViewMode:
		if len(m.customHeaders) == 0 {
			headersContent += "  (no headers)\n"
			headersContent += "\n"
			headersContent += "  Press 'a' to add\n"
		} else {
			for i, h := range m.customHeaders {
				prefix := "  "
				if i == m.selectedCustomHeader {
					prefix = selectedItemStyle.Render("▶ ")
				}
				headerLine := fmt.Sprintf("%s: %s", h.Key, h.Value)
				if h.Key == "" && h.Value == "" {
					headerLine = "(empty)"
				}
				headersContent += prefix + headerLine + "\n"
			}
			headersContent += "\n"
			headersContent += "  a: add | d: delete | e: edit\n"
		}

	case HeadersAddMode:
		headersContent += "  Select header type:\n"
		headersContent += "\n"
		for i, template := range headerTemplates {
			prefix := "  "
			if i == m.selectedTemplate {
				prefix = selectedItemStyle.Render("▶ ")
			}
			headersContent += prefix + template.Name + "\n"
		}
		headersContent += "\n"
		headersContent += "  Enter: select | Esc: cancel\n"

	case HeadersEditMode:
		if len(m.customHeaders) > 0 {
			header := m.customHeaders[m.selectedCustomHeader]
			headersContent += fmt.Sprintf("  Editing: %s\n", header.Key)
			headersContent += "\n"
			headersContent += "  Value:\n"
			headersContent += "  " + m.headerEditInput.View() + "\n"
			headersContent += "\n"
			headersContent += "  Enter: save | Esc: cancel\n"
		}
	}

	headersStyle := borderStyle
	if m.activePane == HeadersPane {
		headersStyle = activeBorderStyle
	}
	headersPane := headersStyle.Width(rightColumnWidth).Height(headersHeight).Render(headersContent)

	leftColumn := lipgloss.JoinVertical(lipgloss.Left, urlPane, bodyPane, resultPane)
	rightColumn := lipgloss.JoinVertical(lipgloss.Left, methodPane, headerPane, headersPane)
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, rightColumn)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("51")).
		Background(lipgloss.Color("235")).
		Padding(0, 1).
		Bold(true)

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("213")).
		Bold(true)

	help := helpStyle.Render(
		keyStyle.Render("Tab") + " Next Pane │ " +
			keyStyle.Render("1-6") + " Jump │ " +
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
