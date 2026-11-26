package tci

import (
	"errors"
	"net/url"
)

func WithBaseURL[T any](baseURL string, set func(*T, *url.URL)) func(*T) error {
	return func(v *T) error {
		u, err := url.Parse(baseURL)
		if err != nil {
			return err
		}

		if u.Scheme != "http" && u.Scheme != "https" {
			return errors.New("invalid base URL scheme")
		}

		if u.Host == "" {
			return errors.New("base URL host is required")
		}

		set(v, u)
		return nil
	}
}
