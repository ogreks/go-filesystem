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

package aliyun

import (
	"bytes"
	"context"
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/noOvertimeGroup/go-filesystem"
)

var _ filesystem.Storage = (*Storage)(nil)

type Storage struct {
	client *oss.Client
}

func NewStorage(client *oss.Client) *Storage {
	return &Storage{
		client: client,
	}
}

// setBucket set client bucket return *oss.Bucket,filesystem.Object,error
func (s *Storage) setBucket(target string) (bucket *oss.Bucket, object filesystem.Object, err error) {
	object, err = filesystem.NewObject(target)
	if err != nil {
		return
	}

	bucket, err = s.client.Bucket(object.Bucket)
	return
}

func (s *Storage) PutFile(ctx context.Context, target string, file io.Reader) error {
	b, object, err := s.setBucket(target)
	if err != nil {
		return err
	}

	return b.PutObject(object.Target, file)
}

func (s *Storage) GetFile(ctx context.Context, target string) (io.Reader, error) {
	b, object, err := s.setBucket(target)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	response, err := b.GetObject(object.Target)
	if err != nil {
		return nil, err
	}

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			// TODO logging set log
			return
		}
	}(response)

	_, err = io.Copy(buf, response)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// Size GetFile read file bytes length return
func (s *Storage) Size(ctx context.Context, target string) (int64, error) {
	f, err := s.GetFile(ctx, target)
	if err != nil {
		return 0, err
	}

	return int64(f.(*bytes.Buffer).Len()), nil
}
