package knxbaosip

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	defaultUrl string = "http://localhost:8888/baos"
)

type Client struct {
	Url string
}

func NewClient(url string) *Client {
	var apiUrl string
	if url == "" {
		apiUrl = defaultUrl
	} else {
		apiUrl = url
	}
	return &Client{Url: apiUrl}
}

func (a *Client) JsonGetServerItem() string {
	getPath := a.Url + "getServerItem"
	fmt.Println(getPath)
	res, err := http.Get(getPath)
	if err != nil {
		log.Fatal(err)
	}
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return string(greeting)
}
