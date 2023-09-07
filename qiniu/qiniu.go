package qiniu

import (
	"context"
	"github.com/noOvertimeGroup/go-filesystem"
	"github.com/qiniu/go-sdk/v7/storage"
)

type Storage struct {
	client  *storage.FormUploader
	upToken string
}

func NewFileSystem(bucket *storage.FormUploader, upToken string) filesystem.FileSystem {
	return &Storage{
		client:  bucket,
		upToken: upToken,
	}
}

func (s *Storage) PutFile(ctx context.Context, target string, file string) error {
	err := s.client.PutFile(ctx, nil, s.upToken, target, file, nil)
	if err != nil {
		return err
	}
	return nil
}
