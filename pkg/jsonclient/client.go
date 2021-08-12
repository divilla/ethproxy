package jsonclient

import (
	"bytes"
	"fmt"
	"github.com/divilla/ethproxy/config"
	"io"
	"io/ioutil"
	"net/http"
)

type (
	JsonHttpClient struct {
		Url string
	}
)

func New(url string) *JsonHttpClient {
	return &JsonHttpClient{
		Url: url,
	}
}

func (c *JsonHttpClient) Post(request string) ([]byte, error) {
	var resp *http.Response
	var err error
	var body []byte

	for i:=0; i<config.FetchRetries; i++ {
		resp, err = http.Post(c.Url, "application/json", bytes.NewBuffer([]byte(request)))
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("http POST request to '%s' with body '%s' failed with: %w", c.Url, request, err)
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return body, nil
}
