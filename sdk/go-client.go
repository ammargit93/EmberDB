package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents an EmberDB client
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// Response represents the server response
type Response struct {
	Namespace string      `json:"namespace"`
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
}

// ErrorResponse represents an error response from the server
type ErrorResponse struct {
	Error string `json:"error"`
}

// NewClient creates a new EmberDB client
func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "http://localhost:9182"
	}
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewClientWithTimeout creates a new EmberDB client with a custom timeout
func NewClientWithTimeout(baseURL string, timeout time.Duration) *Client {
	if baseURL == "" {
		baseURL = "http://localhost:9182"
	}
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// SetKey sets a key-value pair in a namespace
func (c *Client) SetKey(namespace, key string, value interface{}) error {
	data := Response{
		Namespace: namespace,
		Key:       key,
		Value:     value,
	}
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/set", c.baseURL),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to set key: %s", string(body))
	}

	return nil
}

// GetKey retrieves a value from a namespace by key
func (c *Client) GetKey(namespace, key string) (interface{}, error) {
	url := fmt.Sprintf("%s/get/%s/%s", c.baseURL, namespace, key)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("key not found")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get key: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result["value"], nil
}

// GetKeyWithMetadata retrieves a value and its metadata
func (c *Client) GetKeyWithMetadata(namespace, key string) (*Response, error) {
	url := fmt.Sprintf("%s/get/%s/%s", c.baseURL, namespace, key)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("key not found")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get key: %s", string(body))
	}

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// UpdateKey updates an existing key's value
func (c *Client) UpdateKey(namespace, key string, value interface{}) error {
	data := Response{
		Namespace: namespace,
		Key:       key,
		Value:     value,
	}
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("%s/update", c.baseURL),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update key: %s", string(body))
	}

	return nil
}

// DeleteKey deletes a key from a namespace
func (c *Client) DeleteKey(namespace, key string) error {
	url := fmt.Sprintf("%s/delete/%s/%s", c.baseURL, namespace, key)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("key not found")
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete key: %s", string(body))
	}

	return nil
}

// GetAll retrieves all keys and values
func (c *Client) GetAll() (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/getall", c.baseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get all keys: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}

// Exists checks if a key exists in a namespace
func (c *Client) Exists(namespace, key string) (bool, error) {
	_, err := c.GetKey(namespace, key)
	if err != nil {
		if err.Error() == "key not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Close closes the HTTP client
func (c *Client) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
}
