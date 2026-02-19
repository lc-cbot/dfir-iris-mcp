package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func New(baseURL, apiKey string) *Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if os.Getenv("DFIR_IRIS_TLS_SKIP_VERIFY") != "" {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{Transport: transport},
	}
}

type envelope struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type APIError struct {
	StatusCode int
	Status     string
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("DFIR-IRIS API error (HTTP %d): %s - %s", e.StatusCode, e.Status, e.Message)
}

func (c *Client) Get(ctx context.Context, path string, query map[string]string) (json.RawMessage, error) {
	return c.do(ctx, http.MethodGet, path, query, nil)
}

func (c *Client) Post(ctx context.Context, path string, query map[string]string, body interface{}) (json.RawMessage, error) {
	return c.do(ctx, http.MethodPost, path, query, body)
}

func (c *Client) do(ctx context.Context, method, path string, query map[string]string, body interface{}) (json.RawMessage, error) {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	if len(query) > 0 {
		q := u.Query()
		for k, v := range query {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var env envelope
	if err := json.Unmarshal(respBody, &env); err != nil {
		if resp.StatusCode >= 400 {
			return nil, &APIError{StatusCode: resp.StatusCode, Message: string(respBody)}
		}
		return respBody, nil
	}

	if env.Status != "success" {
		msg := env.Message
		if len(env.Data) > 0 && string(env.Data) != "null" {
			msg = msg + " - " + string(env.Data)
		}
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Status:     env.Status,
			Message:    msg,
		}
	}

	return env.Data, nil
}
