// Package urlmap provides utilities for constructing and manipulating URLs for WMS requests.
package urlmap

import (
	"fmt"
	"net/url"
)

// URLmap wraps a URL and provides methods to manipulate query parameters.
type URLmap struct {
	u *url.URL
}

// New creates a new URLmap from a URL string.
func New(u string) (*URLmap, error) {
	// Use url.Parse() to parse a string into a *url.URL type. If your URL is
	// already a url.URL type you can skip this step.
	url, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("parsing URL: %w", err)
	}
	urlObject := URLmap{
		u: url,
	}
	return &urlObject, nil
}

// AddParameter adds or replaces a query parameter in the URL.
func (u *URLmap) AddParameter(param string, value string) {
	// Use the Query() method to get the query string params as a url.Values map.
	values := u.u.Query()
	// Make the changes that you want using the Add(), Set() and Del() methods. If
	// you want to retrieve or check for a specific parameter you can use the Get()
	// and Has() methods respectively.
	values.Del(param)
	values.Add(param, value)
	u.u.RawQuery = values.Encode()
}

func (u *URLmap) String() string {
	return u.u.String()
}
