// Package content supports getcap & getmap packages
package content

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sgaunet/wms/urlmap"
)

// From return data from a URL with Basic Auth
func From(url *urlmap.URLmap) (b *bytes.Reader, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url.String(), nil)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("status Code Error %v", resp.StatusCode)
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return bytes.NewReader(data), nil
}
