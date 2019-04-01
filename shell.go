package main

import (
	"fmt"
)

func main() {
	knx := NewClient("")
	si := knx.GetServerItem()
	fmt.Printf("%s fw:%d sn:%v\n", knx.Url, si.FirmwareVersion, si.SerialNumber)

	dpd := knx.GetDatapointDescription([]int{700, 701, 711, 712, 720, 721, 722})
	for _, d := range dpd {
		fmt.Printf("%d:%d\n", d.Datapoint, d.DatapointType)
	}

	//	println(knx.JsonGetDescriptionString(711))
}
