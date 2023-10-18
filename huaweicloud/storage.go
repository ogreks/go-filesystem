// Copyright (c) 2023 noOvertimeGroup
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package huaweicloud

import (
	"bytes"
	"context"
	"io"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/noOvertimeGroup/go-filesystem"
)

var _ filesystem.Storage = (*Storage)(nil)

type Storage struct {
	client *obs.ObsClient
}

func NewStorage(client *obs.ObsClient) *Storage {
	return &Storage{
		client: client,
	}
}

func (s *Storage) PutFile(ctx context.Context, target string, file io.Reader) error {
	object, err := filesystem.NewObject(target)
	if err != nil {
		return err
	}

	input := &obs.PutObjectInput{}
	input.Bucket = object.Bucket
	input.Key = object.Target
	input.Body = file

	_, err = s.client.PutObject(input)
	return err
}

func (s *Storage) GetFile(ctx context.Context, target string) (io.Reader, error) {
	object, err := filesystem.NewObject(target)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	input := &obs.GetObjectInput{}
	input.Bucket = object.Bucket
	input.Key = object.Target

	response, err := s.client.GetObject(input)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	_, err = io.Copy(buf, response.Body)
	return buf, err
}

func (s *Storage) Size(ctx context.Context, target string) (int64, error) {
	f, err := s.GetFile(ctx, target)
	if err != nil {
		return 0, err
	}

	return int64(f.(*bytes.Buffer).Len()), nil
}
