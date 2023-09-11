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

package qiniu

import (
	"bytes"
	"context"
	"io"

	"github.com/noOvertimeGroup/go-filesystem"
	"github.com/qiniu/go-sdk/v7/storage"
)

type Client struct {
	bucket  *storage.FormUploader
	manager *storage.BucketManager
	upToken string
}

type Storage struct {
	client *Client
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
