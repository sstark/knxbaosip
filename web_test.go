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
	"/baos/getServerItem": "testdata/results/getServerItem.json",
}

func makeTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		out, err := ioutil.ReadFile(ApiTestUrlMap[r.URL.Path])
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Fprintln(w, out)
	}))
}

func TestApi(t *testing.T) {
	ts := makeTestServer()
	defer ts.Close()
	knx := NewClient(ts.URL + "/baos/")
	knx.JsonGetServerItem()
}
