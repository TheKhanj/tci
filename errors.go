package tci

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidCaptcha     = errors.New("invalid captcha")
	ErrCaptchaProvider    = errors.New("captcha provider not configured")
	ErrNotAuthenticated   = errors.New("not authenticated")
)
