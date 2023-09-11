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
	"errors"
	"io"
	"os"
	"path"

	"github.com/noOvertimeGroup/go-filesystem"
)

type Client interface {
	GetDir(filepath string) string
	CreateDir(filepath string) error
	CreateFile(filepath string) error
	Info(filepath string) (os.FileInfo, bool)
	CreateAndGetFileInfo(filepath string, mode int) (*Bucket, error)
	PutFile(target string, file io.Reader) error
	CopyFile(source string, dest string) error
}

type Bucket struct {
	filepath string
	file     *os.File
	fileInfo os.FileInfo
}

// 判断文件文件夹是否存在
func (b *Bucket) Info(filepath string) (os.FileInfo, bool) {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return nil, false
	}
	b.fileInfo = fileInfo
	return fileInfo, true
}

// 创建文件夹
func (b *Bucket) CreateDir(filepath string) error {
	return os.Mkdir(filepath, os.ModePerm)
}

// 判断是否有权限
func (b *Bucket) chmodFile(filepath string, err error) error {
	if os.IsPermission(err) {
		err := os.Chmod(filepath, 0666)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

func (b *Bucket) GetDir(filepath string) string {
	return path.Dir(filepath)
}

// 创建新的文件
func (b *Bucket) createNewFile(filepath string) error {
	_, err := os.Create(filepath)
	return err
}

func (b *Bucket) openFile(filepath string, mode int) (*os.File, error) {
	return os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|mode, 0666)
}

func (b *Bucket) CreateAndGetFileInfo(filepath string, mode int) (*Bucket, error) {
	fileInfo, flag := b.Info(filepath)
	if !flag {
		dir := b.GetDir(filepath)
		_, dirFlag := b.Info(dir)
		if !dirFlag {
			err := b.CreateDir(dir)
			if err != nil {
				return nil, errors.New("文件夹创建失败")
			}
		}
		err := b.createNewFile(filepath)
		if err != nil {
			return nil, errors.New("文件创建失败")
		}
		fileInfo, _ = b.Info(filepath)
	}
	b.fileInfo = fileInfo
	file, err := b.openFile(filepath, mode)
	if err != nil {
		if b.chmodFile(filepath, err) != nil {
			return nil, errors.New("文件无访问权限")
		}
		file, err = b.openFile(filepath, mode)
		if err != nil {
			return nil, errors.New("打开文件失败")
		}
	}
	b.file = file
	return b, nil
}

func (b *Bucket) CreateFile(filepath string) error {
	_, err := os.Create(filepath)
	return err
}

func (b *Bucket) CloseFile() {
	if b.file != nil {
		b.file.Close()
	}
}

func (b *Bucket) PutFile(target string, file io.Reader) error {
	bk1, err := b.CreateAndGetFileInfo(target, os.O_APPEND)
	if err != nil {
		return err
	}
	defer bk1.CloseFile()
	buf := make([]byte, filesystem.BUFFERSIZE)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return errors.New("读取文件失败")
		}
		if n == 0 {
			break
		}
		if _, err := bk1.file.Write(buf[0:n]); err != nil {
			return errors.New("数据写入失败")
		}
	}
	return nil
}

func (b *Bucket) CopyFile(source string, dest string) error {
	bk1, err := b.CreateAndGetFileInfo(source, os.O_APPEND)
	if err != nil {
		return err
	}
	defer bk1.CloseFile()
	bk2, err := b.CreateAndGetFileInfo(dest, os.O_TRUNC)
	if err != nil {
		return err
	}
	defer bk2.CloseFile()
	buf := make([]byte, filesystem.BUFFERSIZE)
	for {
		n, err := bk1.file.Read(buf)
		if err != nil && err != io.EOF {
			return errors.New("读取文件失败")
		}
		if n == 0 {
			break
		}
		if _, err := bk2.file.Write(buf[0:n]); err != nil {
			return errors.New("数据写入失败")
		}
	}
	return nil
}
