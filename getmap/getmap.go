// Package getmap is a package to handle and expand the abilities of Web Map Services
package getmap

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/wroge/wgs84"

	"github.com/sgaunet/wms/getcap"
	"github.com/sgaunet/wms/urlmap"
)

// MaxPixel which can be downloaded with GetMap.
var MaxPixel = 64000000

const (
	wgs84EPSG        = 4326
	millimetersPerInch = 25.4
	utmZoneDivisor   = 6
	utmZoneOffset    = 31
	averageDivisor   = 2
)

// Service is a struct which holds the values for the GetMap request.
type Service struct {
	Capabilities getcap.Abilities
	url          *urlmap.URLmap
	Version      string
	Format       string
	Layers       []string
	Styles       []string
	EPSG         int
}

// New is the constructor which accepts optional parameters.
func New(options ...func(*Service) error) (*Service, error) {
	s := &Service{}
	for _, o := range options {
		if err := o(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

// InvalidInputError is a Error type for invalid inputs.
type InvalidInputError string

func (e InvalidInputError) Error() string {
	return string(e)
}

// GetCapabilities puts random values from the GetCapabilities-Document into the Service.
// URL and Version have to be set.
// Is called within the New constructor.
func (s *Service) GetCapabilities() (getcap.Abilities, error) {
	c, err := getcap.From(s.url, s.Version)
	if err != nil {
		return getcap.Abilities{}, fmt.Errorf("getting capabilities: %w", err)
	}
	ff := c.Formats
	ll := c.Layers
	bb := c.GetBBoxes()
	if len(ff) < 1 || len(ll) < 1 || len(bb) < 1 || c.Version == "" {
		return getcap.Abilities{}, InvalidInputError("Invalid: Please check URL and Version")
	}
	s.Version = c.Version
	s.Format = ff[0]
	s.Layers = []string{ll[0].Name}
	s.Styles = make([]string, len(s.Layers))
	if bb.GetBBox(wgs84EPSG).GetEPSG() == wgs84EPSG {
		s.EPSG = 4326
		return c, nil
	}
	b := bb[0].GetEPSG()
	for i := 1; b == 0; i++ {
		if len(bb) < i+1 {
			return getcap.Abilities{}, InvalidInputError("Invalid: Please check URL and Version")
		}
		b = bb[i].GetEPSG()
	}
	s.EPSG = b
	return c, nil
}

// SetURL is an optional Parameter for the constructor.
func SetURL(url string) func(*Service) error {
	return func(s *Service) error {
		return s.SetURL(url)
	}
}

// SetURL adds a URL to a Service.
func (s *Service) SetURL(url string) error {
	var err error
	s.url, err = urlmap.New(url)
	if err != nil {
		return fmt.Errorf("creating URL map: %w", err)
	}
	c, err := s.GetCapabilities()
	if err != nil {
		return err
	}
	s.Version = c.Version
	s.Capabilities = c
	return nil
}

// AddVersion is an optional Parameter for the constructor.
func AddVersion(version string) func(*Service) error {
	return func(s *Service) error {
		return s.AddVersion(version)
	}
}

// AddVersion adds a version to a Service.
func (s *Service) AddVersion(version string) error {
	s.Version = version
	c, err := s.GetCapabilities()
	if err != nil {
		return err
	}
	s.Capabilities = c
	return nil
}

// InvalidValueError is a Error type for invalid value inputs.
type InvalidValueError struct {
	Field       string
	Value       string
	ValidValues []string
}

func (e InvalidValueError) Error() string {
	return fmt.Sprintf("Invalid %v: %v\nValid %vs: %v", e.Field, e.Value, e.Field, e.ValidValues)
}

// AddFormat is an optional Parameter for the constructor.
func AddFormat(format string) func(*Service) error {
	return func(s *Service) error {
		return s.AddFormat(format)
	}
}

// AddFormat adds a format to a Service.
func (s *Service) AddFormat(format string) error {
	ff := s.Capabilities.Formats
	if !contains(ff, format) {
		return InvalidValueError{"Format", format, ff}
	}
	s.Format = format
	return nil
}

// AddLayers is an optional Parameter for the constructor.
func AddLayers(layers ...string) func(*Service) error {
	return func(s *Service) error {
		return s.AddLayers(layers...)
	}
}

// AddLayers adds layers to a Service.
func (s *Service) AddLayers(layers ...string) error {
	for _, l := range layers {
		cl := s.Capabilities.GetLayer(l)
		if cl.Name == "" {
			return InvalidValueError{"Layer", l, s.Capabilities.GetLayerNames()}
		}
	}
	s.Layers = layers
	s.Styles = make([]string, len(s.Layers))
	return nil
}

// AddStyle is an optional Parameter for the constructor.
func AddStyle(layer, style string) func(*Service) error {
	return func(s *Service) error {
		return s.AddStyle(layer, style)
	}
}

// AddStyle adds a style to a Service.
func (s *Service) AddStyle(layer, style string) error {
	ss := s.Capabilities.GetLayer(layer).Styles
	if !contains(ss, style) {
		return InvalidValueError{"Style", style, ss}
	}
	if len(s.Styles) != len(s.Layers) {
		return ErrAddingStyleFailed
	}
	for i, l := range s.Layers {
		if l == layer {
			s.Styles[i] = style
		}
	}
	return nil
}

// AddEPSG is an optional Parameter for the constructor.
func AddEPSG(epsgCode int) func(*Service) error {
	return func(s *Service) error {
		return s.AddEPSG(epsgCode)
	}
}

// AddEPSG adds an EPSG code to a Service.
func (s *Service) AddEPSG(epsgCode int) error {
	epsgCap := s.Capabilities.GetBBoxes().GetEPSG()
	if len(epsgCap) == 0 {
		return ErrAddingEPSGFailed
	}
	for _, e := range wgs84.EPSG().Codes() {
		redundant := false
		for _, eeC := range epsgCap {
			if eeC == e {
				redundant = true
			}
		}
		if !redundant {
			epsgCap = append(epsgCap, e)
		}
	}
	if !containsInt(epsgCap, epsgCode) {
		eeStr := []string{}
		for _, ee := range epsgCap {
			eeStr = append(eeStr, strconv.Itoa(ee))
		}
		return InvalidValueError{"EPSG", strconv.Itoa(epsgCode), eeStr}
	}
	s.EPSG = epsgCode
	return nil
}

// // Validate validates a Service which is not made by the constructor or methods
// func (s *Service) Validate() (err error) {
// 	n := &Service{}
// 	n.URL = s.URL
// 	err = n.AddVersion(s.Version)
// 	if err != nil {
// 		return
// 	}
// 	err = n.AddFormat(s.Format)
// 	if err != nil {
// 		return
// 	}
// 	err = n.AddLayers(s.Layers...)
// 	if err != nil {
// 		return
// 	}
// 	for i, st := range s.Styles {
// 		if st != "" {
// 			err = n.AddStyle(s.Layers[i], st)
// 			if err != nil {
// 				return
// 			}
// 		}
// 	}
// 	err = n.AddEPSG(s.EPSG)
// 	if err != nil {
// 		return
// 	}
// 	return
// }

func (s *Service) String() string {
	return fmt.Sprintf(`URL: %v
Version: %v
Format: %v
Layers: %v
Styles: %v
EPSG: %v`, s.url.String(), s.Version, s.Format, s.Layers, s.Styles, s.EPSG)
}

// GetFileExt returns the file extension for various formats.
func (s *Service) GetFileExt() string {
	if s.Format == "image/png" {
		return "png"
	}
	if s.Format == "image/jpeg" {
		return "jpeg"
	}
	if s.Format == "image/gif" {
		return "gif"
	}
	if s.Format == "image/tiff" {
		return "tiff"
	}
	return ""
}

// Option calculates width and height for a specific bounding box.
type Option func(*Service, float64, float64, float64, float64) (int, int, error)

// ScaleDPIOption calculates width and height via scale and dpi.
func ScaleDPIOption(scale, dpi int) Option {
	return func(s *Service, minx, miny, maxx, maxy float64) (int, int, error) {
		if scale == 0 || dpi == 0 {
			return 0, 0, ErrSizeMustBeSet
		}
		x1, y1, x2, y2 := utmCoord(minx, miny, maxx, maxy, s.EPSG)
		width := int(math.Round((x2 - x1) / float64(scale) * float64(dpi) * millimetersPerInch))
		height := int(math.Round((y2 - y1) / float64(scale) * float64(dpi) * millimetersPerInch))
		return width, height, nil
	}
}

// WidthHeightOption sets width and height.
func WidthHeightOption(width, height int) Option {
	if width == 0 {
		return HeightOption(height)
	}
	if height == 0 {
		return WidthOption(width)
	}
	return func(_ *Service, _, _, _, _ float64) (int, int, error) {
		return width, height, nil
	}
}

// HeightOption calculates width via height and bounding box.
func HeightOption(height int) Option {
	return func(s *Service, minx, miny, maxx, maxy float64) (int, int, error) {
		if height == 0 {
			return 0, 0, ErrWidthHeightReq
		}
		x1, y1, x2, y2 := utmCoord(minx, miny, maxx, maxy, s.EPSG)
		width := int(math.Round((x2 - x1) / (y2 - y1) * float64(height)))
		return width, height, nil
	}
}

// WidthOption calculates height via width and bounding box.
func WidthOption(width int) Option {
	return func(s *Service, minx, miny, maxx, maxy float64) (int, int, error) {
		if width == 0 {
			return 0, 0, ErrWidthHeightReq
		}
		x1, y1, x2, y2 := utmCoord(minx, miny, maxx, maxy, s.EPSG)
		height := int(math.Round((y2 - y1) / (x2 - x1) * float64(width)))
		return width, height, nil
	}
}

func utmCoord(minx, miny, maxx, maxy float64, e int) (float64, float64, float64, float64) {
	x1, y1, _ := wgs84.EPSG().Transform(e, wgs84EPSG)(minx, miny, 0)
	x2, y2, _ := wgs84.EPSG().Transform(e, wgs84EPSG)(maxx, maxy, 0)
	zone1 := math.Floor(x1/utmZoneDivisor) + utmZoneOffset
	zone2 := math.Floor(x2/utmZoneDivisor) + utmZoneOffset
	northern := !(y1 < 0 || y2 < 0)
	x1, y1, _ = wgs84.Transform(wgs84.LonLat(), wgs84.UTM((zone1+zone2)/averageDivisor, northern))(x1, y1, 0)
	x2, y2, _ = wgs84.Transform(wgs84.LonLat(), wgs84.UTM((zone1+zone2)/averageDivisor, northern))(x2, y2, 0)
	return x1, y1, x2, y2
}

// validateImageSize checks if the requested image size is valid.
func validateImageSize(width, height int) error {
	if width*height > MaxPixel {
		return InvalidInputError("Invalid: Image is too big: " + strconv.Itoa(width*height) +
			" Max Pixel: " + strconv.Itoa(MaxPixel))
	}
	return nil
}


// GetMap returns a bytes.Reader which contains the image data and the width and height of the image.
func (s *Service) GetMap(minx, miny, maxx, maxy float64, o Option) (*bytes.Reader, int, int, error) {
	width, height, err := o(s, minx, miny, maxx, maxy)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("calculating dimensions: %w", err)
	}

	if err := validateImageSize(width, height); err != nil {
		return nil, 0, 0, err
	}

	minx, miny, maxx, maxy, err = s.transformCoordinates(minx, miny, maxx, maxy)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("transforming coordinates: %w", err)
	}

	if err := s.validateBBox(minx, miny, maxx, maxy); err != nil {
		return nil, 0, 0, err
	}

	getmapURL, err := s.buildGetMapURL(minx, miny, maxx, maxy, width, height)
	if err != nil {
		return nil, 0, 0, err
	}

	fmt.Println(getmapURL.String())
	r, err := executeGetMapRequest(getmapURL.String())
	if err != nil {
		return nil, 0, 0, err
	}
	return r, width, height, nil
}

// transformCoordinates transforms coordinates to a supported EPSG if needed.
func (s *Service) transformCoordinates(minx, miny, maxx, maxy float64) (float64, float64, float64, float64, error) {
	epsgCap := s.Capabilities.GetBBoxes().GetEPSG()
	if containsInt(epsgCap, s.EPSG) {
		return minx, miny, maxx, maxy, nil
	}

	from := wgs84.EPSG().Code(s.EPSG)
	if from == nil {
		return 0, 0, 0, 0, InvalidSourceEPSGError{Code: s.EPSG}
	}
	to := wgs84.EPSG().Code(epsgCap[0])
	if to == nil {
		return 0, 0, 0, 0, InvalidTargetEPSGError{Code: epsgCap[0]}
	}
	minx, miny, _ = wgs84.Transform(from, to)(minx, miny, 0)
	maxx, maxy, _ = wgs84.Transform(from, to)(maxx, maxy, 0)
	s.EPSG = epsgCap[0]
	return minx, miny, maxx, maxy, nil
}

// validateBBox checks if the bounding box is within valid bounds.
func (s *Service) validateBBox(minx, miny, maxx, maxy float64) error {
	bbox := s.Capabilities.GetBBox(s.EPSG)
	if minx < bbox.MinX || minx > bbox.MaxX || maxx < bbox.MinX || maxx > bbox.MaxX ||
		miny < bbox.MinY || miny > bbox.MaxY || maxy < bbox.MinY || maxy > bbox.MaxY {
		return InvalidInputError("Invalid: BBox is out of bounds: " +
			fmt.Sprintf("%v,%v,%v,%v", minx, miny, maxx, maxy) + "\nValid BBox: " + bbox.String())
	}
	return nil
}

// buildGetMapURL creates the GetMap request URL.
func (s *Service) buildGetMapURL(minx, miny, maxx, maxy float64, width, height int) (*urlmap.URLmap, error) {
	getmapURL, err := urlmap.New(s.url.String())
	if err != nil {
		return nil, fmt.Errorf("creating GetMap URL: %w", err)
	}
	getmapURL.AddParameter("SERVICE", "WMS")
	getmapURL.AddParameter("REQUEST", "GetMap")
	getmapURL.AddParameter("VERSION", s.Version)
	getmapURL.AddParameter("FORMAT", s.Format)
	getmapURL.AddParameter("LAYERS", strings.Join(s.Layers, ","))
	getmapURL.AddParameter("STYLES", strings.Join(s.Styles, ","))

	if s.Version == "1.3.0" {
		getmapURL.AddParameter("CRS", fmt.Sprintf("EPSG:%d", s.EPSG))
	} else {
		getmapURL.AddParameter("SRS", fmt.Sprintf("EPSG:%d", s.EPSG))
	}
	getmapURL.AddParameter("HEIGHT", strconv.Itoa(height))
	getmapURL.AddParameter("WIDTH", strconv.Itoa(width))
	getmapURL.AddParameter("BBOX", fmt.Sprintf("%.7f,%.7f,%.7f,%.7f", minx, miny, maxx, maxy))
	return getmapURL, nil
}

// executeGetMapRequest performs the HTTP request to get the map.
func executeGetMapRequest(url string) (*bytes.Reader, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating HTTP request: %w", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing HTTP request: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = fmt.Errorf("closing response body: %w", cerr)
		}
	}()
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}
	return bytes.NewReader(buf), nil
}

func contains(xx []string, y string) bool {
	for _, x := range xx {
		if x == y {
			return true
		}
	}
	return false
}

func containsInt(xx []int, y int) bool {
	for _, x := range xx {
		if x == y {
			return true
		}
	}
	return false
}
