
knxbaosip
=========

Simple Go client implementation of the KNX IP BAOS application layer web
interface, as specified for accessing the following home automation devices
from Weinzierl:

  - KNX IP BAOS 771
  - KNX IP BAOS 772
  - KNX IP BAOS 773
  - KNX IP BAOS 774
  - KNX IP BAOS 777 

Original API documentation:

  - https://www.weinzierl.de/images/download/documents/baos/KNX_IP_BAOS_WebServices.pdf

For setting, currently only DPT1 and DPT5 are supported.

The values (as returned by GetDatapointValue()) are not cast to a certain data
type. Instead they are returned as a json.RawMessage. You can simply cast them
to string or use your own conversion for now. Maybe later some higher level
methods will be added.


Example
=======

See https://github.com/sstark/kxsh for an example using this package.


Interface
=========

    func NewClient(url string) *Client
        NewClient creates a new client object using the given URL to access the
        knx baos ip gateway.


Low Level Functions
-------------------

    func (a *Client) ApiGetJson(serviceQuery string) (error, []byte)
        ApiGetJson queries the baos gateway with a given service query and
        returns the result as a slice of bytes

    func (a *Client) JsonGetDatapointDescription(datapoint int, count int) (error, JsonResult)
        JsonGetDatapointDescription fetches <count> consecutive datapoints from
        the server and returns the raw json data.

    func (a *Client) JsonGetDatapointValue(datapoint int, count int) (error, JsonResult)
        JsonGetDatapointValue fetches <count> consecutive datapoints from the
        server and returns the raw json data.

    func (a *Client) JsonGetDescriptionString(datapoint int, count int) (error, JsonResult)
        JsonGetDescriptionString fetches <count> consecutive datapoints from the
        server and returns the raw json data.

    func (a *Client) JsonGetServerItem() (error, JsonResult)
        JsonGetServerItem returns some basic gateway information as a generic
        json result object


Mid Level Functions
-------------------

    func (a *Client) GetDatapointDescription(datapoints []int) (error, []JsonDatapointDescription)
        GetDatapointDescription takes a list of datapoints and tries to fetch
        them with as little calls to JsonGetDatapointDescription as possible.

    func (a *Client) GetDatapointValue(datapoints []int) (error, []JsonDatapointValue)
        GetDatapointValue takes a list of datapoints and tries to fetch them
        with as little calls to JsonGetDatapointValue as possible. The actual
        value is a json.RawMessage in each elements Value field.

    func (a *Client) GetDescriptionString(datapoints []int) (error, []JsonDescriptionString)
        GetDescriptionString takes a list of datapoints and tries to fetch them
        with as little calls to JsonGetDescriptionString as possible.

    func (a *Client) GetServerItem() (error, JsonServerItem)
        GetServerItem returns some basic gateway information as a specific
        result object

    func (a *Client) SetDatapointValue(datapoint int, format int, value interface{}) (error, JsonResult)
        SetDatapointValue sets a given datapoint to the given value. The type of
        the value depends on the supplied DPT format.


Types
-----

    type JsonDatapointDescription struct {
        Datapoint          int
        ValueType          int
        ConfigurationFlags int
        DatapointType      int
    }

    type JsonDatapointValue struct {
        Datapoint int
        Format    string
        Length    int
        State     int
        Value json.RawMessage
    }

    type JsonDescriptionString struct {
        Datapoint   int
        Description string
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

Constants
---------

DPT1 DPT2 DPT3 DPT4 DPT5 DPT6 DPT7 DPT8 DPT9 DPT10 DPT11 DPT12 DPT13 DPT14 DPT15 DPT16 DPT17 DPT18 correspond to the integer values 1 to 18.

