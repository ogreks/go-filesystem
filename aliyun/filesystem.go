package aliyun

import (
	"context"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/noOvertimeGroup/go-filesystem"
	"io"
)

type Storage struct {
	client *oss.Bucket
}

func NewFileSystem(bucket *oss.Bucket) filesystem.FileSystem {
	return &Storage{
		client: bucket,
	}
}

func (s *Storage) PutFile(ctx context.Context, newPath string, file io.ReadCloser) error {
	defer file.Close()

	err := s.client.PutObject(newPath, file)
	if err != nil {
		return err
	}

	return nil
}
