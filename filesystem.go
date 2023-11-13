package filesystem

import (
	"io"
)

var _ Operator = (*Filesystem)(nil)

type Option func(*Filesystem)

func WithConfig(c Config) Option {
	return func(f *Filesystem) {
		f.Config = c
	}
}

type Filesystem struct {
	Adapter Adapter
	Config  Config
}

// NewFileSystem creates a new filesystem.
func NewFileSystem(adapter Adapter, options ...Option) *Filesystem {
	f := &Filesystem{
		Adapter: adapter,
	}

	for _, option := range options {
		option(f)
	}

	return f
}

// FileExists if a file exists.
func (f *Filesystem) FileExists(path string) (bool, error) {
	return f.Adapter.FileExists(path)
}

// DirectoryExists if a directory exists.
func (f *Filesystem) DirectoryExists(path string) (bool, error) {
	return f.Adapter.DirectoryExists(path)
}

// Has if a file exists.
func (f *Filesystem) Has(path string) (bool, error) {
	fileExists, err := f.Adapter.FileExists(path)
	if err != nil {
		return false, err
	}

	directoryExists, err := f.Adapter.DirectoryExists(path)
	if err != nil {
		return false, err
	}

	return fileExists || directoryExists, nil
}

// Write a new file.
func (f *Filesystem) Write(path string, contents io.Reader, config map[string]any) (bool, error) {
	return f.Adapter.Write(
		path,
		contents,
		f.Config.Extend(NewConfig(config)),
	)
}

// Read a file.
func (f *Filesystem) Read(path string) (io.ReadCloser, error) {
	return f.Adapter.Read(path)
}

// Delete the file at a given path.
func (f *Filesystem) Delete(path string) (bool, error) {
	return f.Adapter.Delete(path)
}

// DeleteDir Delete a directory.
func (f *Filesystem) DeleteDir(path string) (bool, error) {
	return f.Adapter.DeleteDir(path)
}

// CreateDir Create a directory.
func (f *Filesystem) CreateDir(path string, config map[string]any) (bool, error) {
	return f.Adapter.CreateDir(path, f.Config.Extend(NewConfig(config)))
}

// MimeType Get the mime-type of a given file.
func (f *Filesystem) MimeType(path string) (string, error) {
	return f.Adapter.MimeType(path)
}

// LastModified Get the last modified time of a file as a UNIX timestamp.
func (f *Filesystem) LastModified(path string) (int64, error) {
	return f.Adapter.LastModified(path)
}

// FileSize Get the file size of a given file.
func (f *Filesystem) FileSize(path string) (int64, error) {
	return f.Adapter.FileSize(path)
}

// ListContents List contents of a directory.
func (f *Filesystem) ListContents(directory string, recursive bool) ([]map[string]any, error) {
	return f.Adapter.ListContents(directory, recursive)
}

// Move a file to a new location.
func (f *Filesystem) Move(source string, destination string) (bool, error) {
	if source == destination {
		return false, ERR_SOURCE_SAME
	}

	return f.Adapter.Move(source, destination)
}

// Copy a file to a new location.
func (f *Filesystem) Copy(source string, destination string) (bool, error) {
	if source == destination {
		return false, ERR_SOURCE_SAME
	}

	return f.Adapter.Copy(source, destination)
}
