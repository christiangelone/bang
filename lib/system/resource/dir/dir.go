package dir

import (
	"fmt"
	"github.com/christiangelone/bang/lib/system/resource/file"
	"path/filepath"
)

type Props struct {
	Name  string
	Path  string
}

type Dir struct {
	name  string
	path  string
	files []file.File
}

func New(props Props) *Dir {
	return &Dir{
		name:  props.Name,
		path:  props.Path,
		files: nil,
	}
}

func (d *Dir) Read() ([]byte, error) {
	return nil, fmt.Errorf("you cannot read from a directory (dir: %s)", d.FullPath())
}

func (d *Dir) IsDir() bool {
	return false
}

func (d *Dir) Path() string {
	return d.path
}

func (d *Dir) Name() string {
	return d.name
}

func (d *Dir) FullPath() string {
	return filepath.Join(d.path, d.name)
}
