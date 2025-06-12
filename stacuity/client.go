// Copyright (c) HashiCorp, Inc.

package stacuity

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// HostURL - Default Stacuity API URL
const HostURL string = "http://localhost:5050/api/v1"

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

func New(text string) error {

	return &errorString{text}
}

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func NewClient(host, authToken *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    HostURL,
		Token:      *authToken,
	}

	if host != nil {
		c.HostURL = *host
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
