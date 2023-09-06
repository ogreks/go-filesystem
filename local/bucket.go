package local

import (
	"errors"
	"os"
	"path"
)

type Bucket struct {
	Filepath string
	File     *os.File
	FileInfo os.FileInfo
}

//判断文件文件夹是否存在
func (bucket *Bucket) IsExistPath(path string) (os.FileInfo, bool) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, false
	}
	bucket.FileInfo = fileInfo
	return fileInfo, true
}

//创建文件夹
func (bucket *Bucket) CreateDir(path string) error {
	return os.Mkdir(path, os.ModePerm)
}

func (bucket *Bucket) CreateFile() error {
	_, flag := bucket.IsExistPath(bucket.Filepath)
	if !flag {
		dir := path.Dir(bucket.Filepath)
		err := bucket.CreateDir(dir)
		if err != nil {
			return errors.New("创建文件加失败")
		}
		file, err := os.Create(bucket.Filepath)
		if err != nil {
			return err
		}
		bucket.File = file
	} else {
		file, err := os.OpenFile(bucket.Filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			if os.IsPermission(err) {
				err := os.Chmod(bucket.Filepath, 0666)
				if err != nil {
					return err
				}
				file, err = os.OpenFile(bucket.Filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			} else {
				return err
			}
		}
		bucket.File = file
	}
	return nil
}
