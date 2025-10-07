package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestInitialModel tests that the initial model is set up correctly
func TestInitialModel(t *testing.T) {
	m := initialModel()

	if m.activePane != URLPane {
		t.Errorf("Expected initial active pane to be URLPane, got %v", m.activePane)
	}

	if m.selectedMethod != 0 {
		t.Errorf("Expected initial selected method to be 0 (GET), got %d", m.selectedMethod)
	}

	if m.selectedHeader != 0 {
		t.Errorf("Expected initial selected header to be 0, got %d", m.selectedHeader)
	}

	if m.executing {
		t.Error("Expected executing to be false initially")
	}

	if m.urlInput.Value() != "" {
		t.Errorf("Expected URL input to be empty, got %s", m.urlInput.Value())
	}

	if m.bodyInput.Value() != "" {
		t.Errorf("Expected body input to be empty, got %s", m.bodyInput.Value())
	}
}

// TestPaneNavigation tests switching between panes with number keys
func TestPaneNavigation(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		expectedPane Pane
	}{
		{"Switch to URL pane", "1", URLPane},
		{"Switch to Method pane", "2", MethodPane},
		{"Switch to Body pane", "3", BodyPane},
		{"Switch to Header pane", "4", HeaderPane},
		{"Switch to Response pane", "5", ResponsePane},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := initialModel()
			m.width = 100
			m.height = 40

			// First, use Tab to get to MethodPane (a non-text-input pane)
			// so that number keys work for pane switching
			tabMsg := tea.KeyMsg{Type: tea.KeyTab}
			newModel, _ := m.Update(tabMsg)
			m = newModel.(model)

			// Now test number key navigation from MethodPane
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			newModel, _ = m.Update(msg)
			m = newModel.(model)

			if m.activePane != tt.expectedPane {
				t.Errorf("Expected active pane to be %v, got %v", tt.expectedPane, m.activePane)
			}
		})
	}
}

// TestMethodSelection tests navigating through HTTP methods
func TestMethodSelection(t *testing.T) {
	m := initialModel()
	m.width = 100
	m.height = 40
	m.activePane = MethodPane

	// Test moving down
	downMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
	newModel, _ := m.Update(downMsg)
	m = newModel.(model)

	if m.selectedMethod != 1 {
		t.Errorf("Expected selected method to be 1 after pressing j, got %d", m.selectedMethod)
	}

	// Test moving up
	upMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")}
	newModel, _ = m.Update(upMsg)
	m = newModel.(model)

	if m.selectedMethod != 0 {
		t.Errorf("Expected selected method to be 0 after pressing k, got %d", m.selectedMethod)
	}

	// Test boundary - can't go below 0
	newModel, _ = m.Update(upMsg)
	m = newModel.(model)

	if m.selectedMethod != 0 {
		t.Errorf("Expected selected method to stay at 0, got %d", m.selectedMethod)
	}

	// Test moving to last method
	for i := 0; i < len(httpMethods); i++ {
		newModel, _ = m.Update(downMsg)
		m = newModel.(model)
	}

	if m.selectedMethod >= len(httpMethods) {
		t.Errorf("Selected method should not exceed array bounds, got %d", m.selectedMethod)
	}
}

// TestHeaderSelection tests navigating through content types
func TestHeaderSelection(t *testing.T) {
	m := initialModel()
	m.width = 100
	m.height = 40
	m.activePane = HeaderPane

	// Test moving down
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ := m.Update(downMsg)
	m = newModel.(model)

	if m.selectedHeader != 1 {
		t.Errorf("Expected selected header to be 1 after pressing down, got %d", m.selectedHeader)
	}

	// Test moving up
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ = m.Update(upMsg)
	m = newModel.(model)

	if m.selectedHeader != 0 {
		t.Errorf("Expected selected header to be 0 after pressing up, got %d", m.selectedHeader)
	}
}

// TestWindowSizeUpdate tests that window size updates correctly update model dimensions
func TestWindowSizeUpdate(t *testing.T) {
	m := initialModel()

	msg := tea.WindowSizeMsg{Width: 120, Height: 50}
	newModel, _ := m.Update(msg)
	m = newModel.(model)

	if m.width != 120 {
		t.Errorf("Expected width to be 120, got %d", m.width)
	}

	if m.height != 50 {
		t.Errorf("Expected height to be 50, got %d", m.height)
	}

	// Check that input widths were updated
	if m.urlInput.Width < 10 {
		t.Errorf("URL input width should be set, got %d", m.urlInput.Width)
	}
}

// TestQuitKeys tests that q and ctrl+c trigger quit
func TestQuitKeys(t *testing.T) {
	tests := []struct {
		name string
		key  tea.KeyMsg
	}{
		{"Quit with q", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}},
		{"Quit with ctrl+c", tea.KeyMsg{Type: tea.KeyCtrlC}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := initialModel()
			m.width = 100
			m.height = 40

			_, cmd := m.Update(tt.key)

			if cmd == nil {
				t.Error("Expected quit command to be returned")
			}

			// Note: We can't directly compare cmd to tea.Quit as they're both functions,
			// but we verified that cmd is not nil, which indicates a quit command was returned
		})
	}
}

// TestHTTPRequestExecution tests the HTTP request execution with a mock server
func TestHTTPRequestExecution(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify method
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Verify Content-Type header
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type to be application/json, got %s", contentType)
		}

		// Read and verify body
		body, _ := io.ReadAll(r.Body)
		expectedBody := `{"test": "data"}`
		if strings.TrimSpace(string(body)) != expectedBody {
			t.Errorf("Expected body to be %s, got %s", expectedBody, string(body))
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}))
	defer server.Close()

	// Execute the request command
	cmd := executeRequest("POST", server.URL, `{"test": "data"}`, "application/json", []Header{})
	msg := cmd()

	// Check the response
	respMsg, ok := msg.(responseMsg)
	if !ok {
		t.Fatal("Expected responseMsg type")
	}

	if respMsg.err != nil {
		t.Fatalf("Expected no error, got %v", respMsg.err)
	}

	if respMsg.statusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", respMsg.statusCode)
	}

	if !strings.Contains(respMsg.body, "success") {
		t.Errorf("Expected response body to contain 'success', got %s", respMsg.body)
	}
}

// TestHTTPRequestGET tests GET request execution
func TestHTTPRequestGET(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "hello"})
	}))
	defer server.Close()

	cmd := executeRequest("GET", server.URL, "", "application/json", []Header{})
	msg := cmd()

	respMsg := msg.(responseMsg)

	if respMsg.err != nil {
		t.Fatalf("Expected no error, got %v", respMsg.err)
	}

	if respMsg.statusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", respMsg.statusCode)
	}
}

// TestHTTPRequestError tests handling of request errors
func TestHTTPRequestError(t *testing.T) {
	// Use an invalid URL to trigger an error
	cmd := executeRequest("GET", "http://invalid-url-that-does-not-exist-12345.com", "", "application/json", []Header{})
	msg := cmd()

	respMsg, ok := msg.(responseMsg)
	if !ok {
		t.Fatal("Expected responseMsg type")
	}

	if respMsg.err == nil {
		t.Error("Expected an error for invalid URL")
	}
}

// TestResponseMessageHandling tests that the model correctly handles response messages
func TestResponseMessageHandling(t *testing.T) {
	m := initialModel()
	m.width = 100
	m.height = 40
	m.executing = true

	// Test successful response
	respMsg := responseMsg{
		body:       "Test response",
		statusCode: 200,
		err:        nil,
	}

	newModel, _ := m.Update(respMsg)
	m = newModel.(model)

	if m.executing {
		t.Error("Expected executing to be false after receiving response")
	}

	viewportContent := strings.TrimSpace(m.responseViewport.View())
	if !strings.HasPrefix(viewportContent, "Test response") {
		t.Errorf("Expected response to start with 'Test response', got %s", viewportContent)
	}

	if m.statusCode != 200 {
		t.Errorf("Expected status code to be 200, got %d", m.statusCode)
	}

	// Test error response
	m.executing = true
	errorMsg := responseMsg{
		err: http.ErrServerClosed,
	}

	newModel, _ = m.Update(errorMsg)
	m = newModel.(model)

	if m.executing {
		t.Error("Expected executing to be false after receiving error")
	}

	errorContent := m.responseViewport.View()
	if !strings.Contains(errorContent, "Error") {
		t.Errorf("Expected response to contain 'Error', got %s", errorContent)
	}

	if m.statusCode != 0 {
		t.Errorf("Expected status code to be 0 on error, got %d", m.statusCode)
	}
}

// TestEnterKeyExecutesRequest tests that pressing Enter in URL pane executes request
func TestEnterKeyExecutesRequest(t *testing.T) {
	m := initialModel()
	m.width = 100
	m.height = 40
	m.activePane = URLPane
	m.urlInput.SetValue("https://example.com")

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, cmd := m.Update(enterMsg)
	m = newModel.(model)

	if !m.executing {
		t.Error("Expected executing to be true after pressing Enter")
	}

	if cmd == nil {
		t.Error("Expected a command to be returned for request execution")
	}
}

// TestEnterKeyInResponsePane tests that pressing Enter in response pane also executes request
func TestEnterKeyInResponsePane(t *testing.T) {
	m := initialModel()
	m.width = 100
	m.height = 40
	m.activePane = ResponsePane
	m.urlInput.SetValue("https://example.com")

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, cmd := m.Update(enterMsg)
	m = newModel.(model)

	if !m.executing {
		t.Error("Expected executing to be true after pressing Enter in response pane")
	}

	if cmd == nil {
		t.Error("Expected a command to be returned for request execution")
	}
}

// TestEnterKeyEmptyURL tests that pressing Enter with empty URL does nothing
func TestEnterKeyEmptyURL(t *testing.T) {
	m := initialModel()
	m.width = 100
	m.height = 40
	m.activePane = URLPane
	// URL is empty

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, cmd := m.Update(enterMsg)
	m = newModel.(model)

	if m.executing {
		t.Error("Expected executing to remain false with empty URL")
	}

	// cmd might be nil or a no-op
	_ = cmd
}

// TestJSONFormatting tests that JSON responses are formatted
func TestJSONFormatting(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Send compact JSON
		w.Write([]byte(`{"key":"value","nested":{"data":"test"}}`))
	}))
	defer server.Close()

	cmd := executeRequest("GET", server.URL, "", "application/json", []Header{})
	msg := cmd()

	respMsg := msg.(responseMsg)

	if respMsg.err != nil {
		t.Fatalf("Expected no error, got %v", respMsg.err)
	}

	// Check that JSON is formatted (contains newlines and indentation)
	if !strings.Contains(respMsg.body, "\n") {
		t.Error("Expected formatted JSON to contain newlines")
	}

	if !strings.Contains(respMsg.body, "  ") {
		t.Error("Expected formatted JSON to contain indentation")
	}
}

// TestDifferentHTTPMethods tests different HTTP methods
func TestDifferentHTTPMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != method {
					t.Errorf("Expected method %s, got %s", method, r.Method)
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			body := ""
			if method == "POST" || method == "PUT" || method == "PATCH" {
				body = `{"test": "data"}`
			}

			cmd := executeRequest(method, server.URL, body, "application/json", []Header{})
			msg := cmd()

			respMsg := msg.(responseMsg)
			if respMsg.err != nil {
				t.Errorf("Expected no error for %s, got %v", method, respMsg.err)
			}
		})
	}
}

// TestContentTypeHeader tests different content types
func TestContentTypeHeader(t *testing.T) {
	testCases := []struct {
		name        string
		contentType string
	}{
		{"JSON", "application/json"},
		{"XML", "application/xml"},
		{"Plain text", "text/plain"},
		{"Form urlencoded", "application/x-www-form-urlencoded"},
		{"Multipart", "multipart/form-data"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedCT := r.Header.Get("Content-Type")
				if receivedCT != tc.contentType {
					t.Errorf("Expected Content-Type %s, got %s", tc.contentType, receivedCT)
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			cmd := executeRequest("POST", server.URL, "test body", tc.contentType, []Header{})
			msg := cmd()

			respMsg := msg.(responseMsg)
			if respMsg.err != nil {
				t.Errorf("Expected no error, got %v", respMsg.err)
			}
		})
	}
}

// TestViewRendersWithoutPanic tests that View doesn't panic
func TestViewRendersWithoutPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("View panicked: %v", r)
		}
	}()

	m := initialModel()
	// Before window size
	view := m.View()
	if view != "Loading..." {
		t.Errorf("Expected 'Loading...' before window size, got %s", view)
	}

	// After window size
	m.width = 120
	m.height = 40
	view = m.View()

	if view == "" {
		t.Error("Expected non-empty view")
	}

	if view == "Loading..." {
		t.Error("Expected view to render after window size is set")
	}
}

// TestViewContainsExpectedElements tests that the view contains expected UI elements
func TestViewContainsExpectedElements(t *testing.T) {
	m := initialModel()
	m.width = 120
	m.height = 40

	view := m.View()

	expectedElements := []string{
		"Method",
		"URL",
		"Body",
		"Content-Type",
		"Result",
		"Send",
		"Quit",
	}

	for _, element := range expectedElements {
		if !strings.Contains(view, element) {
			t.Errorf("Expected view to contain '%s'", element)
		}
	}
}

// TestMethodListInView tests that all HTTP methods appear in the view
func TestMethodListInView(t *testing.T) {
	m := initialModel()
	m.width = 120
	m.height = 40
	m.activePane = MethodPane

	view := m.View()

	for _, method := range httpMethods {
		if !strings.Contains(view, method) {
			t.Errorf("Expected view to contain method '%s'", method)
		}
	}
}

// TestExecutingStatePreventsDuplicateRequests tests that multiple requests can't be triggered while executing
func TestExecutingStatePreventsDuplicateRequests(t *testing.T) {
	m := initialModel()
	m.width = 100
	m.height = 40
	m.activePane = URLPane
	m.urlInput.SetValue("https://example.com")
	m.executing = true // Already executing

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, cmd := m.Update(enterMsg)
	m = newModel.(model)

	if cmd != nil {
		t.Error("Expected no command when already executing")
	}
}
