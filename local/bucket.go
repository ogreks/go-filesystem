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
	"os"
	"path"
)

type Bucket struct {
	filepath string
	file     *os.File
	fileInfo os.FileInfo
}

//判断文件文件夹是否存在
func (bucket *Bucket) info(filepath string) (os.FileInfo, bool) {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return nil, false
	}
	bucket.fileInfo = fileInfo
	return fileInfo, true
}

//创建文件夹
func (bucket *Bucket) CreateDir(filepath string) error {
	return os.Mkdir(filepath, os.ModePerm)
}

//判断是否有权限
func (bucket *Bucket) chmodFile(filepath string, err error) error {
	if os.IsPermission(err) {
		err := os.Chmod(filepath, 0666)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

func (bucket *Bucket) GetDir(filepath string) string {
	return path.Dir(filepath)
}

//创建新的文件
func (bucket *Bucket) createNewFile(filepath string) error {
	_, err := os.Create(filepath)
	if err != nil {
		return err
	}
	return err
}

func (bucket *Bucket) openFile(filepath string) (*os.File, error) {
	return os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
}

func (bucket *Bucket) createAndGetFileInfo(filepath string) error {
	fileInfo, flag := bucket.info(filepath)
	if !flag {
		dir := bucket.GetDir(filepath)
		err := bucket.CreateDir(dir)
		if err != nil {
			return errors.New("文件夹创建失败")
		}
		err = bucket.createNewFile(filepath)
		if err != nil {
			return errors.New("文件创建失败")
		}
		fileInfo, _ = bucket.info(filepath)
	}
	bucket.fileInfo = fileInfo
	file, err := bucket.openFile(filepath)
	if err != nil {
		if bucket.chmodFile(filepath, err) != nil {
			return errors.New("文件无访问权限")
		}
		file, err = bucket.openFile(filepath)
		if err != nil {
			return errors.New("打开文件失败")
		}
	}
	bucket.file = file
	return nil
}

func (bucket *Bucket) CreateFile() error {
	return bucket.createAndGetFileInfo(bucket.filepath)
}