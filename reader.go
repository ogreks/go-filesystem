package filesystem

import "io"

type Reader interface {
	FileExists(path string) (bool, error)
	DirectoryExists(path string) (bool, error)
	Has(path string) (bool, error)
	Read(path string) (io.ReadCloser, error)
	LastModified(path string) (int64, error)
	FileSize(path string) (int64, error)
	MimeType(path string) (string, error)
	ListContents(directory string, recursive bool) ([]map[string]any, error)
}
