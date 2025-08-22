package cmd

import "errors"

// Define static errors.
var (
	ErrEmptyURL       = errors.New("url is empty")
	ErrInvalidBBox    = errors.New("invalid BBox")
	ErrInvalidBBoxReq = errors.New("invalid: Add BBox")
)