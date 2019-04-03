package knxbaosip

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

type JsonDatapointValue struct {
	Datapoint int
	Format    string
	Length    int
	State     int
	// depends on Format
	Value json.RawMessage
}

// NewClient creates a new client object using the given URL to access the knx
// baos ip gateway.
func NewClient(url string) *Client {
	var apiUrl string
	if url == "" {
		apiUrl = defaultUrl
	} else {
		apiUrl = url
	}
	return &Client{Url: apiUrl}
}

// ApiGetJson queries the baos gateway with a given service query and returns
// the result as a slice of bytes
func (a *Client) ApiGetJson(serviceQuery string) (error, []byte) {
	getPath := fmt.Sprintf("%s%s", a.Url, serviceQuery)
	res, err := http.Get(getPath)
	if err != nil {
		return fmt.Errorf("http GET error: %s", err), nil
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return fmt.Errorf("http read error: %s", err), nil
	}
	return nil, result
}

// JsonGetServerItem returns some basic gateway information as a generic json
// result object
func (a *Client) JsonGetServerItem() (error, JsonResult) {
	var m JsonResult
	var err error
	uri := fmt.Sprintf("getServerItem?ItemStart=1&ItemCount=%d", GetServerItemCount)
	err, out := a.ApiGetJson(uri)
	if err != nil {
		return err, m
	}
	err = json.Unmarshal(out, &m)
	if err != nil {
		return fmt.Errorf("Error decoding message: %s", err), m
	}
	if m.Result == false {
		return fmt.Errorf("%s: BAOS error: %s", uri, m.Error), m
	}
	return nil, m
}

// GetServerItem returns some basic gateway information as a specific
// result object
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
	err = json.Unmarshal(out, &m)
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
	var m []JsonDatapointDescription
	var r JsonResult
	var err error
	for _, chunk := range makeChunks(datapoints) {
		err, r = a.JsonGetDatapointDescription(chunk[0], chunk[1])
		if err != nil {
			return err, m
		}
		var t []JsonDatapointDescription
		err = json.Unmarshal(r.Data, &t)
		if err != nil {
			return fmt.Errorf("Error decoding data from message: %s", err), m
		}
		m = append(m, t...)
	}
	return nil, m
}

// JsonGetDescriptionString fetches <count> consecutive datapoints from the server
// and returns the raw json data.
func (a *Client) JsonGetDescriptionString(datapoint int, count int) (error, JsonResult) {
	var m JsonResult
	var err error
	uri := fmt.Sprintf("getDescriptionString?DatapointStart=%d&DatapointCount=%d", datapoint, count)
	err, out := a.ApiGetJson(uri)
	if err != nil {
		return err, m
	}
	err = json.Unmarshal(out, &m)
	if err != nil {
		return fmt.Errorf("Error decoding message: %s", err), m
	}
	if m.Result == false {
		return fmt.Errorf("%s: BAOS error: %s", uri, m.Error), m

	}
	return nil, m
}

// GetDescriptionString takes a list of datapoints and tries to fetch them with as little
// calls to JsonGetDescriptionString as possible.
func (a *Client) GetDescriptionString(datapoints []int) (error, []JsonDescriptionString) {
	var m []JsonDescriptionString
	var r JsonResult
	var err error
	for _, chunk := range makeChunks(datapoints) {
		err, r = a.JsonGetDescriptionString(chunk[0], chunk[1])
		if err != nil {
			return err, m
		}
		var t []JsonDescriptionString
		err = json.Unmarshal(r.Data, &t)
		if err != nil {
			return fmt.Errorf("Error decoding data from message: %s", err), m
		}
		m = append(m, t...)
	}
	return nil, m
}

// JsonGetDatapointValue fetches <count> consecutive datapoints from the server
// and returns the raw json data.
func (a *Client) JsonGetDatapointValue(datapoint int, count int) (error, JsonResult) {
	var m JsonResult
	var err error
	uri := fmt.Sprintf("getDatapointValue?DatapointStart=%d&DatapointCount=%d&Format=Default", datapoint, count)
	err, out := a.ApiGetJson(uri)
	if err != nil {
		return err, m
	}
	err = json.Unmarshal(out, &m)
	if err != nil {
		return fmt.Errorf("Error decoding message: %s", err), m
	}
	if m.Result == false {
		return fmt.Errorf("%s: BAOS error: %s", uri, m.Error), m

	}
	return nil, m
}

// GetDatapointValue takes a list of datapoints and tries to fetch them with as little
// calls to JsonGetDatapointValue as possible.
// The actual value is a json.RawMessage in each elements Value field.
func (a *Client) GetDatapointValue(datapoints []int) (error, []JsonDatapointValue) {
	var m []JsonDatapointValue
	var r JsonResult
	var err error
	for _, chunk := range makeChunks(datapoints) {
		err, r = a.JsonGetDatapointValue(chunk[0], chunk[1])
		if err != nil {
			return err, m
		}
		// this needs to be reset on every iteration
		var t []JsonDatapointValue
		err = json.Unmarshal(r.Data, &t)
		if err != nil {
			return fmt.Errorf("Error decoding data from message: %s", err), m
		}
		m = append(m, t...)
	}
	return nil, m
}
