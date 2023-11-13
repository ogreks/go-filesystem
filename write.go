package filesystem

import "io"

type Writer interface {
	Write(path string, contents io.Reader, config map[string]any) (bool, error)
	Delete(path string) (bool, error)
	DeleteDir(path string) (bool, error)
	CreateDir(path string, config map[string]any) (bool, error)
	Move(source string, destination string) (bool, error)
	Copy(source string, destination string) (bool, error)
}
