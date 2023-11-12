package filesystem

import "io"

type Adapter interface {
	FileExists(path string) (bool, error)
	DirectoryExists(path string) (bool, error)
	Write(path string, contents io.Reader, config Config) (bool, error)
	Read(path string) (io.ReadCloser, error)
	Delete(path string) (bool, error)
	DeleteDir(path string) (bool, error)
	CreateDir(path string, config Config) (bool, error)
	SetVisibility(path string, visibility string) (bool, error)
	Visibility(path string) (string, error)
	MimeType(path string) (string, error)
	LastModified(path string) (int64, error)
	FileSize(path string) (int64, error)
	ListContents(directory string, recursive bool) ([]map[string]any, error)
	Move(source string, destination string, config Config) (bool, error)
	Copy(source string, destination string, config Config) (bool, error)
}
