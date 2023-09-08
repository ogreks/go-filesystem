package qiniu

import (
	"bytes"
	"context"
	"github.com/noOvertimeGroup/go-filesystem"
	"github.com/qiniu/go-sdk/v7/storage"
	"io"
)

type Client struct {
	bucket  *storage.FormUploader
	manager *storage.BucketManager
	upToken string
}

type Storage struct {
	client  *Client
	upToken string
}

func NewStorage(client *Client) filesystem.Storage {
	return &Storage{
		client: client,
	}
}

func (s *Storage) PutFile(ctx context.Context, target string, file io.Reader) error {

	buf := &bytes.Buffer{}
	size, err := buf.ReadFrom(file)
	if err != nil {
		return err
	}
	fileBytes := buf.Bytes()
	data := bytes.NewReader(fileBytes)
	return s.client.bucket.Put(ctx, nil, s.client.upToken, target, data, int64(size), nil)
}

func (s *Storage) GetFile(ctx context.Context, target string) (io.Reader, error) {
	fileInfo, sErr := s.client.manager.Stat("", target)
	b := new(bytes.Buffer)
	b.WriteString(fileInfo.String())
	return b, sErr
	//return io.Reader, sErr
}
