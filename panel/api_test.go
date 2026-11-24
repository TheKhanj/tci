package panel

import (
	"context"
	"os"
	"testing"

	"github.com/thekhanj/tci"
)

func TestApiGet(t *testing.T) {
	username := os.Getenv("TCI_USERNAME")
	password := os.Getenv("TCI_PASSWORD")

	if username == "" || password == "" {
		t.Skip("TCI_USERNAME and TCI_PASSWORD must be set")
	}

	client, err := tci.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	session, err := tci.NewSession(client)
	if err != nil {
		t.Fatal(err)
	}

	if err := session.Login(context.Background(), username, password); err != nil {
		t.Fatalf("Login() failed: %v", err)
	}

	api, err := NewApi()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := api.Close(); err != nil {
			t.Errorf("Api.Close() failed: %v", err)
		}
	}()

	got, err := api.Get(context.Background(), session)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	t.Logf("Panel: %+v", got)
}
