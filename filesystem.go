package filesystem

import "io"

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

func NewFileSystem(adapter Adapter, options ...Option) *Filesystem {
	f := &Filesystem{
		Adapter: adapter,
	}

	for _, option := range options {
		option(f)
	}

	return f
}

func (f *Filesystem) FileExists(path string) (bool, error) {
	return f.Adapter.FileExists(path)
}

func (f *Filesystem) DirectoryExists(path string) (bool, error) {
	return f.Adapter.DirectoryExists(path)
}

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

func (f *Filesystem) Write(path string, contents io.Reader, config map[string]any) (bool, error) {
	return f.Adapter.Write(
		path,
		contents,
		f.Config.Extend(NewConfig(config)),
	)
}

func (f *Filesystem) Read(path string) (io.ReadCloser, error) {
	return f.Adapter.Read(path)
}

func (f *Filesystem) Delete(path string) (bool, error) {
	return f.Adapter.Delete(path)
}

func (f *Filesystem) DeleteDir(path string) (bool, error) {
	return f.Adapter.DeleteDir(path)
}

func (f *Filesystem) CreateDir(path string, config map[string]any) (bool, error) {
	return f.Adapter.CreateDir(path, f.Config.Extend(NewConfig(config)))
}

func (f *Filesystem) SetVisibility(path string, visibility string) (bool, error) {
	return f.Adapter.SetVisibility(path, visibility)
}

func (f *Filesystem) Visibility(path string) (string, error) {
	return f.Adapter.Visibility(path)
}

func (f *Filesystem) MimeType(path string) (string, error) {
	return f.Adapter.MimeType(path)
}

func (f *Filesystem) LastModified(path string) (int64, error) {
	return f.Adapter.LastModified(path)
}

func (f *Filesystem) FileSize(path string) (int64, error) {
	return f.Adapter.FileSize(path)
}

func (f *Filesystem) ListContents(directory string, recursive bool) ([]map[string]any, error) {
	return f.Adapter.ListContents(directory, recursive)
}

func (f *Filesystem) resolveConfigForMoveAndCopy(config map[string]any) (Config, error) {
	retainVisibility, err := f.Config.Get("retain_visibility", true)
	if err != nil {
		return nil, err
	}

	nc := f.Config.Extend(NewConfig(config))

	if retainVisibility.(bool) {
		visibility, err := f.Visibility(config["visibility"].(string))
		if err != nil {
			return nil, err
		}

		nc.WithDefault("visibility", visibility)
	}

	return nc, nil
}

func (f *Filesystem) Move(source string, destination string, config map[string]any) (bool, error) {
	nc, err := f.resolveConfigForMoveAndCopy(config)
	if err != nil {
		return false, err
	}

	if source == destination {
		return true, nil
	}

	return f.Adapter.Move(source, destination, nc)
}

func (f *Filesystem) Copy(source string, destination string, config map[string]any) (bool, error) {
	nc, err := f.resolveConfigForMoveAndCopy(config)
	if err != nil {
		return false, err
	}

	if source == destination {
		return true, nil
	}

	return f.Adapter.Copy(source, destination, nc)
}
