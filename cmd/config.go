package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/thekhanj/tci"
	"github.com/urfave/cli/v3"
)

func ConfigFlag(path *string) cli.Flag {
	return &cli.StringFlag{
		Name:        "config",
		Aliases:     []string{"c"},
		Usage:       "Path to configuration file",
		Value:       DefaultConfigPath(),
		Destination: path,
	}
}

func DefaultConfigPath() string {
	candidates := []string{}

	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(cwd, "tci.json"))
	}

	if dir, err := os.UserConfigDir(); err == nil {
		candidates = append(candidates, filepath.Join(dir, "tci", "tci.json"))
	}

	if runtime.GOOS != "windows" {
		candidates = append(candidates, "/etc/tci/tci.json")
	}

	for _, path := range candidates {
		info, err := os.Stat(path)
		if err == nil && info.Mode().IsRegular() {
			return path
		}
		if err != nil && !os.IsNotExist(err) {
			return path
		}
	}

	if dir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(dir, "tci", "tci.json")
	}

	return "tci.json"
}

func ReadConfig(path string) (tci.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return tci.Config{}, fmt.Errorf("read config %s: %w", path, err)
	}

	var config tci.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return tci.Config{}, fmt.Errorf("parse config %s: %w", path, err)
	}

	return config, nil
}
