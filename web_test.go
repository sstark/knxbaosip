package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var ApiTestUrlMap = map[string]string{
	"/baos/getServerItem?ItemStart=1&ItemCount=18":                               "testdata/results/getServerItem.json",
	"/baos/getDatapointDescription?DatapointStart=1&DatapointCount=33":           "testdata/results/getDatapointDescription-1-33.json",
	"/baos/getDescriptionString?DatapointStart=1&DatapointCount=33":              "testdata/results/getDescriptionString-1-33.json",
	"/baos/getDescriptionString?DatapointStart=543&DatapointCount=3":             "testdata/results/invalidParam.json",
	"/baos/getDatapointValue?DatapointStart=711&DatapointCount=1&Format=Default": "testdata/results/getDatapointValue-711.json",
}

func makeTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		out, err := ioutil.ReadFile(ApiTestUrlMap[r.URL.RequestURI()])
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Fprint(w, string(out))
	}))
}

func setup(t *testing.T) (func(t *testing.T), *Client) {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	ts := makeTestServer()
	return func(t *testing.T) { ts.Close() }, NewClient(ts.URL + "/baos/")
}

func TestGetServerItem(t *testing.T) {
	tearDown, knx := setup(t)
	defer tearDown(t)
	_, si := knx.GetServerItem()
	got := si.ApplicationId
	wanted := 1801
	if got != wanted {
		t.Errorf("got %d, wanted %d", got, wanted)
	}
}

func TestGetDatapointDescription(t *testing.T) {
	tearDown, knx := setup(t)
	defer tearDown(t)
	dps := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33}
	_, ds := knx.GetDatapointDescription(dps)
	got := ds[7].DatapointType
	wanted := 3
	if got != wanted {
		t.Errorf("got %d, wanted %d", got, wanted)
	}
}

func TestGetDescriptionString(t *testing.T) {
	tearDown, knx := setup(t)
	defer tearDown(t)
	dps := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33}
	_, ds := knx.GetDescriptionString(dps)
	got := ds[10].Description
	wanted := "Jalo. N4.015 Auf/Ab"
	if got != wanted {
		t.Errorf("got %s, wanted %s", got, wanted)
	}
}

func TestGetDescriptionStringInvalidParam(t *testing.T) {
	tearDown, knx := setup(t)
	defer tearDown(t)
	_, ds := knx.JsonGetDescriptionString(543, 3)
	got := ds.Error
	wanted := "InvalidParam"
	if got != wanted {
		t.Errorf("got %s, wanted %s", got, wanted)
	}
}

func TestGetDatapointValue(t *testing.T) {
	tearDown, knx := setup(t)
	defer tearDown(t)
	dps := []int{711}
	_, dpv := knx.GetDatapointValue(dps)
	// the value is a json.RawMessage, so we cast it for comparison
	got := string(dpv[0].Value)
	wanted := "3198"
	if got != wanted {
		t.Errorf("got %s, wanted %s", got, wanted)
	}
}
