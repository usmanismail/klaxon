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

//Make a graphite reader which looks reads metrics from graphite instance(s) referenced by BaseUrl,
//and uses the last n seconds to determine find value, where 'n' is equal to lookback
func MakeGraphiteReader(baseUrlStr string, lookback string, context appengine.Context) (GraphiteReader, error) {

	baseUrl, err := url.Parse(baseUrlStr + "/render")

	var reader GraphiteHandler
	if err == nil {
		reader = GraphiteHandler{baseUrl, "-" + lookback + "sec", urlfetch.Client(context)}
	}

	return reader, err
}

//Read the metric value specified by target
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
