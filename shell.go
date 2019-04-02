package main

import (
	"fmt"
	"log"
	"strings"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	knx := NewClient("")

	err, si := knx.GetServerItem()
	if err != nil {
		log.Fatal(err)
	}
	sn := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(si.SerialNumber)), "."), "[]")
	fmt.Printf("%s fw:%d sn:%v\n", knx.Url, si.FirmwareVersion, sn)

	datapoints := []int{700, 701, 711, 712, 720, 721, 722}

	err, ds := knx.GetDescriptionString(datapoints)
	if err != nil {
		log.Fatal(err)
	}
	err, dpv := knx.GetDatapointValue(datapoints)
	if err != nil {
		log.Fatal(err)
	}
	for i, d := range dpv {
		desc := ds[i].Description
		fmt.Printf("%5d %5s \"%s\": %s\n", d.Datapoint, d.Format, desc, string(d.Value))
	}
}
