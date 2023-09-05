package local

import (
	"context"
	"io"
)

type LocalStorage struct {
	client LocalClient
}

func (localStorage *LocalStorage) NewFileSystem(file string) *LocalStorage {
	return localStorage
}

func (localStorage *LocalStorage) PutFile(ctx context.Context, newPath string, file io.ReadCloser) error {
	defer file.Close()
	return nil
}
