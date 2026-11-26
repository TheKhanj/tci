package tci

import (
	"context"
	"os"
	"testing"
)

func TestClientSerializer(t *testing.T) {
	username := os.Getenv("TCI_USERNAME")
	password := os.Getenv("TCI_PASSWORD")

	if username == "" || password == "" {
		t.Skip("TCI_USERNAME and TCI_PASSWORD must be set")
	}

	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	session, err := NewSession(client, SessionWithCaptchaProvider(YADCaptchaProvider))
	if err != nil {
		t.Fatal(err)
	}

	if err := session.Login(context.Background(), username, password); err != nil {
		t.Fatalf("Login() failed: %v", err)
	}

	authenticated, err := session.Authenticated(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if !authenticated {
		t.Fatal("Authenticated() = false, want true")
	}

	serializer, err := NewClientSerializer()
	if err != nil {
		t.Fatal(err)
	}

	data, err := serializer.Marshal(client)
	if err != nil {
		t.Fatalf("Marshal() failed: %v", err)
	}

	restored, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	if err := serializer.Unmarshal(data, restored); err != nil {
		t.Fatalf("Unmarshal() failed: %v", err)
	}

	session, err = NewSession(restored)
	if err != nil {
		t.Fatal(err)
	}

	authenticated, err = session.Authenticated(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if !authenticated {
		t.Fatal("Authenticated() after deserialization = false, want true")
	}
}
