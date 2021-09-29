package storage

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
)

type File struct {
	Name string
	URL  string
}

type Storage interface {
	Save(ctx context.Context, data []byte, filename string) (*File, error)
	Delete(ctx context.Context, filename string) error
}

func New(location string, logger zerolog.Logger) Storage {
	store, err := NewFileStorage(logger, location)
	if err != nil {
		fmt.Println("failed to setup file storage")
	}

	return store
}
