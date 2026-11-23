package main

import (
	"context"
	"net/http"
)

type Credentials struct {
	Username string
	Password string
}

var defaultCaptchaProvider = YADCaptchaProvider

type SessionOption func(*Session)

func SessionWithCaptchaProvider(provider CaptchaProvider) SessionOption {
	return func(c *Session) {
		c.captchaProvider = provider
	}
}

func NewSession(client *Client, opts ...SessionOption) (*Session, error) {
	session := Session{
		client:          client,
		captchaProvider: defaultCaptchaProvider,
	}

	for _, opt := range opts {
		opt(&session)
	}

	if err := session.init(); err != nil {
		return nil, err
	}

	return &session, nil
}

type Session struct {
	client          *Client
	captchaProvider CaptchaProvider
}

func (s *Session) Do(req *http.Request) (*http.Response, error) {
	return s.client.Do(req)
}

func (s *Session) Login(
	ctx context.Context,
	credentials Credentials,
) error {
	// Implement the TCI login sequence here.
	return nil
}

func (s *Session) Authenticated() bool {
	// Implement authentication state verification here.
	return false
}

func (s *Session) init() error {
}
