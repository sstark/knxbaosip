package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	defaultUrl         string = "http://localhost:8888/baos/"
	GetServerItemCount int    = 18
)

type Client struct {
	Url string
}

type JsonResult struct {
	Result  bool
	Service string
	Error   string
	Data    json.RawMessage
}

type JsonServerItem struct {
	HardwareType               []int
	HardwareVersion            int
	FirmwareVersion            int
	KnxManufacturerCodeDev     int
	KnxManufacturerCodeApp     int
	ApplicationId              int
	ApplicationVersion         int
	SerialNumber               []int
	TimeSinceReset             int
	BusConnectionState         int
	MaximalBufferSize          int
	LengthOfDescriptionString  int
	Baudrate                   int
	CurrentBufferSize          int
	ProgrammingMode            int
	ProtocolVersion            int
	IndicationSending          int
	ProtocolVersionWebServices int
}

type JsonDatapointDescription struct {
	Datapoint          int
	ValueType          int
	ConfigurationFlags int
	DatapointType      int
}

type JsonDescriptionString struct {
	Datapoint   int
	Description string
}

func NewClient(url string) *Client {
	var apiUrl string
	if url == "" {
		apiUrl = defaultUrl
	} else {
		apiUrl = url
	}
	return &Client{Url: apiUrl}
}

func (a *Client) ApiGetJson(serviceQuery string) (error, string) {
	getPath := fmt.Sprintf("%s%s", a.Url, serviceQuery)
	res, err := http.Get(getPath)
	if err != nil {
		return fmt.Errorf("http GET error: %s", err), ""
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return fmt.Errorf("http read error: %s", err), ""
	}
	return nil, string(result)
}

func (a *Client) JsonGetServerItem() (error, JsonResult) {
	var m JsonResult
	var err error
	uri := fmt.Sprintf("getServerItem?ItemStart=1&ItemCount=%d", GetServerItemCount)
	err, out := a.ApiGetJson(uri)
	if err != nil {
		return err, m
	}
	j := []byte(out)
	err = json.Unmarshal(j, &m)
	if err != nil {
		return fmt.Errorf("Error decoding message: %s", err), m
	}
	if m.Result == false {
		return fmt.Errorf("%s: BAOS error: %s", uri, m.Error), m
	}
	return nil, m
}

func (a *Client) GetServerItem() (error, JsonServerItem) {
	var t JsonResult
	var m JsonServerItem
	var err error
	err, t = a.JsonGetServerItem()
	if err != nil {
		return err, m
	}
	err = json.Unmarshal(t.Data, &m)
	if err != nil {
		return fmt.Errorf("Error decoding data from message: %s", err), m
	}
	return nil, m
}

// JsonGetDatapointDescription fetches <count> consecutive datapoints from the server
// and returns the raw json data.
func (a *Client) JsonGetDatapointDescription(datapoint int, count int) (error, JsonResult) {
	var m JsonResult
	var err error
	uri := fmt.Sprintf("getDatapointDescription?DatapointStart=%d&DatapointCount=%d", datapoint, count)
	err, out := a.ApiGetJson(uri)
	if err != nil {
		return err, m
	}
	j := []byte(out)
	err = json.Unmarshal(j, &m)
	if err != nil {
		return fmt.Errorf("Error decoding message: %s", err), m
	}
	if m.Result == false {
		return fmt.Errorf("%s: BAOS error: %s", uri, m.Error), m
	}
	return nil, m
}

// GetDatapointDescription takes a list of datapoints and tries to fetch them with as little
// calls to JsonGetDatapointDescription as possible.
func (a *Client) GetDatapointDescription(datapoints []int) (error, []JsonDatapointDescription) {
	var m, t []JsonDatapointDescription
	var r JsonResult
	var err error
	for _, chunk := range makeChunks(datapoints) {
		err, r = a.JsonGetDatapointDescription(chunk[0], chunk[1])
		if err != nil {
			return err, m
		}
		err = json.Unmarshal(r.Data, &t)
		if err != nil {
			return fmt.Errorf("Error decoding data from message: %s", err), m
		}
		m = append(m, t...)
	}
	return nil, m
}

func (a *Client) JsonGetDescriptionString(datapoint int) (error, string) {
	return a.ApiGetJson(fmt.Sprintf("getDescriptionString?DatapointStart=%d&DatapointCount=1", datapoint))
}
