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
	"time"

	"github.com/noOvertimeGroup/go-filesystem"
	qiniuSdk "github.com/qiniu/go-sdk/v7/storage"
)

var _ filesystem.Storage = (*Storage)(nil)

type Storage struct {
	client *qiniuSdk.BucketManager
	domain string // 外链域名
}

func NewStorage(client *qiniuSdk.BucketManager, domain string) *Storage {
	return &Storage{
		client: client,
		domain: domain,
	}
}

func (s *Storage) PutFile(ctx context.Context, target string, file io.Reader) error {
	object, err := filesystem.NewObject(target)
	if err != nil {
		return err
	}

	// qiniu 内部使用了 io.ReadSeeker 导致直接传递 io.Reader 无法使用
	buf := &bytes.Buffer{}
	size, err := buf.ReadFrom(file)
	if err != nil {
		return err
	}

	putPolicy := qiniuSdk.PutPolicy{
		Scope: object.Bucket + ":" + object.Target,
	}
	uploadToken := putPolicy.UploadToken(s.client.Mac)
	from := qiniuSdk.NewFormUploader(s.client.Cfg)
	return from.Put(ctx, nil, uploadToken, object.Target, buf, size, nil)
}

func (s *Storage) GetFile(ctx context.Context, target string) (io.Reader, error) {
	object, err := filesystem.NewObject(target)
	if err != nil {
		return nil, err
	}

	url := qiniuSdk.MakePrivateURL(s.client.Mac, s.domain, object.Target, time.Now().Add(time.Second*10).Unix())
	response, err := s.client.Client.Get(url)
	return response.Body, nil
}

func (s *Storage) Size(ctx context.Context, target string) (int64, error) {
	object, err := filesystem.NewObject(target)
	if err != nil {
		return 0, err
	}

	info, err := s.client.Stat(object.Bucket, object.Target)
	return info.Fsize, err
}
