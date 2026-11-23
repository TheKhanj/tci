package main

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const defaultBaseURL = "https://internet.tci.ir"

type ClientOption func(*Client) error

func ClientWithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		baseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}

		c.baseURL = baseURL
		return nil
	}
}

func ClientWithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) error {
		c.httpClient = client
		return nil
	}
}

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
}

func NewClient(opts ...ClientOption) (*Client, error) {
	client := Client{}

	for _, opt := range opts {
		if err := opt(&client); err != nil {
			return nil, err
		}
	}

	if client.baseURL == nil {
		if err := ClientWithBaseURL(defaultBaseURL)(&client); err != nil {
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
