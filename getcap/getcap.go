// Package getcap parses a GetCapabilities-Request
package getcap

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/sgaunet/wms/content"
	"github.com/sgaunet/wms/urlmap"
)

// Formats of a GetCapabilities-Request.
type Formats []string

// Styles of a GetCapabilities-Request.
type Styles []string

// Layers of a GetCapabilities-Request.
type Layers []Layer

// BBoxes of a GetCapabilities-Request.
type BBoxes []BBox

// Abilities (Capabilities) of a GetCapabilities-Request.
// This struct supports both WMS 1.1.x and 1.3.0 formats.
type Abilities struct {
	Version  string
	Name     string
	Title    string
	Abstract string
	Formats  Formats
	Layers   Layers
	BBoxes   BBoxes
}

// abilities11x represents WMS 1.1.x format with WMT_MS_Capabilities root element.
type abilities11x struct {
	XMLName  xml.Name `xml:"WMT_MS_Capabilities"`
	Version  string   `xml:"version,attr"`
	Name     string   `xml:"Service>Name"`
	Title    string   `xml:"Service>Title"`
	Abstract string   `xml:"Service>Abstract"`
	Formats  Formats  `xml:"Capability>Request>GetMap>Format"`
	Layers   Layers   `xml:"Capability>Layer>Layer"`
	BBoxes   BBoxes   `xml:"Capability>Layer>BoundingBox"`
}

// abilities130 represents WMS 1.3.0 format with WMS_Capabilities root element.
type abilities130 struct {
	XMLName  xml.Name `xml:"WMS_Capabilities"`
	Version  string   `xml:"version,attr"`
	Name     string   `xml:"Service>Name"`
	Title    string   `xml:"Service>Title"`
	Abstract string   `xml:"Service>Abstract"`
	Formats  Formats  `xml:"Capability>Request>GetMap>Format"`
	Layers   Layers   `xml:"Capability>Layer>Layer"`
	BBoxes   BBoxes   `xml:"Capability>Layer>BoundingBox"`
}

// Layer of a GetCapabilities-Request.
type Layer struct {
	Name   string `xml:"Name"`
	Styles Styles `xml:"Style>Name"`
	BBoxes BBoxes `xml:"BoundingBox"`
}

// BBox of a GetCapabilities-Request.
type BBox struct {
	SRS  string  `xml:"SRS,attr"`
	CRS  string  `xml:"CRS,attr"`
	MinX float64 `xml:"minx,attr"`
	MinY float64 `xml:"miny,attr"`
	MaxX float64 `xml:"maxx,attr"`
	MaxY float64 `xml:"maxy,attr"`
}

// From Capabilities of a WMS service.
func From(url *urlmap.URLmap, version string) (Abilities, error) {
	requestURL, err := urlmap.New(url.String())
	if err != nil {
		return Abilities{}, fmt.Errorf("creating request URL: %w", err)
	}
	// request := url + "?SERVICE=WMS&REQUEST=GetCapabilities"
	requestURL.AddParameter("SERVICE", "WMS")
	requestURL.AddParameter("REQUEST", "GetCapabilities")
	if version != "" {
		requestURL.AddParameter("VERSION", version)
	}

	// fmt.Println(requestURL.String())
	reader, err := content.From(requestURL)
	// reader, err := content.From(request.Request.String(), request.User.Username(), request.User.Password())
	if err != nil {
		return Abilities{}, fmt.Errorf("fetching capabilities: %w", err)
	}
	c, err := Read(reader)
	if err != nil {
		return Abilities{}, fmt.Errorf("parsing capabilities: %w", err)
	}
	return c, nil
}

// Read Capabilities from a GetCapabilities-Document.
// Supports both WMS 1.1.x (WMT_MS_Capabilities) and 1.3.0 (WMS_Capabilities) formats.
func Read(data io.Reader) (Abilities, error) {
	// Read the data into a buffer so we can try parsing multiple formats
	buf, err := io.ReadAll(data)
	if err != nil {
		return Abilities{}, fmt.Errorf("reading data: %w", err)
	}

	// Try WMS 1.3.0 format first (WMS_Capabilities)
	var c130 abilities130
	decoder := xml.NewDecoder(bytes.NewReader(buf))
	err = decoder.Decode(&c130)
	if err == nil {
		// Successfully parsed as WMS 1.3.0
		return Abilities{
			Version:  c130.Version,
			Name:     c130.Name,
			Title:    c130.Title,
			Abstract: c130.Abstract,
			Formats:  c130.Formats,
			Layers:   c130.Layers,
			BBoxes:   c130.BBoxes,
		}, nil
	}

	// Try WMS 1.1.x format (WMT_MS_Capabilities)
	var c11x abilities11x
	decoder = xml.NewDecoder(bytes.NewReader(buf))
	err = decoder.Decode(&c11x)
	if err == nil {
		// Successfully parsed as WMS 1.1.x
		return Abilities{
			Version:  c11x.Version,
			Name:     c11x.Name,
			Title:    c11x.Title,
			Abstract: c11x.Abstract,
			Formats:  c11x.Formats,
			Layers:   c11x.Layers,
			BBoxes:   c11x.BBoxes,
		}, nil
	}

	return Abilities{}, fmt.Errorf("decoding XML: unable to parse as WMS 1.3.0 or 1.1.x format: %w", err)
}
