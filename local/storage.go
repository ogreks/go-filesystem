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

package local

import (
	"bytes"
	"context"
	"io"
	"os"

	"github.com/noOvertimeGroup/go-filesystem"
)

type Storage struct {
	client Client
}

func NewStorage(c Client) filesystem.Storage {
	return &Storage{
		client: c,
	}
}

//复制文件，类似上传文件
func (s *Storage) PutFile(tx context.Context, target string, file io.Reader) error {
	return s.client.PutFile(target, file)
}

func (s *Storage) GetFile(ctx context.Context, target string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	b, err := s.client.CreateAndGetFileInfo(target, os.O_APPEND)
	defer b.CloseFile()
	_, err = io.Copy(buf, b.file)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
