// Package getcap parses a GetCapabilities-Request
package getcap

import (
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
type Abilities struct {
	XMLName  xml.Name `xml:"WMT_MS_Capabilities"`
	Version  string  `xml:"version,attr"`
	Name     string  `xml:"Service>Name"`
	Title    string  `xml:"Service>Title"`
	Abstract string  `xml:"Service>Abstract"`
	Formats  Formats `xml:"Capability>Request>GetMap>Format"`
	Layers   Layers  `xml:"Capability>Layer>Layer"`
	BBoxes   BBoxes  `xml:"Capability>Layer>BoundingBox"`
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
func Read(data io.Reader) (Abilities, error) {
	var c Abilities
	decoder := xml.NewDecoder(data)
	err := decoder.Decode(&c)
	if err != nil {
		return c, fmt.Errorf("decoding XML: %w", err)
	}
	return c, nil
}
