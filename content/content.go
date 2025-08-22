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

// From return data from a URL with Basic Auth.
func From(url *urlmap.URLmap) (*bytes.Reader, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
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
