package tci

import "os"

func Executable(blob []byte) (path string, cleanup func(), err error) {
	path = ""
	cleanup = func() {}
	err = nil

	file, err := os.CreateTemp("", "tci-tmp-bin-*")
	if err != nil {
		return path, cleanup, err
	}

	_, err = file.Write(blob)
	if err != nil {
		file.Close()
		os.Remove(file.Name())
		return path, cleanup, err
	}

	err = file.Chmod(0o755)
	file.Close()
	if err != nil {
		os.Remove(file.Name())
		return path, cleanup, err
	}

	return file.Name(), func() {
		os.Remove(file.Name())
	}, nil
}
