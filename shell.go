package main

import (
	"fmt"
)

func main() {
	knx := NewClient("http://localhost:8888/baos/")

	si := knx.GetServerItem()
	fmt.Printf("firmware %d, serialnumber %v\n", si.FirmwareVersion, si.SerialNumber)

	dpd := knx.GetDatapointDescription([]int{700, 701, 711, 712, 720, 721, 722})
	fmt.Printf("%+v\n", dpd)

	//	println(knx.JsonGetDescriptionString(711))
}
