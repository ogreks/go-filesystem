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
	"io"
	"os"
	"time"
)

// File represents a file in the filesystem.
type File interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Writer
	io.WriterAt

	Name() string
	Readdir(count int) ([]os.FileInfo, error)
	ReadDirNames(n int) ([]string, error)
	Stat() (os.FileInfo, error)
	Sync() error
	Truncate(size int64) error
	WriteString(s string) (ret int, err error)
}

// FS is the filesystem interface.
// Any simulated or real filesystem should implement this interface.
type FS interface {
	Create(name string) (File, error)
	Mkdir(name string, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Open(name string) (File, error)
	OpenFile(name string, flag int, perm os.FileMode) (File, error)
	Remove(name string) error
	RemoveAll(path string) error
	Rename(oldName, newName string) error
	Name() string
	Chmod(name string, mode os.FileMode) error
	Chown(name string, uid, gid int) error
	ChtTimes(name string, atime time.Time, mtime time.Time) error
}
