package local

import (
	"context"
	"github.com/noOvertimeGroup/go-filesystem"
	"io"
)

type Storage struct {
	client *Bucket
}

func NewFileSystem(path string) filesystem.FileSystem {
	return &Storage{
		client: &Bucket{filepath: path},
	}
}

func (localStorage *Storage) PutFile(ctx context.Context, newPath string, file io.ReadCloser) error {
	defer file.Close()
	return nil
}
