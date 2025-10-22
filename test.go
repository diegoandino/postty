package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"postty/src/model"
	"postty/src/services"
	"postty/src/types"
)

// TestInitialModel tests that the initial model is set up correctly
func TestInitialModel(t *testing.T) {
	m := model.New()

	if m.ActivePane != types.URLPane {
		t.Errorf("Expected initial active pane to be URLPane, got %v", m.ActivePane)
	}

	if m.SelectedMethod != 0 {
		t.Errorf("Expected initial selected method to be 0 (GET), got %d", m.SelectedMethod)
	}

	if m.SelectedHeader != 0 {
		t.Errorf("Expected initial selected header to be 0, got %d", m.SelectedHeader)
	}

	if m.Executing {
		t.Error("Expected executing to be false initially")
	}

	if m.URLInput.Value() != "" {
		t.Errorf("Expected URL input to be empty, got %s", m.URLInput.Value())
	}

	if m.BodyInput.Value() != "" {
		t.Errorf("Expected body input to be empty, got %s", m.BodyInput.Value())
	}
}

// TestPaneNavigation tests switching between panes with number keys
func TestPaneNavigation(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		expectedPane types.Pane
	}{
		{"Switch to URL pane", "1", types.URLPane},
		{"Switch to Method pane", "2", types.MethodPane},
		{"Switch to Body pane", "3", types.BodyPane},
		{"Switch to Header pane", "4", types.HeaderPane},
		{"Switch to Response pane", "5", types.ResponsePane},
		{"Switch to Headers pane", "6", types.HeadersPane},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := app{model: model.New()}
			a.model.Width = 100
			a.model.Height = 40

			// First, use Tab to get to MethodPane (a non-text-input pane)
			// so that number keys work for pane switching
			tabMsg := tea.KeyMsg{Type: tea.KeyTab}
			newApp, _ := a.Update(tabMsg)
			a = newApp.(app)

			// Now test number key navigation from MethodPane
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			newApp, _ = a.Update(msg)
			a = newApp.(app)

			if a.model.ActivePane != tt.expectedPane {
				t.Errorf("Expected active pane to be %v, got %v", tt.expectedPane, a.model.ActivePane)
			}
		})
	}
}

// TestMethodSelection tests navigating through HTTP methods
func TestMethodSelection(t *testing.T) {
	a := app{model: model.New()}
	a.model.Width = 100
	a.model.Height = 40
	a.model.ActivePane = types.MethodPane

	// Test moving down
	downMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
	newApp, _ := a.Update(downMsg)
	a = newApp.(app)

	if a.model.SelectedMethod != 1 {
		t.Errorf("Expected selected method to be 1 after pressing j, got %d", a.model.SelectedMethod)
	}

	// Test moving up
	upMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")}
	newApp, _ = a.Update(upMsg)
	a = newApp.(app)

	if a.model.SelectedMethod != 0 {
		t.Errorf("Expected selected method to be 0 after pressing k, got %d", a.model.SelectedMethod)
	}

	// Test boundary - can't go below 0
	newApp, _ = a.Update(upMsg)
	a = newApp.(app)

	if a.model.SelectedMethod != 0 {
		t.Errorf("Expected selected method to stay at 0, got %d", a.model.SelectedMethod)
	}

	// Test moving to last method
	for i := 0; i < len(types.HTTPMethods); i++ {
		newApp, _ = a.Update(downMsg)
		a = newApp.(app)
	}

	if a.model.SelectedMethod >= len(types.HTTPMethods) {
		t.Errorf("Selected method should not exceed array bounds, got %d", a.model.SelectedMethod)
	}
}

// TestHeaderSelection tests navigating through content types
func TestHeaderSelection(t *testing.T) {
	a := app{model: model.New()}
	a.model.Width = 100
	a.model.Height = 40
	a.model.ActivePane = types.HeaderPane

	// Test moving down
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	newApp, _ := a.Update(downMsg)
	a = newApp.(app)

	if a.model.SelectedHeader != 1 {
		t.Errorf("Expected selected header to be 1 after pressing down, got %d", a.model.SelectedHeader)
	}

	// Test moving up
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	newApp, _ = a.Update(upMsg)
	a = newApp.(app)

	if a.model.SelectedHeader != 0 {
		t.Errorf("Expected selected header to be 0 after pressing up, got %d", a.model.SelectedHeader)
	}
}

// TestWindowSizeUpdate tests that window size updates correctly update model dimensions
func TestWindowSizeUpdate(t *testing.T) {
	a := app{model: model.New()}

	msg := tea.WindowSizeMsg{Width: 120, Height: 50}
	newApp, _ := a.Update(msg)
	a = newApp.(app)

	if a.model.Width != 120 {
		t.Errorf("Expected width to be 120, got %d", a.model.Width)
	}

	if a.model.Height != 50 {
		t.Errorf("Expected height to be 50, got %d", a.model.Height)
	}

	// Check that input widths were updated
	if a.model.URLInput.Width < 10 {
		t.Errorf("URL input width should be set, got %d", a.model.URLInput.Width)
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
	cmd := services.ExecuteRequest("POST", server.URL, `{"test": "data"}`, "application/json", []types.Header{})
	msg := cmd()

	// Check the response
	respMsg, ok := msg.(types.ResponseMsg)
	if !ok {
		t.Fatal("Expected types.ResponseMsg type")
	}

	if respMsg.Err != nil {
		t.Fatalf("Expected no error, got %v", respMsg.Err)
	}

	if respMsg.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", respMsg.StatusCode)
	}

	if !strings.Contains(respMsg.Body, "success") {
		t.Errorf("Expected response body to contain 'success', got %s", respMsg.Body)
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

	cmd := services.ExecuteRequest("GET", server.URL, "", "application/json", []types.Header{})
	msg := cmd()

	respMsg := msg.(types.ResponseMsg)

	if respMsg.Err != nil {
		t.Fatalf("Expected no error, got %v", respMsg.Err)
	}

	if respMsg.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", respMsg.StatusCode)
	}
}

// TestHTTPRequestError tests handling of request errors
func TestHTTPRequestError(t *testing.T) {
	// Use an invalid URL to trigger an error
	cmd := services.ExecuteRequest("GET", "http://invalid-url-that-does-not-exist-12345.com", "", "application/json", []types.Header{})
	msg := cmd()

	respMsg, ok := msg.(types.ResponseMsg)
	if !ok {
		t.Fatal("Expected types.ResponseMsg type")
	}

	if respMsg.Err == nil {
		t.Error("Expected an error for invalid URL")
	}
}

// TestResponseMessageHandling tests that the model correctly handles response messages
func TestResponseMessageHandling(t *testing.T) {
	a := app{model: model.New()}
	a.model.Width = 100
	a.model.Height = 40
	a.model.Executing = true

	// Test successful response
	respMsg := types.ResponseMsg{
		Body:       "Test response",
		StatusCode: 200,
		Err:        nil,
	}

	newApp, _ := a.Update(respMsg)
	a = newApp.(app)

	if a.model.Executing {
		t.Error("Expected executing to be false after receiving response")
	}

	viewportContent := strings.TrimSpace(a.model.ResponseViewport.View())
	if !strings.HasPrefix(viewportContent, "Test response") {
		t.Errorf("Expected response to start with 'Test response', got %s", viewportContent)
	}

	if a.model.StatusCode != 200 {
		t.Errorf("Expected status code to be 200, got %d", a.model.StatusCode)
	}

	// Test error response
	a.model.Executing = true
	errorMsg := types.ResponseMsg{
		Err: http.ErrServerClosed,
	}

	newApp, _ = a.Update(errorMsg)
	a = newApp.(app)

	if a.model.Executing {
		t.Error("Expected executing to be false after receiving error")
	}

	errorContent := a.model.ResponseViewport.View()
	if !strings.Contains(errorContent, "Error") {
		t.Errorf("Expected response to contain 'Error', got %s", errorContent)
	}

	if a.model.StatusCode != 0 {
		t.Errorf("Expected status code to be 0 on error, got %d", a.model.StatusCode)
	}
}

// TestEnterKeyExecutesRequest tests that pressing Enter in URL pane executes request
func TestEnterKeyExecutesRequest(t *testing.T) {
	a := app{model: model.New()}
	a.model.Width = 100
	a.model.Height = 40
	a.model.ActivePane = types.URLPane
	a.model.URLInput.SetValue("https://example.com")

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newApp, cmd := a.Update(enterMsg)
	a = newApp.(app)

	if !a.model.Executing {
		t.Error("Expected executing to be true after pressing Enter")
	}

	if cmd == nil {
		t.Error("Expected a command to be returned for request execution")
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

			cmd := services.ExecuteRequest(method, server.URL, body, "application/json", []types.Header{})
			msg := cmd()

			respMsg := msg.(types.ResponseMsg)
			if respMsg.Err != nil {
				t.Errorf("Expected no error for %s, got %v", method, respMsg.Err)
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

			cmd := services.ExecuteRequest("POST", server.URL, "test body", tc.contentType, []types.Header{})
			msg := cmd()

			respMsg := msg.(types.ResponseMsg)
			if respMsg.Err != nil {
				t.Errorf("Expected no error, got %v", respMsg.Err)
			}
		})
	}
}

// TestCustomHeaders tests that custom headers are properly sent
func TestCustomHeaders(t *testing.T) {
	customHeaders := []types.Header{
		{Key: "Authorization", Value: "Bearer test-token"},
		{Key: "X-Custom-Header", Value: "custom-value"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("Expected Authorization header to be 'Bearer test-token', got %s", auth)
		}

		custom := r.Header.Get("X-Custom-Header")
		if custom != "custom-value" {
			t.Errorf("Expected X-Custom-Header to be 'custom-value', got %s", custom)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cmd := services.ExecuteRequest("GET", server.URL, "", "application/json", customHeaders)
	msg := cmd()

	respMsg := msg.(types.ResponseMsg)
	if respMsg.Err != nil {
		t.Fatalf("Expected no error, got %v", respMsg.Err)
	}
}
