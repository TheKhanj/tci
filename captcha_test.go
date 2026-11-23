package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestYADCaptchaProvider(t *testing.T) {
	if _, err := exec.LookPath("yad"); err != nil {
		t.Skip("yad not installed")
	}

	if os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" {
		t.Skip("no graphical display available")
	}

	image, err := os.ReadFile("captchas/0001.png")
	if err != nil {
		t.Fatal(err)
	}

	got, err := YADCaptchaProvider(t.Context(), image)
	if err != nil {
		t.Fatalf("YADCaptchaProvider() error: %v", err)
	}

	if got == "" {
		t.Fatal("YADCaptchaProvider() returned empty value")
	}

	want, err := os.ReadFile("labels/0001")
	if err != nil {
		t.Fatal(err)
	}

	if got != strings.TrimSpace(string(want)) {
		t.Errorf("YADCaptchaProvider() = %q, want %q", got, strings.TrimSpace(string(want)))
	}
}
