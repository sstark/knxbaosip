package main

import (
	"fmt"
	"log"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	knx := NewClient("")

	err, si := knx.GetServerItem()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s fw:%d sn:%v\n", knx.Url, si.FirmwareVersion, si.SerialNumber)

	err, dpd := knx.GetDatapointDescription([]int{700, 701, 711, 712, 720, 721, 722})
	if err != nil {
		log.Fatal(err)
	}
	for _, d := range dpd {
		fmt.Printf("%d:%d\n", d.Datapoint, d.DatapointType)
	}

	err, ds := knx.GetDescriptionString([]int{700, 701, 711, 712, 720, 721, 722})
	if err != nil {
		log.Fatal(err)
	}
	for _, d := range ds {
		fmt.Printf("%d:%s\n", d.Datapoint, d.Description)
	}
}
