// Package getmap is a package to handle and expand the abilities of Web Map Services
package getmap

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/wroge/wgs84"

	"github.com/sgaunet/wms/getcap"
	"github.com/sgaunet/wms/urlmap"
)

// MaxPixel which can be downloaded with GetMap
var MaxPixel = 64000000

// Service is a struct which holds the values for the GetMap request
type Service struct {
	Capabilities getcap.Abilities
	url          *urlmap.URLmap
	Version      string
	Format       string
	Layers       []string
	Styles       []string
	EPSG         int
}

// New is the constructor which accepts optional parameters
func New(options ...func(*Service) error) (s *Service, err error) {
	s = &Service{}
	for _, o := range options {
		err = o(s)
		if err != nil {
			return
		}
	}
	return
}

// InvalidInput is a Error type for invalid inputs
type InvalidInput string

func (e InvalidInput) Error() string {
	return string(e)
}

// GetCapabilities puts random values from the GetCapabilities-Document into the Service
// URL and Version have to be set
// Is called within the New constructor
func (s *Service) GetCapabilities() (c getcap.Abilities, err error) {
	c, err = getcap.From(s.url, s.Version)
	if err != nil {
		return
	}
	ff := c.Formats
	ll := c.Layers
	bb := c.GetBBoxes()
	if len(ff) < 1 || len(ll) < 1 || len(bb) < 1 || c.Version == "" {
		err = InvalidInput("Invalid: Please check URL and Version")
		return
	}
	s.Version = c.Version
	s.Format = ff[0]
	s.Layers = []string{ll[0].Name}
	s.Styles = make([]string, len(s.Layers))
	if bb.GetBBox(4326).GetEPSG() == 4326 {
		s.EPSG = 4326
		return
	}
	b := bb[0].GetEPSG()
	for i := 1; b == 0; i++ {
		if len(bb) < i+1 {
			err = InvalidInput("Invalid: Please check URL and Version")
			return
		}
		b = bb[i].GetEPSG()
	}
	s.EPSG = b
	return
}

// SetURL is an optional Parameter for the constructor
func SetURL(url string) func(*Service) error {
	return func(s *Service) error {
		return s.SetURL(url)
	}
}

// SetURL adds a URL to a Service
func (s *Service) SetURL(url string) (err error) {
	s.url, err = urlmap.New(url)
	if err != nil {
		return
	}
	c, err := s.GetCapabilities()
	s.Version = c.Version
	s.Capabilities = c
	return
}

// AddVersion is an optional Parameter for the constructor
func AddVersion(version string) func(*Service) error {
	return func(s *Service) error {
		return s.AddVersion(version)
	}
}

// AddVersion adds a version to a Service
func (s *Service) AddVersion(version string) (err error) {
	s.Version = version
	c, err := s.GetCapabilities()
	s.Capabilities = c
	return
}

// InvalidValue is a Error type for invalid value inputs
type InvalidValue struct {
	Field       string
	Value       string
	ValidValues []string
}

func (e InvalidValue) Error() string {
	return fmt.Sprintf("Invalid %v: %v\nValid %vs: %v", e.Field, e.Value, e.Field, e.ValidValues)
}

// AddFormat is an optional Parameter for the constructor
func AddFormat(format string) func(*Service) error {
	return func(s *Service) error {
		return s.AddFormat(format)
	}
}

// AddFormat adds a format to a Service
func (s *Service) AddFormat(format string) (err error) {
	ff := s.Capabilities.Formats
	if !contains(ff, format) {
		err = InvalidValue{"Format", format, ff}
		return
	}
	s.Format = format
	return
}

// AddLayers is an optional Parameter for the constructor
func AddLayers(layers ...string) func(*Service) error {
	return func(s *Service) error {
		return s.AddLayers(layers...)
	}
}

// AddLayers adds layers to a Service
func (s *Service) AddLayers(layers ...string) (err error) {
	for _, l := range layers {
		cl := s.Capabilities.GetLayer(l)
		if cl.Name == "" {
			return InvalidValue{"Layer", l, s.Capabilities.GetLayerNames()}
		}
	}
	s.Layers = layers
	s.Styles = make([]string, len(s.Layers))
	return
}

// AddStyle is an optional Parameter for the constructor
func AddStyle(layer, style string) func(*Service) error {
	return func(s *Service) error {
		return s.AddStyle(layer, style)
	}
}

// AddStyle adds a style to a Service
func (s *Service) AddStyle(layer, style string) (err error) {
	ss := s.Capabilities.GetLayer(layer).Styles
	if !contains(ss, style) {
		return InvalidValue{"Style", style, ss}
	}
	if len(s.Styles) != len(s.Layers) {
		return errors.New("Adding Style failed")
	}
	for i, l := range s.Layers {
		if l == layer {
			s.Styles[i] = style
		}
	}
	return
}

// AddEPSG is an optional Parameter for the constructor
func AddEPSG(epsgCode int) func(*Service) error {
	return func(s *Service) error {
		return s.AddEPSG(epsgCode)
	}
}

// AddEPSG adds an EPSG code to a Service
func (s *Service) AddEPSG(epsgCode int) (err error) {
	epsgCap := s.Capabilities.GetBBoxes().GetEPSG()
	if len(epsgCap) == 0 {
		return errors.New("Adding EPSG failed")
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
		return InvalidValue{"EPSG", strconv.Itoa(epsgCode), eeStr}
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

// GetFileExt returns the file extension for various formats
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

// Option calculates width and height for a specific bounding box
type Option func(*Service, float64, float64, float64, float64) (width, height int, err error)

// ScaleDPIOption calculates width and height via scale and dpi
func ScaleDPIOption(scale, dpi int) Option {
	return func(s *Service, minx, miny, maxx, maxy float64) (width, height int, err error) {
		if scale == 0 || dpi == 0 {
			err = errors.New("Size must be set (Width, Height, Scale/DPI)")
			return
		}
		x1, y1, x2, y2 := utmCoord(minx, miny, maxx, maxy, s.EPSG)
		width = int(math.Round((x2 - x1) / float64(scale) * float64(dpi) * 25.4))
		height = int(math.Round((y2 - y1) / float64(scale) * float64(dpi) * 25.4))
		return
	}
}

// WidthHeightOption sets width and height
func WidthHeightOption(width, height int) Option {
	if width == 0 {
		return HeightOption(height)
	}
	if height == 0 {
		return WidthOption(width)
	}
	return func(s *Service, minx, miny, maxx, maxy float64) (widthN, heightN int, err error) {
		return width, height, nil
	}
}

// HeightOption calculates width via height and bounding box
func HeightOption(height int) Option {
	return func(s *Service, minx, miny, maxx, maxy float64) (widthN, heightN int, err error) {
		if height == 0 {
			err = errors.New("Width or Height must be set")
			return
		}
		x1, y1, x2, y2 := utmCoord(minx, miny, maxx, maxy, s.EPSG)
		width := int(math.Round((x2 - x1) / (y2 - y1) * float64(height)))
		return width, height, nil
	}
}

// WidthOption calculates height via width and bounding box
func WidthOption(width int) Option {
	return func(s *Service, minx, miny, maxx, maxy float64) (widthN, heightN int, err error) {
		if width == 0 {
			err = errors.New("Width or Height must be set")
			return
		}
		x1, y1, x2, y2 := utmCoord(minx, miny, maxx, maxy, s.EPSG)
		height := int(math.Round((y2 - y1) / (x2 - x1) * float64(width)))
		return width, height, nil
	}
}

func utmCoord(minx, miny, maxx, maxy float64, e int) (x1, y1, x2, y2 float64) {
	x1, y1, _ = wgs84.EPSG().Transform(e, 4326)(minx, miny, 0)
	x2, y2, _ = wgs84.EPSG().Transform(e, 4326)(maxx, maxy, 0)
	zone1 := math.Floor(x1/6) + 31
	zone2 := math.Floor(x2/6) + 31
	northern := true
	if y1 < 0 || y2 < 0 {
		northern = false
	}
	x1, y1, _ = wgs84.Transform(wgs84.LonLat(), wgs84.UTM((zone1+zone2)/2, northern))(x1, y1, 0)
	x2, y2, _ = wgs84.Transform(wgs84.LonLat(), wgs84.UTM((zone1+zone2)/2, northern))(x2, y2, 0)
	return
}

// GetMap returns a bytes.Reader which contains the image data and the width and height of the image
func (s *Service) GetMap(minx, miny, maxx, maxy float64, o Option) (r *bytes.Reader, width, height int, err error) {
	width, height, err = o(s, minx, miny, maxx, maxy)
	if err != nil {
		return
	}
	if width*height > MaxPixel {
		err = InvalidInput("Invalid: Image is too big: " + strconv.Itoa(width*height) + " Max Pixel: " + strconv.Itoa(MaxPixel))
		return
	}
	epsgCap := s.Capabilities.GetBBoxes().GetEPSG()
	if !containsInt(epsgCap, s.EPSG) {
		from := wgs84.EPSG().Code(s.EPSG)
		if from == nil {
			return nil, 0, 0, err
		}
		to := wgs84.EPSG().Code(epsgCap[0])
		if to == nil {
			return nil, 0, 0, err
		}
		minx, miny, _ = wgs84.Transform(from, to)(minx, miny, 0)
		maxx, maxy, _ = wgs84.Transform(from, to)(maxx, maxy, 0)
		s.EPSG = epsgCap[0]
	}
	bbox := s.Capabilities.GetBBox(s.EPSG)
	if minx < bbox.MinX || minx > bbox.MaxX || maxx < bbox.MinX || maxx > bbox.MaxX || miny < bbox.MinY || miny > bbox.MaxY || maxy < bbox.MinY || maxy > bbox.MaxY {
		err = InvalidInput("Invalid: BBox is out of bounds: " + fmt.Sprintf("%v,%v,%v,%v", minx, miny, maxx, maxy) + "\nValid BBox: " + bbox.String())
		return
	}
	// request := fmt.Sprintf("%v?SERVICE=WMS&REQUEST=GetMap&VERSION=%v&FORMAT=%v&LAYERS=%v&STYLES=%v", s.url.String(), s.Version, s.Format, strings.Join(s.Layers, ","), strings.Join(s.Styles, ","))
	getmapUrl, _ := urlmap.New(s.url.String())
	getmapUrl.AddParameter("SERVICE", "WMS")
	getmapUrl.AddParameter("REQUEST", "GetMap")
	getmapUrl.AddParameter("VERSION", s.Version)
	getmapUrl.AddParameter("FORMAT", s.Format)
	getmapUrl.AddParameter("LAYERS", strings.Join(s.Layers, ","))
	getmapUrl.AddParameter("STYLES", strings.Join(s.Styles, ","))

	if s.Version == "1.3.0" {
		// request += fmt.Sprintf("&CRS=EPSG:%v", s.EPSG)
		getmapUrl.AddParameter("CRS", fmt.Sprintf("EPSG:%d", s.EPSG))
	} else {
		getmapUrl.AddParameter("SRS", fmt.Sprintf("EPSG:%d", s.EPSG))
		// request += fmt.Sprintf("&SRS=EPSG:%v", s.EPSG)
	}
	if width != 0 {
		// ERR
	}
	if height != 0 {
		// ERR
	}
	getmapUrl.AddParameter("HEIGHT", fmt.Sprintf("%d", height))
	getmapUrl.AddParameter("WIDTH", fmt.Sprintf("%d", width))
	getmapUrl.AddParameter("BBOX", fmt.Sprintf("%.7f,%.7f,%.7f,%.7f", minx, miny, maxx, maxy))

	// request += fmt.Sprintf("&WIDTH=%v&HEIGHT=%v&BBOX=%.7f,%.7f,%.7f,%.7f", width, height, minx, miny, maxx, maxy)
	fmt.Println(getmapUrl.String())
	req, err := http.NewRequest("GET", getmapUrl.String(), nil)
	if err != nil {
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	r = bytes.NewReader(buf)
	return r, width, height, err
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
