// Package content supports getcap & getmap packages
package content

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/sgaunet/wms/urlmap"
)

// Option is a functional option for configuring HTTP requests.
type Option func(*Config)

// Config holds the configuration for HTTP requests.
type Config struct {
	Username string
	Password string
}

// WithBasicAuth returns an Option that sets HTTP Basic Authentication credentials.
func WithBasicAuth(username, password string) Option {
	return func(c *Config) {
		c.Username = username
		c.Password = password
	}
}

// From fetches data from a URL with optional configuration.
// Options can include authentication credentials via WithBasicAuth.
func From(url *urlmap.URLmap, opts ...Option) (*bytes.Reader, error) {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}

	client := &http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Apply Basic Authentication if credentials are provided
	if cfg.Username != "" {
		req.SetBasicAuth(cfg.Username, cfg.Password)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("closing response body: %w", cerr)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, StatusCodeError{Code: resp.StatusCode}
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}
	return bytes.NewReader(data), nil
}
