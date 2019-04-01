package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	defaultUrl         string = "http://localhost:8888/baos"
	GetServerItemCount int    = 18
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

func (a *Client) ApiGetJSON(service string, count int) string {
	getPath := fmt.Sprintf("%s%s?ItemStart=1&ItemCount=%d", a.Url, service, count)
	res, err := http.Get(getPath)
	if err != nil {
		log.Fatal(err)
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return string(result)
}

func (a *Client) JsonGetServerItem() string {
	return a.ApiGetJSON("getServerItem", 18)
}
