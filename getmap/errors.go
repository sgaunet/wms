package getmap

import (
	"errors"
	"fmt"
)

// Define static errors.
var (
	ErrAddingStyleFailed = errors.New("adding Style failed")
	ErrAddingEPSGFailed  = errors.New("adding EPSG failed")
	ErrSizeMustBeSet     = errors.New("size must be set (Width, Height, Scale/DPI)")
	ErrWidthHeightReq    = errors.New("width or Height must be set")
	ErrInvalidEPSG       = errors.New("invalid EPSG code")
)

// InvalidSourceEPSGError represents an invalid source EPSG code error.
type InvalidSourceEPSGError struct {
	Code int
}

func (e InvalidSourceEPSGError) Error() string {
	return fmt.Sprintf("invalid source EPSG: %d", e.Code)
}

// InvalidTargetEPSGError represents an invalid target EPSG code error.
type InvalidTargetEPSGError struct {
	Code int
}

func (e InvalidTargetEPSGError) Error() string {
	return fmt.Sprintf("invalid target EPSG: %d", e.Code)
}