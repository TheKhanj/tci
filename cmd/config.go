package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/thekhanj/tci"
)

func ReadConfig() (tci.Config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return tci.Config{}, err
	}

	paths := []string{
		"tci.json",
		filepath.Join(configDir, "tci", "tci.json"),
		"/etc/tci/tci.json",
	}

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return tci.Config{}, fmt.Errorf("read config %s: %w", path, err)
		}

		var config tci.Config
		if err := json.Unmarshal(data, &config); err != nil {
			return tci.Config{}, fmt.Errorf("parse config %s: %w", path, err)
		}

		return config, nil
	}

	return tci.Config{}, fmt.Errorf("no configuration file found")
}
