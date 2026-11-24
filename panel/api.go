package panel

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"

	"github.com/thekhanj/tci"
)

//go:embed html2json
var html2json []byte

type Api struct {
	baseURL *url.URL
	binary  string
	cleanup func() error
}

type ApiOption func(*Api) error

func ApiWithBaseURL(baseURL string) ApiOption {
	return func(a *Api) error {
		u, err := url.Parse(baseURL)
		if err != nil {
			return err
		}

		a.baseURL = u
		return nil
	}
}

func NewApi(opts ...ApiOption) (*Api, error) {
	api := Api{}

	for _, opt := range opts {
		if err := opt(&api); err != nil {
			return nil, err
		}
	}

	if api.baseURL == nil {
		if err := ApiWithBaseURL(tci.DefaultBaseURL)(&api); err != nil {
			return nil, err
		}
	}

	binary, cleanup, err := tci.TempExecutable(html2json)
	if err != nil {
		return nil, err
	}

	api.binary = binary
	api.cleanup = cleanup

	return &api, nil
}

func (a *Api) Close() error {
	if a.cleanup == nil {
		return nil
	}

	return a.cleanup()
}

func (a *Api) Get(ctx context.Context, session *tci.Session) (Panel, error) {
	var panel Panel

	u := *a.baseURL
	u.Path = "/panel/"

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		u.String(),
		nil,
	)
	if err != nil {
		return panel, err
	}

	res, err := session.Do(req)
	if err != nil {
		return panel, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return panel, fmt.Errorf("GET %s: %s", req.URL, res.Status)
	}

	html, err := io.ReadAll(res.Body)
	if err != nil {
		return panel, err
	}

	cmd := exec.CommandContext(ctx, a.binary)
	cmd.Stdin = bytes.NewReader(html)

	jsonData, err := cmd.Output()
	if err != nil {
		return panel, err
	}

	if err := json.Unmarshal(jsonData, &panel); err != nil {
		return panel, err
	}

	return panel, nil
}
