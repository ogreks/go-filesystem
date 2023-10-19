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
	"errors"
	"fmt"
	"io"

	"github.com/noOvertimeGroup/go-filesystem"
	"github.com/qiniu/go-sdk/v7/storage"
)

var _ filesystem.Storage = (*Storage)(nil)

// TODO 此结构体不应该存在
type Client struct {
	bucket  *storage.FormUploader
	manager *storage.BucketManager

	upToken string
}

type Storage struct {
	client *Client
}

func NewStorage(client *Client) *Storage {
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

	bucketName := ctx.Value("bucketName")
	bucketName1, ok := bucketName.(string)
	if !ok {
		return nil, errors.New("bucketName 必须是string 类型")
	}
	fileInfo, sErr := s.client.manager.Stat(bucketName1, target)
	fmt.Print(fileInfo.String())
	fmt.Print("-------------")
	b := new(bytes.Buffer)
	b.WriteString(fileInfo.String())
	return b, sErr
	//return io.Reader, sErr
}

func (s *Storage) GetFile1(ctx context.Context, target string) {
	bucketName := ctx.Value("bucketName")
	bucketName1, ok := bucketName.(string)

	key := "文件保存的 key.jpg"
	bucket := "对象所在的 bucket"
	NewBucketManager := s.client.manager

	domain := "https://image.example.com"
	key := "这是一个测试文件.jpg"
	publicAccessURL := storage.MakePublicURL(domain, key)
	fmt.Println(publicAccessURL)

}

func (s *Storage) Size(ctx context.Context, target string) (int64, error) {
	bucketName := ctx.Value("bucketName")
	bucketName1, ok := bucketName.(string)
	if !ok {
		return 0, errors.New("bucketName 必须是string 类型")
	}
	fileInfo, err := s.client.manager.Stat(bucketName1, target)
	if err != nil {
		return 0, err
	}
	return fileInfo.Fsize, nil
}
