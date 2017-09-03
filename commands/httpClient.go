package commands

import (
	"net/http"
	"io/ioutil"
	"bytes"
	"log"
)

type ClientMock struct {
}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	log.Printf("hi from ClientMock Do()")
	zip := []byte{80, 75, 05, 06, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00}

	resp := http.Response{
		Body: ioutil.NopCloser(bytes.NewBuffer(zip)),
	}
	return &resp, nil
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
