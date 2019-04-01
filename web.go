package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	defaultUrl         string = "http://localhost:8888/baos"
	GetServerItemCount int    = 18
)

type Client struct {
	Url string
}

type JsonResult struct {
	Result  bool
	Service string
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

func NewClient(url string) *Client {
	var apiUrl string
	if url == "" {
		apiUrl = defaultUrl
	} else {
		apiUrl = url
	}
	return &Client{Url: apiUrl}
}

func (a *Client) ApiGetJson(serviceQuery string) string {
	getPath := fmt.Sprintf("%s%s", a.Url, serviceQuery)
	fmt.Println(getPath)
	res, err := http.Get(getPath)
	if err != nil {
		log.Fatal(err)
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return string(result)
}

func (a *Client) JsonGetServerItem() JsonResult {
	var m JsonResult
	j := []byte(a.ApiGetJson(fmt.Sprintf("getServerItem?ItemStart=1&ItemCount=%d", GetServerItemCount)))
	err := json.Unmarshal(j, &m)
	if err != nil {
		log.Fatal(err)
	}
	return m
}

func (a *Client) GetServerItem() JsonServerItem {
	var m JsonServerItem
	err := json.Unmarshal(a.JsonGetServerItem().Data, &m)
	if err != nil {
		log.Fatal(err)
	}
	return m
}

func (a *Client) JsonGetDataPointDescription(datapoint int) string {
	return a.ApiGetJson(fmt.Sprintf("getDatapointDescription?DatapointStart=%d&DatapointCount=1", datapoint))
}

func (a *Client) JsonGetDescriptionString(datapoint int) string {
	return a.ApiGetJson(fmt.Sprintf("getDescriptionString?DatapointStart=%d&DatapointCount=1", datapoint))
}
