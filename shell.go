package main

func main() {
	knx := NewClient("http://localhost:8888/baos/")
	m := knx.JsonGetServerItem()
	println(m.Data)
	println(knx.JsonGetDataPointDescription(711))
	println(knx.JsonGetDescriptionString(711))
}
