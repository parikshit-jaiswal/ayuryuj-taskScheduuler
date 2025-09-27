package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"ayuryuj-task/internal/models"

	"github.com/google/uuid"
)

// HTTPExecutorService handles HTTP request execution
type HTTPExecutorService struct {
	client *http.Client
}

// NewHTTPExecutorService creates a new HTTP executor service
func NewHTTPExecutorService() *HTTPExecutorService {
	return &HTTPExecutorService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ExecuteTask executes a task's HTTP action and returns the result
func (s *HTTPExecutorService) ExecuteTask(task models.Task) *models.TaskResult {
	start := time.Now()
	runAt := start

	result := &models.TaskResult{
		ID:        uuid.New(),
		TaskID:    task.ID,
		RunAt:     runAt,
		CreatedAt: time.Now(),
	}

	// Create HTTP request
	req, err := s.createHTTPRequest(task.Action)
	if err != nil {
		result.Success = false
		result.StatusCode = 0
		errorMsg := fmt.Sprintf("Failed to create HTTP request: %v", err)
		result.ErrorMessage = &errorMsg
		result.DurationMs = time.Since(start).Milliseconds()
		return result
	}

	// Execute HTTP request
	resp, err := s.client.Do(req)
	duration := time.Since(start)
	result.DurationMs = duration.Milliseconds()

	if err != nil {
		result.Success = false
		result.StatusCode = 0
		errorMsg := fmt.Sprintf("HTTP request failed: %v", err)
		result.ErrorMessage = &errorMsg
		return result
	}
	defer resp.Body.Close()

	// Read response
	result.StatusCode = resp.StatusCode
	result.Success = resp.StatusCode >= 200 && resp.StatusCode < 300

	// Read response headers
	responseHeaders := make(map[string]string)
	for name, values := range resp.Header {
		if len(values) > 0 {
			responseHeaders[name] = values[0]
		}
	}

	headersJSON, _ := json.Marshal(responseHeaders)
	result.ResponseHeaders = json.RawMessage(headersJSON)

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to read response body: %v", err)
		result.ErrorMessage = &errorMsg
		result.Success = false
	} else {
		result.ResponseBody = string(bodyBytes)
	}

	return result
}

// createHTTPRequest creates an HTTP request from an action
func (s *HTTPExecutorService) createHTTPRequest(action models.Action) (*http.Request, error) {
	var body io.Reader

	// Handle payload
	if len(action.Payload) > 0 {
		body = bytes.NewReader(action.Payload)
	}

	// Create request
	req, err := http.NewRequest(action.Method, action.URL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	if action.Headers != nil {
		for key, value := range action.Headers {
			req.Header.Set(key, value)
		}
	}

	// Set default content type if payload exists and no content type is set
	if len(action.Payload) > 0 && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}
