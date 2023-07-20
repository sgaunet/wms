package getmap

import (
	"strings"
	"testing"

	"github.com/sgaunet/wms/getcap"
	"github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {
	assert := assert.New(t)

	r := strings.NewReader(test)
	c, err := getcap.Read(r)
	assert.Nil(err, "Should be nil")
	assert.Equal(c.GetLayer("OSM-Overlay-WMS").Name, "OSM-Overlay-WMS", "Should be equal")
	assert.Equal(c.GetBBox(4326).GetEPSG(), 4326, "Should be equal")

	s := &Service{}
	s.Capabilities = c
	err = s.AddFormat("image/png")
	assert.Nil(err, "Should be nil")
	err = s.AddLayers("OSM-WMS")
	assert.Nil(err, "Should be nil")
	err = s.AddStyle("OSM-WMS", "default")
	assert.Nil(err, "Should be nil")
	err = s.AddEPSG(4326)
	assert.Nil(err, "Should be nil")
}

var test = `<?xml version="1.0"?>
<!DOCTYPE WMT_MS_Capabilities SYSTEM "http://schemas.opengis.net/wms/1.1.1/WMS_MS_Capabilities.dtd"
 [
 <!ELEMENT VendorSpecificCapabilities EMPTY>
 ]>  <!-- end of DOCTYPE declaration -->
<WMT_MS_Capabilities version="1.1.1">
<Service>
  <Name>OGC:WMS</Name>
  <Title>OpenStreetMap WMS</Title>
  <Abstract>OpenStreetMap WMS, bereitgestellt durch terrestris GmbH und Co. KG. Beschleunigt mit MapProxy (http://mapproxy.org/)</Abstract>
  <OnlineResource xmlns:xlink="http://www.w3.org/1999/xlink" xlink:href="http://www.terrestris.de"/>
  <ContactInformation>
      <ContactPersonPrimary>
        <ContactPerson>Johannes Weskamm</ContactPerson>
        <ContactOrganization>terrestris GmbH und Co. KG</ContactOrganization>
      </ContactPersonPrimary>
      <ContactPosition>Technical Director</ContactPosition>
      <ContactAddress>
        <AddressType>postal</AddressType>
        <Address>Kölnstr. 99</Address>
        <City>Bonn</City>
        <StateOrProvince></StateOrProvince>
        <PostCode>53111</PostCode>
        <Country>Germany</Country>
      </ContactAddress>
      <ContactVoiceTelephone>+49(0)228 962 899 51</ContactVoiceTelephone>
      <ContactFacsimileTelephone>+49(0)228 962 899 57</ContactFacsimileTelephone>
      <ContactElectronicMailAddress>info@terrestris.de</ContactElectronicMailAddress>
  </ContactInformation>
  <Fees>None</Fees>
  <AccessConstraints>(c) OpenStreetMap contributors (http://www.openstreetmap.org/copyright) (c) OpenStreetMap Data (http://openstreetmapdata.com) (c) Natural Earth Data (http://www.naturalearthdata.com) (c) ASTER GDEM 30m (https://asterweb.jpl.nasa.gov/gdem.asp) (c) SRTM 450m by ViewfinderPanoramas (http://viewfinderpanoramas.org/) (c) Great Lakes Bathymetry by NGDC (http://www.ngdc.noaa.gov/mgg/greatlakes/) (c) SRTM 30m by NASA EOSDIS Land Processes Distributed Active Archive Center (LP DAAC, https://lpdaac.usgs.gov/)</AccessConstraints>
</Service>
<Capability>
  <Request>
    <GetCapabilities>
      <Format>application/vnd.ogc.wms_xml</Format>
      <DCPType>
        <HTTP>
          <Get><OnlineResource xmlns:xlink="http://www.w3.org/1999/xlink" xlink:href="http://ows.terrestris.de/osm/service?"/></Get>
        </HTTP>
      </DCPType>
    </GetCapabilities>
    <GetMap>
        <Format>image/jpeg</Format>
        <Format>image/png</Format>
      <DCPType>
        <HTTP>
          <Get><OnlineResource xmlns:xlink="http://www.w3.org/1999/xlink" xlink:href="http://ows.terrestris.de/osm/service?"/></Get>
        </HTTP>
      </DCPType>
    </GetMap>
    <GetFeatureInfo>
      <Format>text/plain</Format>
      <Format>text/html</Format>
      <Format>application/vnd.ogc.gml</Format>
      <DCPType>
        <HTTP>
          <Get><OnlineResource xmlns:xlink="http://www.w3.org/1999/xlink" xlink:href="http://ows.terrestris.de/osm/service?"/></Get>
        </HTTP>
      </DCPType>
    </GetFeatureInfo>
    <GetLegendGraphic>
        <Format>image/jpeg</Format>
        <Format>image/png</Format>
        <DCPType>
            <HTTP>
                <Get><OnlineResource xmlns:xlink="http://www.w3.org/1999/xlink" xlink:href="http://ows.terrestris.de/osm/service?"/></Get>
            </HTTP>
        </DCPType>
    </GetLegendGraphic>
  </Request>
  <Exception>
    <Format>application/vnd.ogc.se_xml</Format>
    <Format>application/vnd.ogc.se_inimage</Format>
    <Format>application/vnd.ogc.se_blank</Format>
  </Exception>
  <Layer queryable="1">
    <Title>OpenStreetMap WMS</Title>
    <SRS>EPSG:900913</SRS>
    <SRS>EPSG:3857</SRS>
    <SRS>EPSG:25832</SRS>
    <SRS>EPSG:25833</SRS>
    <SRS>EPSG:29192</SRS>
    <SRS>EPSG:29193</SRS>
    <SRS>EPSG:31466</SRS>
    <SRS>EPSG:31467</SRS>
    <SRS>EPSG:31468</SRS>
    <SRS>EPSG:32648</SRS>
    <SRS>EPSG:4326</SRS>
    <SRS>EPSG:4674</SRS>
    <SRS>EPSG:3068</SRS>
    <SRS>EPSG:3034</SRS>
    <SRS>EPSG:3035</SRS>
    <SRS>EPSG:2100</SRS>
    <SRS>EPSG:31463</SRS>
    <SRS>EPSG:4258</SRS>
    <SRS>EPSG:4839</SRS>
    <SRS>EPSG:2180</SRS>
    <SRS>EPSG:21781</SRS>
    <SRS>EPSG:2056</SRS>
    <LatLonBoundingBox minx="-180" miny="-88" maxx="180" maxy="88" />
    <BoundingBox SRS="EPSG:900913" minx="-20037508.3428" miny="-25819498.5135" maxx="20037508.3428" maxy="25819498.5135" />
    <BoundingBox SRS="EPSG:4326" minx="-180" miny="-88" maxx="180" maxy="88" />
    <BoundingBox SRS="EPSG:3857" minx="-20037508.3428" miny="-25819498.5135" maxx="20037508.3428" maxy="25819498.5135" />
    <Layer queryable="1">
      <Name>OSM-WMS</Name>
      <Title>OpenStreetMap WMS - by terrestris</Title>
      <LatLonBoundingBox minx="-180" miny="-88" maxx="180" maxy="88" />
      <BoundingBox SRS="EPSG:900913" minx="-20037508.3428" miny="-25819498.5135" maxx="20037508.3428" maxy="25819498.5135" />
      <BoundingBox SRS="EPSG:4326" minx="-180" miny="-88" maxx="180" maxy="88" />
      <BoundingBox SRS="EPSG:3857" minx="-20037508.3428" miny="-25819498.5135" maxx="20037508.3428" maxy="25819498.5135" />
      <Style>
          <Name>default</Name>
          <Title>default</Title>
          <LegendURL width="155" height="344">
              <Format>image/png</Format>
              <OnlineResource xmlns:xlink="http://www.w3.org/1999/xlink" xlink:type="simple" xlink:href="http://ows.terrestris.de/osm/service?styles=&amp;layer=OSM-WMS&amp;service=WMS&amp;format=image%2Fpng&amp;sld_version=1.1.0&amp;request=GetLegendGraphic&amp;version=1.1.1"/>
          </LegendURL>
      </Style>
    </Layer>
    <Layer queryable="1">
      <Name>OSM-Overlay-WMS</Name>
      <Title>OSM Overlay WMS - by terrestris</Title>
      <LatLonBoundingBox minx="-180" miny="-88" maxx="180" maxy="88" />
      <BoundingBox SRS="EPSG:900913" minx="-20037508.3428" miny="-25819498.5135" maxx="20037508.3428" maxy="25819498.5135" />
      <BoundingBox SRS="EPSG:4326" minx="-180" miny="-88" maxx="180" maxy="88" />
      <BoundingBox SRS="EPSG:3857" minx="-20037508.3428" miny="-25819498.5135" maxx="20037508.3428" maxy="25819498.5135" />
    </Layer>
    <Layer queryable="1">
      <Name>TOPO-WMS</Name>
      <Title>Topographic WMS - by terrestris</Title>
      <LatLonBoundingBox minx="-180" miny="-88" maxx="180" maxy="88" />
      <BoundingBox SRS="EPSG:900913" minx="-20037508.3428" miny="-25819498.5135" maxx="20037508.3428" maxy="25819498.5135" />
      <BoundingBox SRS="EPSG:4326" minx="-180" miny="-88" maxx="180" maxy="88" />
      <BoundingBox SRS="EPSG:3857" minx="-20037508.3428" miny="-25819498.5135" maxx="20037508.3428" maxy="25819498.5135" />
    </Layer>
    <Layer queryable="1">
      <Name>TOPO-OSM-WMS</Name>
      <Title>Topographic OSM WMS - by terrestris</Title>
      <LatLonBoundingBox minx="-180" miny="-88" maxx="180" maxy="88" />
      <BoundingBox SRS="EPSG:900913" minx="-20037508.3428" miny="-25819498.5135" maxx="20037508.3428" maxy="25819498.5135" />
      <BoundingBox SRS="EPSG:4326" minx="-180" miny="-88" maxx="180" maxy="88" />
      <BoundingBox SRS="EPSG:3857" minx="-20037508.3428" miny="-25819498.5135" maxx="20037508.3428" maxy="25819498.5135" />
    </Layer>
    <Layer>
      <Name>SRTM30-Hillshade</Name>
      <Title>SRTM30 Hillshade - by terrestris</Title>
      <LatLonBoundingBox minx="-180" miny="-56" maxx="180" maxy="60" />
      <BoundingBox SRS="EPSG:900913" minx="-20037508.3428" miny="-7558415.65608" maxx="20037508.3428" maxy="8399737.88982" />
      <BoundingBox SRS="EPSG:4326" minx="-180" miny="-56" maxx="180" maxy="60" />
      <BoundingBox SRS="EPSG:3857" minx="-20037508.3428" miny="-7558415.65608" maxx="20037508.3428" maxy="8399737.88982" />
    </Layer>
    <Layer>
      <Name>SRTM30-Colored</Name>
      <Title>SRTM30 Colored - by terrestris</Title>
      <LatLonBoundingBox minx="-180" miny="-56" maxx="180" maxy="60" />
      <BoundingBox SRS="EPSG:900913" minx="-20037508.3428" miny="-7558415.65608" maxx="20037508.3428" maxy="8399737.88982" />
      <BoundingBox SRS="EPSG:4326" minx="-180" miny="-56" maxx="180" maxy="60" />
      <BoundingBox SRS="EPSG:3857" minx="-20037508.3428" miny="-7558415.65608" maxx="20037508.3428" maxy="8399737.88982" />
    </Layer>
    <Layer>
      <Name>SRTM30-Colored-Hillshade</Name>
      <Title>SRTM30 Colored Hillshade - by terrestris</Title>
      <LatLonBoundingBox minx="-180" miny="-56" maxx="180" maxy="60" />
      <BoundingBox SRS="EPSG:900913" minx="-20037508.3428" miny="-7558415.65608" maxx="20037508.3428" maxy="8399737.88982" />
      <BoundingBox SRS="EPSG:4326" minx="-180" miny="-56" maxx="180" maxy="60" />
      <BoundingBox SRS="EPSG:3857" minx="-20037508.3428" miny="-7558415.65608" maxx="20037508.3428" maxy="8399737.88982" />
    </Layer>
    <Layer>
      <Name>SRTM30-Contour</Name>
      <Title>SRTM30 Contour Lines - by terrestris</Title>
      <LatLonBoundingBox minx="-180" miny="-56" maxx="180" maxy="60" />
      <BoundingBox SRS="EPSG:900913" minx="-20037508.3428" miny="-7558415.65608" maxx="20037508.3428" maxy="8399737.88982" />
      <BoundingBox SRS="EPSG:4326" minx="-180" miny="-56" maxx="180" maxy="60" />
      <BoundingBox SRS="EPSG:3857" minx="-20037508.3428" miny="-7558415.65608" maxx="20037508.3428" maxy="8399737.88982" />
    </Layer>
  </Layer>
</Capability>
</WMT_MS_Capabilities>`
