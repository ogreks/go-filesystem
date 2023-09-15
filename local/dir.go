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
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

var _ DirEntryFS = (*DirEntry)(nil)

type DirEntryFS interface {
	Info(dir string) (fs.FileInfo, error)
	Create(dir string) error
	CreateOverlay(dir string) error
	Delete(dir string) error
}

type DirEntry struct {
	driver string
}

func (d *DirEntry) name(dir string) string {
	return path.Join(d.driver, dir)
}

func (d *DirEntry) create(dir string) error {
	err := os.MkdirAll(filepath.Join(d.driver, dir), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (d *DirEntry) Info(dir string) (fs.FileInfo, error) {
	return os.Stat(d.name(dir))
}

// Create use fs create new folder.
func (d *DirEntry) Create(dir string) error {
	return d.create(dir)
}

// CreateOverlay
// if folder exist ignore errors. create folder success
func (d *DirEntry) CreateOverlay(dir string) error {
	err := d.create(d.name(dir))
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

// Delete folder
func (d *DirEntry) Delete(dir string) error {
	return os.Remove(d.name(dir))
}

// NewDirEntry returns DirEntryFS or an error requires an absolute path given
// driver is root path
func NewDirEntry(driver string) (*DirEntry, error) {
	if !filepath.IsAbs(driver) {
		return nil, errors.New("the given path is not an absolute path")
	}

	return &DirEntry{
		driver: driver,
	}, nil
}
