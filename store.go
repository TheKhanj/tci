package tci

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Store[T any] struct {
	dir        string
	filename   string
	serializer Serializer[T]
}

type StoreOption[T any] func(*Store[T]) error

func StoreWithDir[T any](dir string) StoreOption[T] {
	return func(s *Store[T]) error {
		if dir == "" {
			return errors.New("store directory is empty")
		}
		s.dir = dir
		return nil
	}
}

func StoreWithFilename[T any](name string) StoreOption[T] {
	return func(s *Store[T]) error {
		if name == "" || filepath.Base(name) != name {
			return errors.New("invalid store filename")
		}
		s.filename = name
		return nil
	}
}

func GetDataDir() (string, error) {
	if os.Geteuid() == 0 {
		return "/var/lib", nil
	}

	if dir := os.Getenv("XDG_STATE_HOME"); dir != "" {
		return dir, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home directory: %w", err)
	}

	return filepath.Join(home, ".local", "state"), nil
}

func NewStore[T any](
	serializer Serializer[T],
	opts ...StoreOption[T],
) (*Store[T], error) {
	if serializer == nil {
		return nil, errors.New("serializer is nil")
	}

	dir, err := GetDataDir()
	if err != nil {
		return nil, err
	}

	s := &Store[T]{
		dir:        filepath.Join(dir, "tci"),
		filename:   "data",
		serializer: serializer,
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Store[T]) Load(ctx context.Context) (*T, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filepath.Join(s.dir, s.filename))
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", s.filename, err)
	}

	var value T
	if err := s.serializer.Unmarshal(data, &value); err != nil {
		return nil, fmt.Errorf("deserialize %s: %w", s.filename, err)
	}

	return &value, nil
}

func (s *Store[T]) Save(ctx context.Context, value *T) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if value == nil {
		return errors.New("value is nil")
	}

	data, err := s.serializer.Marshal(value)
	if err != nil {
		return fmt.Errorf("serialize %s: %w", s.filename, err)
	}

	if err := os.MkdirAll(s.dir, 0o700); err != nil {
		return fmt.Errorf("create store directory: %w", err)
	}

	path := filepath.Join(s.dir, s.filename)

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write %s: %w", s.filename, err)
	}

	return nil
}

func (s *Store[T]) Delete(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	err := os.Remove(filepath.Join(s.dir, s.filename))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("delete %s: %w", s.filename, err)
	}

	return nil
}
