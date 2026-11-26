package config

import (
	"path/filepath"
	"runtime"
)

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
