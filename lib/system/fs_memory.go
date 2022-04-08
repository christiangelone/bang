package system

import (
	"fmt"
	"github.com/christiangelone/bang/lib/system/resource"
	"github.com/christiangelone/bang/lib/system/resource/dir"
	"github.com/christiangelone/bang/lib/system/resource/file"
	"strings"
)

type inMemoryFileSystem struct{}

var fsMap = map[string]resource.Resource{}

func (fs *inMemoryFileSystem) validPath(path string) bool {
	pathParts := strings.Split(path, "/")
	for i := 0; i < len(pathParts); i++ {
		partialPath := strings.Join(pathParts[0:i], "/")
		if _, ok := fsMap[partialPath]; !ok {
			return false
		}
	}
	return true
}

func (fs *inMemoryFileSystem) Exists(filePath string) bool {
	_, ok := fsMap[filePath]
	return ok
}

func (fs *inMemoryFileSystem) Get(filePath string) (resource.Resource, error) {
	r, ok := fsMap[filePath]
	if ok {
		return r, nil
	}
	return nil, fmt.Errorf("file '%s' not found", filePath)
}

func (fs *inMemoryFileSystem) CreateFile(f *file.File) error {
	if fs.validPath(f.FullPath()) {
		fsMap[f.FullPath()] = f
		return nil
	}
	return fmt.Errorf("invalid path '%s'", f.FullPath())
}

func (fs *inMemoryFileSystem) CreateDir(d *dir.Dir) error {
	if fs.validPath(d.FullPath()) {
		fsMap[d.FullPath()] = d
		return nil
	}
	return fmt.Errorf("invalid path '%s'", d.FullPath())
}
