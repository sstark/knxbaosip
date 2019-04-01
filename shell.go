package main

func main() {
	knx := NewClient("http://localhost:8888/baos/")
	println(knx.JsonGetServerItem())
	println(knx.JsonGetDataPointDescription(711))
	println(knx.JsonGetDescriptionString(711))
}
