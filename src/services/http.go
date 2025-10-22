package services

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"postty/src/types"
)

// ExecuteRequest creates a command to execute an HTTP request
func ExecuteRequest(method, url, body, contentType string, customHeaders []types.Header) tea.Cmd {
	return func() tea.Msg {
		var req *http.Request
		var err error

		if body != "" && (method == "POST" || method == "PUT" || method == "PATCH") {
			req, err = http.NewRequest(method, url, strings.NewReader(body))
		} else {
			req, err = http.NewRequest(method, url, nil)
		}

		if err != nil {
			return types.ResponseMsg{Err: err}
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
			return types.ResponseMsg{Err: err}
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return types.ResponseMsg{Err: err}
		}

		var prettyJSON bytes.Buffer
		if strings.Contains(resp.Header.Get("Content-Type"), "json") {
			if err := json.Indent(&prettyJSON, bodyBytes, "", "  "); err == nil {
				bodyBytes = prettyJSON.Bytes()
			}
		}

		return types.ResponseMsg{
			Body:       string(bodyBytes),
			StatusCode: resp.StatusCode,
		}
	}
}
