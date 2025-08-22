package content

import "fmt"

// StatusCodeError represents an HTTP status code error.
type StatusCodeError struct {
	Code int
}

func (e StatusCodeError) Error() string {
	return fmt.Sprintf("status Code Error %v", e.Code)
}