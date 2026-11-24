package tci

import "os"

func TempExecutable(blob []byte) (string, func() error, error) {
	file, err := os.CreateTemp("", "tci-tmp-bin-*")
	if err != nil {
		return "", nil, err
	}

	path := file.Name()
	remove := func() error {
		return os.Remove(path)
	}

	defer func() {
		if err != nil {
			file.Close()
			remove()
		}
	}()

	if _, err = file.Write(blob); err != nil {
		return "", nil, err
	}

	if err = file.Chmod(0o755); err != nil {
		return "", nil, err
	}

	if err = file.Close(); err != nil {
		return "", nil, err
	}

	return path, remove, nil
}
