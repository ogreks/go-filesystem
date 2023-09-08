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

type Storage struct {
	client *oss.Bucket
}

func NewStorage(bucket *oss.Bucket) filesystem.Storage {
	return &Storage{
		client: bucket,
	}
}

func (s *Storage) PutFile(ctx context.Context, target string, file io.Reader) error {
	err := s.client.PutObject(target, file)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetFile(ctx context.Context, target string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	response, err := s.client.GetObject(target)
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
