package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
)

type i3WindowEvent struct {
	Change    string `json:"change"`
	Container struct {
		ID               int64 `json:"id"`
		WindowProperties struct {
			Title string `json:"title"`
		} `json:"window_properties"`
	} `json:"container"`
}

func watchI3Window(title string) (<-chan error, error) {
	cmd := exec.Command(
		"i3-msg", "-t", "subscribe", "-m", `["window"]`,
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	done := make(chan error, 1)

	go func() {
		defer cmd.Process.Kill()

		scanner := bufio.NewScanner(stdout)

		for scanner.Scan() {
			var event i3WindowEvent
			if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
				continue
			}

			if event.Change != "new" ||
				event.Container.WindowProperties.Title != title {
				continue
			}

			err := exec.Command(
				"i3-msg",
				fmt.Sprintf(
					"[con_id=%d] floating enable",
					event.Container.ID,
				),
			).Run()

			done <- err
			return
		}

		done <- scanner.Err()
	}()

	return done, nil
}
