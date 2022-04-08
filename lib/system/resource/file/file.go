package file

import (
	"path/filepath"
)

type Props struct {
	Name  string
	Path  string
	Data  []byte
}

type File struct {
	name  string
	path  string
	data  []byte
}

func New(props Props) *File {
	return &File{
		name:  props.Name,
		path:  props.Path,
		data:  props.Data,
	}
}

func (f *File) Read() ([]byte, error) {
	return f.data, nil
}

func (f *File) IsDir() bool {
	return false
}

func (f *File) Path() string {
	return f.path
}

func (f *File) Name() string {
	return f.name
}

func (f *File) FullPath() string {
	return filepath.Join(f.path, f.name)
}
