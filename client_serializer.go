package tci

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const (
	cookiePHPSESSID = "PHPSESSID"
	cookieSession1  = "cookiesession1"
)

type ClientSerializerOption func(*ClientSerializer) error

func ClientSerializerWithBaseURL(baseURL string) ClientSerializerOption {
	return WithBaseURL(baseURL, func(s *ClientSerializer, u *url.URL) {
		s.baseURL = u
	})
}

func NewClientSerializer(
	opts ...ClientSerializerOption,
) (*ClientSerializer, error) {
	s := ClientSerializer{}

	for _, opt := range opts {
		if err := opt(&s); err != nil {
			return nil, err
		}
	}

	if s.baseURL == nil {
		if err := ClientSerializerWithBaseURL(DefaultBaseURL)(&s); err != nil {
			return nil, err
		}
	}

	return &s, nil
}

type ClientSerializer struct {
	baseURL *url.URL
}

type serializedClient struct {
	Cookies []*http.Cookie `json:"cookies"`
}

func (s *ClientSerializer) Marshal(c *Client) ([]byte, error) {
	if c == nil || c.httpClient == nil {
		return nil, errors.New("nil client")
	}

	if c.httpClient.Jar == nil {
		return nil, errors.New("client has no cookie jar")
	}

	var cookies []*http.Cookie

	for _, cookie := range c.httpClient.Jar.Cookies(s.baseURL) {
		switch cookie.Name {
		case cookiePHPSESSID, cookieSession1:
			cookies = append(cookies, cookie)
		}
	}

	return json.Marshal(serializedClient{
		Cookies: cookies,
	})
}

func (s *ClientSerializer) Unmarshal(data []byte, c *Client) error {
	if c == nil {
		return errors.New("nil client")
	}

	var state serializedClient
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}

	if c.httpClient == nil {
		return errors.New("client has no HTTP client")
	}

	if c.httpClient.Jar == nil {
		jar, err := cookiejar.New(nil)
		if err != nil {
			return err
		}
		c.httpClient.Jar = jar
	}

	c.httpClient.Jar.SetCookies(s.baseURL, state.Cookies)

	return nil
}

var _ Serializer[Client] = (*ClientSerializer)(nil)
