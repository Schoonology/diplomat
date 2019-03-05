package http

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type client struct {
	Address string
	client  http.Client
}

func NewClient(address string) *client {
	return &client{
		Address: address,
		client: http.Client{
			Transport: &http.Transport{
				DisableCompression: true,
			},
		},
	}
}

func (c *client) Do(request *Request) (*Response, error) {
	nativeAddress, err := url.Parse(c.Address)
	if err != nil {
		return nil, err
	}

	nativeAddress, err = nativeAddress.Parse(request.Path)
	if err != nil {
		return nil, err
	}

	nativeRequest, err := http.NewRequest(request.Method, nativeAddress.String(), nil)
	if err != nil {
		return nil, err
	}

	for key, value := range request.Headers {
		nativeRequest.Header.Set(key, value)
	}

	nativeResponse, err := c.client.Do(nativeRequest)
	if err != nil {
		return nil, err
	}

	defer nativeResponse.Body.Close()

	response := NewResponse(nativeResponse.StatusCode, strings.Join(strings.Split(nativeResponse.Status, " ")[1:], " "))
	for key, value := range nativeResponse.Header {
		response.Headers[key] = value[0]
	}

	response.Body, err = ioutil.ReadAll(nativeResponse.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}
