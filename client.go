package main

import (
	"net/http"
	"net/http/cookiejar"
)

type ClientOption func(*Client) error

func ClientWithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) error {
		c.httpClient = client
		return nil
	}
}

type Client struct {
	httpClient *http.Client
}

func NewClient(opts ...ClientOption) (*Client, error) {
	client := Client{}

	for _, opt := range opts {
		if err := opt(&client); err != nil {
			return nil, err
		}
	}

	if client.httpClient == nil {
		jar, err := cookiejar.New(nil)
		if err != nil {
			return nil, err
		}

		client.httpClient = &http.Client{
			Jar: jar,
		}
	}

	return &client, nil
}

func (s *Client) Do(req *http.Request) (*http.Response, error) {
	return s.httpClient.Do(req)
}
