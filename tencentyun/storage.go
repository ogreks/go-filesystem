package tencentyun

import (
	"bytes"
	"context"
	"errors"
	"github.com/noOvertimeGroup/go-filesystem"
	"github.com/noOvertimeGroup/go-filesystem/internal/errs"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"io/fs"
)

type Storage struct {
	client *cos.Client
}

func NewStorage(client *cos.Client) filesystem.Storage {
	return &Storage{
		client: client,
	}
}

func (s *Storage) PutFile(ctx context.Context, target string, file fs.File) error {
	fileInfo, err := file.Stat()
	if err != nil {
		if errors.Is(err, fs.ErrClosed) {
			return errs.ErrFileClose
		}
		return err
	}

	if fileInfo.Size() > filesystem.FileLimitSize {
		return errs.ErrFileLimit
	}
	// TODO return http.Response handle error
	_, err = s.client.Object.Put(ctx, target, file, &cos.ObjectPutOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetFile(ctx context.Context, target string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	response, err := s.client.Object.Get(ctx, target, &cos.ObjectGetOptions{})
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// TODO logging set log
		}
	}(response.Body)

	_, err = io.Copy(buf, response.Body)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
