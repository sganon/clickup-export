package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const ClickupAPIBaseURL = "https://api.clickup.com/api/v3"

type ClickupClient struct {
	token string
}

func NewClickupClient(token string) *ClickupClient {
	return &ClickupClient{
		token: token,
	}
}

func (c *ClickupClient) request(ctx context.Context, method, path string, params url.Values, body any, dest any) error {
	endpoint := fmt.Sprintf("%s%s", ClickupAPIBaseURL, path)
	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("invalid endpoint %s: %w", endpoint, err)
	}
	u.RawQuery = params.Encode()

	var b io.Reader
	if body != nil {
		rawBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("error marshalling request body: %w", err)
		}
		b = bytes.NewBuffer(rawBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), b)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", c.token)
	if req.Method == "POST" || req.Method == "PUT" {
		req.Header.Set("Content-Type", "application/json")
	}

	client := http.Client{
		Timeout: time.Minute,
	}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

	if dest != nil {
		if err := json.NewDecoder(res.Body).Decode(dest); err != nil {
			return fmt.Errorf("error decoding response body: %w", err)
		}
	}

	return nil
}
