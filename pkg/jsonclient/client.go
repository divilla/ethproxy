package jsonclient

import (
	"bytes"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/divilla/ethproxy/config"
	"github.com/divilla/ethproxy/interfaces"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
)

type (
	JsonHttpClient struct {
		url    string
		logger interfaces.ErrorLogger
	}
)

func New(logger interfaces.ErrorLogger) *JsonHttpClient {
	return &JsonHttpClient{
		logger: logger,
	}
}
func (c *JsonHttpClient) Url(url string) error {
	if !govalidator.IsURL(url) {
		return errors.Errorf("'%s' is not valid url", url)
	}
	c.url = url

	return nil
}

func (c *JsonHttpClient) Post(request string) ([]byte, error) {
	var resp *http.Response
	var err error
	var body []byte

	for i := 0; i < config.FetchRetries; i++ {
		resp, err = http.Post(c.url, "application/json", bytes.NewBuffer([]byte(request)))
		if err != nil {
			c.logger.Errorf("unable to fetch '%s', with body '%s', retry '%d/%d', with error: '%e'", c.url, request, i+1, config.FetchRetries, err)
		} else {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("http POST request to '%s', with body '%s', with '%d' retries, failed with: %w", c.url, request, config.FetchRetries, err)
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
