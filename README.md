
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

Currently only getting values is supporting, not setting.

The values (as returned by GetDatapointValue()) are not cast to a certain data
type. Instead they are returned as a json.RawMessage. You can simply cast them
to string or use your own conversion for now. Maybe later some higher level
methods will be added.

Example
=======



Interface
=========

    func NewClient(url string) *Client

Low Level Functions
-------------------

    func (a *Client) ApiGetJson(serviceQuery string) (error, string)

    func (a *Client) JsonGetDatapointDescription(datapoint int, count int) (error, JsonResult)
        JsonGetDatapointDescription fetches <count> consecutive datapoints from
        the server and returns the raw json data.

    func (a *Client) JsonGetDatapointValue(datapoint int, count int) (error, JsonResult)

    func (a *Client) JsonGetDescriptionString(datapoint int, count int) (error, JsonResult)
        JsonGetDescriptionString fetches <count> consecutive datapoints from the
        server and returns the raw json data.

    func (a *Client) JsonGetServerItem() (error, JsonResult)


Mid Level Functions
-------------------

    func (a *Client) GetDatapointDescription(datapoints []int) (error, []JsonDatapointDescription)
        GetDatapointDescription takes a list of datapoints and tries to fetch
        them with as little calls to JsonGetDatapointDescription as possible.

    func (a *Client) GetDatapointValue(datapoints []int) (error, []JsonDatapointValue)

    func (a *Client) GetDescriptionString(datapoints []int) (error, []JsonDescriptionString)
        GetDescriptionString takes a list of datapoints and tries to fetch them
        with as little calls to JsonGetDescriptionString as possible.

    func (a *Client) GetServerItem() (error, JsonServerItem)


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

