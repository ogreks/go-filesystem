package local

import (
	"os"
	"path"
)

type Bucket struct {
	filepath string
	file     *os.File
	fileInfo os.FileInfo
}

//判断文件文件夹是否存在
func (bucket *Bucket) IsExistPath(path string) (os.FileInfo, bool) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, false
	}
	bucket.fileInfo = fileInfo
	return fileInfo, true
}

//创建文件夹
func (bucket *Bucket) CreateDir(path string) error {
	return os.Mkdir(path, os.ModePerm)
}

//创建新的文件
func (bucket *Bucket) createNewFile() error {
	_, err := os.Create(bucket.filepath)
	if err != nil {
		return err
	}
	return err
}

//判断是否有权限
func (bucket *Bucket) chmodPerm(err error) error {
	if os.IsPermission(err) {
		err := os.Chmod(bucket.filepath, 0666)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

func (bucket *Bucket) CreateFile() error {
	_, flag := bucket.IsExistPath(bucket.filepath)
	if !flag {
		dir := path.Dir(bucket.filepath)
		err := bucket.CreateDir(dir)
		if err != nil {
			return err
		}
		err = bucket.createNewFile()
		if err != nil {
			return err
		}
	}
	file, err := os.OpenFile(bucket.filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil && bucket.chmodPerm(err) != nil {
		return err
	}
	file, err = os.OpenFile(bucket.filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	bucket.file = file
	return nil
}
