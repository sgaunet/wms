package content

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sgaunet/wms/urlmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromWithBasicAuth(t *testing.T) {
	// Expected credentials
	expectedUser := "testuser"
	expectedPass := "testpass"
	expectedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(expectedUser+":"+expectedPass))

	// Create test server that checks for auth header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != expectedAuth {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("authenticated response"))
	}))
	defer server.Close()

	// Create URL map
	urlMap, err := urlmap.New(server.URL)
	require.NoError(t, err, "Failed to create URL map")

	// Test with authentication
	reader, err := From(urlMap, WithBasicAuth(expectedUser, expectedPass))
	require.NoError(t, err, "Request with auth should succeed")
	require.NotNil(t, reader, "Reader should not be nil")

	// Read response
	data := make([]byte, 100)
	n, _ := reader.Read(data)
	assert.Equal(t, "authenticated response", string(data[:n]))
}

func TestFromWithoutAuth(t *testing.T) {
	// Create test server that doesn't require auth
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		assert.Empty(t, auth, "No auth header should be sent")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("public response"))
	}))
	defer server.Close()

	// Create URL map
	urlMap, err := urlmap.New(server.URL)
	require.NoError(t, err, "Failed to create URL map")

	// Test without authentication
	reader, err := From(urlMap)
	require.NoError(t, err, "Request without auth should succeed")
	require.NotNil(t, reader, "Reader should not be nil")

	// Read response
	data := make([]byte, 100)
	n, _ := reader.Read(data)
	assert.Equal(t, "public response", string(data[:n]))
}

func TestFromWithEmptyUsername(t *testing.T) {
	// Create test server that checks no auth is sent when username is empty
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		assert.Empty(t, auth, "No auth header should be sent when username is empty")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("response"))
	}))
	defer server.Close()

	// Create URL map
	urlMap, err := urlmap.New(server.URL)
	require.NoError(t, err, "Failed to create URL map")

	// Test with empty username (should not send auth header)
	reader, err := From(urlMap, WithBasicAuth("", "password"))
	require.NoError(t, err, "Request should succeed")
	require.NotNil(t, reader, "Reader should not be nil")
}

func TestFromUnauthorized(t *testing.T) {
	// Create test server that always returns 401
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	// Create URL map
	urlMap, err := urlmap.New(server.URL)
	require.NoError(t, err, "Failed to create URL map")

	// Test with wrong credentials
	_, err = From(urlMap, WithBasicAuth("wrong", "credentials"))
	require.Error(t, err, "Should return error for 401 response")
	assert.IsType(t, StatusCodeError{}, err, "Error should be StatusCodeError")
}

func TestWithBasicAuthOption(t *testing.T) {
	// Test that the option correctly sets credentials
	cfg := &Config{}
	opt := WithBasicAuth("user1", "pass1")
	opt(cfg)

	assert.Equal(t, "user1", cfg.Username)
	assert.Equal(t, "pass1", cfg.Password)
}
