package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultBaseURL = "https://internet.tci.ir"

var defaultCaptchaProvider = YADCaptchaProvider

type SessionOption func(*Session) error

func SessionWithCaptchaProvider(provider CaptchaProvider) SessionOption {
	return func(s *Session) error {
		s.captchaProvider = provider
		return nil
	}
}

func SessionWithBaseURL(baseURL string) SessionOption {
	return func(s *Session) error {
		baseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}

		s.baseURL = baseURL
		return nil
	}
}

func NewSession(client *Client, opts ...SessionOption) (*Session, error) {
	session := Session{
		client:          client,
		captchaProvider: defaultCaptchaProvider,
	}

	for _, opt := range opts {
		if err := opt(&session); err != nil {
			return nil, err
		}
	}

	if session.baseURL == nil {
		if err := SessionWithBaseURL(defaultBaseURL)(&session); err != nil {
			return nil, err
		}
	}

	if err := session.init(); err != nil {
		return nil, err
	}

	return &session, nil
}

type Session struct {
	client          *Client
	baseURL         *url.URL
	captchaProvider CaptchaProvider
}

func (s *Session) Do(req *http.Request) (*http.Response, error) {
	return s.client.Do(req)
}

func (s *Session) Login(
	ctx context.Context,
	username, password string,
) error {
	timestamp := time.Now().Unix()

	captcha, err := s.fetchCaptcha(ctx, timestamp)
	if err != nil {
		return fmt.Errorf("fetch captcha: %w", err)
	}

	answer, err := s.captchaProvider(ctx, captcha)
	if err != nil {
		return fmt.Errorf("solve captcha: %w", err)
	}

	values := url.Values{
		"redirect":     {""},
		"username":     {username},
		"password":     {password},
		"captcha":      {answer},
		"LoginFromWeb": {""},
	}

	target := fmt.Sprintf("%s/panel/login/%d", defaultBaseURL, timestamp)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		target,
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", defaultBaseURL)
	req.Header.Set("Referer", defaultBaseURL+"/panel/")

	resp, err := s.Do(req)
	if err != nil {
		return fmt.Errorf("submit login: %w", err)
	}
	defer resp.Body.Close()

	ok, err := s.Authenticated(ctx)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotAuthenticated
	}

	return nil
}

func (s *Session) Authenticated(ctx context.Context) (bool, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		s.baseURL.String()+"/panel/",
		nil,
	)
	if err != nil {
		return false, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	return !strings.Contains(string(body), "ورود به سامانه"), nil
}

func (s *Session) init() error {
	req, err := http.NewRequest(
		http.MethodGet,
		s.baseURL.String()+"/",
		nil,
	)
	if err != nil {
		return err
	}

	resp, err := s.Do(req)
	if err != nil {
		return fmt.Errorf("initialize session: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func (s *Session) fetchCaptcha(
	ctx context.Context,
	timestamp int64,
) ([]byte, error) {
	target := defaultBaseURL + "/panel/captcha/?r=" + fmt.Sprintf(
		"%d", timestamp,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		target,
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Referer", defaultBaseURL+"/panel/")

	resp, err := s.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected captcha status: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}
