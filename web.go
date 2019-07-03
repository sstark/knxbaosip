package knxbaosip

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	_ int = iota
	DPT1
	DPT2
	DPT3
	DPT4
	DPT5
	DPT6
	DPT7
	DPT8
	DPT9
	DPT10
	DPT11
	DPT12
	DPT13
	DPT14
	DPT15
	DPT16
	DPT17
	DPT18
	defaultUrl         string = "http://localhost:8888/baos/"
	GetServerItemCount int    = 18
)

var (
	AuthError = errors.New("Authorisation required")
)

type Client struct {
	Url    string
	Logger *log.Logger
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
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	if url == "" {
		apiUrl = defaultUrl
	} else {
		apiUrl = url
	}
	if !strings.HasSuffix(apiUrl, "/") {
		apiUrl = apiUrl + "/"
	}
	return &Client{Url: apiUrl, Logger: logger}
}

func (a *Client) Debugf(format string, args ...interface{}) {
	if os.Getenv("KNXBAOSIP_DEBUG") == "1" {
		a.Logger.Output(2, "<DEBUG> "+fmt.Sprintf(format, args...))
	}
}

func (a *Client) SetDebugLevel(level int) {
	os.Setenv("KNXBAOSIP_DEBUG", fmt.Sprintf("%d", level))
}

// ApiGetJson queries the baos gateway with a given service query and returns
// the result as a slice of bytes
func (a *Client) ApiGetJson(serviceQuery string) (error, []byte) {
	getPath := fmt.Sprintf("%s%s", a.Url, serviceQuery)
	a.Debugf(getPath)
	res, err := http.Get(getPath)
	a.Debugf("%v", err)
	if err != nil {
		return fmt.Errorf("http GET error: %s", err), nil
	}
	a.Debugf("status: %s", res.Status)
	if res.StatusCode == 401 {
		return AuthError, nil
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

// SetDatapointValue sets a given datapoint to the given value. The type of the
// value depends on the supplied DPT format.
func (a *Client) SetDatapointValue(datapoint int, format int, value interface{}) (error, JsonResult) {
	var m JsonResult
	var err error
	var uri string
	switch val := value.(type) {
	case string:
		switch format {
		case DPT1:
			uri = fmt.Sprintf("setDatapointValue?Datapoint=%d&Format=DPT%d&Length=1&Value=%s", datapoint, format, val)
		default:
			return errors.New("unsupported DPT format"), m
		}
	case int:
		switch format {
		case DPT5:
			uri = fmt.Sprintf("setDatapointValue?Datapoint=%d&Format=DPT%d&Length=1&Value=%d", datapoint, format, val)
		default:
			return errors.New("unsupported DPT format"), m
		}
	default:
		return errors.New("unsupported value type"), m
	}
	a.Debugf(uri)
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
