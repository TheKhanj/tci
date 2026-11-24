package tci

import (
	"context"
	"os"
	"testing"
)

func TestSessionLogin(t *testing.T) {
	username := os.Getenv("TCI_USERNAME")
	password := os.Getenv("TCI_PASSWORD")

	if username == "" || password == "" {
		t.Skip("TCI_USERNAME and TCI_PASSWORD must be set")
	}

	tests := []struct {
		name     string
		username string
		password string
		captcha  CaptchaProvider
		wantAuth bool
	}{
		{
			name:     "valid credentials",
			username: username,
			password: password,
			captcha:  YADCaptchaProvider,
			wantAuth: true,
		},
		{
			name:     "wrong username",
			username: username + "-invalid",
			password: password,
			captcha:  YADCaptchaProvider,
			wantAuth: false,
		},
		{
			name:     "wrong password",
			username: username,
			password: password + "-invalid",
			captcha:  YADCaptchaProvider,
			wantAuth: false,
		},
		{
			name:     "wrong captcha",
			username: username,
			password: password,
			captcha: func(context.Context, []byte) (string, error) {
				return "000000", nil
			},
			wantAuth: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient()
			if err != nil {
				t.Fatal(err)
			}

			session, err := NewSession(
				client,
				SessionWithCaptchaProvider(tt.captcha),
			)
			if err != nil {
				t.Fatal(err)
			}

			err = session.Login(context.Background(), tt.username, tt.password)
			if err != nil && tt.wantAuth {
				t.Fatalf("Login() failed: %v", err)
			}

			if err == nil && !tt.wantAuth {
				t.Fatal("Login() succeeded, expected authentication failure")
			}

			got, err := session.Authenticated(t.Context())
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.wantAuth {
				t.Errorf("Authenticated() = %v, want %v", got, tt.wantAuth)
			}
		})
	}
}
