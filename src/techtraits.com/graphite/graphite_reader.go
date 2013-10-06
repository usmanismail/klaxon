package graphite

import (
	"appengine"
	"appengine/urlfetch"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type GraphiteReader interface {
	ReadValue(target string) (float64, error)
}

type GraphiteHandler struct {
	BaseUrl        *url.URL
	LookBack       string
	GraphiteClient *http.Client
}

func MakeGraphiteReader(baseUrlStr string, context appengine.Context) (GraphiteReader, error) {

	baseUrl, err := url.Parse(baseUrlStr + "/render")

	var reader GraphiteHandler
	if err == nil {
		reader = GraphiteHandler{baseUrl, "-300sec", urlfetch.Client(context)}
	}

	return reader, err
}

func (this GraphiteHandler) ReadValue(target string) (float64, error) {

	var value float64

	graphiteUrl, err := url.Parse(this.BaseUrl.String() + "?target=" + target + "&format=csv&from=" + this.LookBack)
	if err == nil {
		var resp *http.Response
		resp, err = this.GraphiteClient.Get(graphiteUrl.String())
		if err == nil {
			defer resp.Body.Close()
			var content []byte
			content, err = ioutil.ReadAll(resp.Body)
			recordStrings := strings.Split(strings.TrimSpace(string(content[:])), "\n")
			err = errors.New("no-known-value")
			for _, recordString := range recordStrings {
				records := strings.Split(strings.TrimSpace(recordString), ",")
				if len(records) == 3 && len(records[2]) != 0 {
					err = nil
					value, _ = strconv.ParseFloat(records[2], 64)
				}
			}
		}
	}
	return value, err
}
