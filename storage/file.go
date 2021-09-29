package storage

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileStorage struct {
	Storage
	root string
	Log  zerolog.Logger
}

func NewFileStorage(logger zerolog.Logger, rootDir string) (*FileStorage, error) {
	path, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if !s.Mode().IsDir() {
		return nil, errors.New("root dir for file storage must be a dir")
	}

	return &FileStorage{
		root: rootDir,
		Log:  logger,
	}, nil
}

func (s *FileStorage) Save(ctx context.Context, data []byte, filename string) (*File, error) {
	url := s.root + "/" + filename
	err := ioutil.WriteFile(url, data, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return &File{
		Name: filename,
		URL:  url,
	}, nil
}

func (s *FileStorage) Delete(ctx context.Context, filename string) error {
	p := filepath.Join(s.root, filename)
	return os.Remove(p)
}
