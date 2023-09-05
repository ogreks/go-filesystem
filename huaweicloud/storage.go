package huaweicloud

import (
	"bytes"
	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/noOvertimeGroup/go-filesystem"
	"io"
	"path"
	"strings"
)

type Storage struct {
	client *obs.ObsClient
}

func NewStorage(client *obs.ObsClient) filesystem.Storage {
	return &Storage{
		client: client,
	}
}

func (s *Storage) PutFile(target string, file io.Reader) error {
	if !path.IsAbs(target) {
		return errors.New("给定服务路径不是相对路径")
	}

	index := strings.Index(target, "/")
	bucket := target[:index]
	target = target[index:]

	input := &obs.PutObjectInput{}
	input.Bucket = bucket
	input.Key = target
	input.Body = file

	_, err := s.client.PutObject(input)
	if err != nil {
		return err
	}
	// TODO 可以根据返回值进一步判断错误
	return nil
}

func (s *Storage) GetFile(target string) (io.Reader, error) {
	if !path.IsAbs(target) {
		return nil, errors.New("给定服务路径不是相对路径")
	}

	buf := new(bytes.Buffer)
	index := strings.Index(target, "/")
	bucket := target[:index]
	target = target[index:]

	input := &obs.GetObjectInput{}
	input.Bucket = bucket
	input.Key = target

	response, err := s.client.GetObject(input)
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
