package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CaptchaProvider func(captcha []byte) (string, error)

func YADCaptchaProvider(captcha []byte) (string, error) {
	f, err := os.CreateTemp("", "captcha-*.png")
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name())

	if _, err := f.Write(captcha); err != nil {
		f.Close()
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", err
	}

	title := fmt.Sprintf("Captcha %s", filepath.Base(f.Name()))

	floating, err := watchI3Window(title)
	if err != nil {
		c := make(chan error, 0)
		floating = c
		close(c)
	}

	cmd := exec.Command(
		"yad",
		"--title="+title,
		"--image="+f.Name(),
		"--entry",
		"--entry-label=Enter CAPTCHA:",
		"--button=Cancel:1",
		"--button=Submit:0",
		"--width=400",
	)

	out, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	defer cmd.Process.Kill()

	if err := <-floating; err != nil {
		return "", err
	}

	data, err := io.ReadAll(out)
	if err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}
