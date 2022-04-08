package system

import (
	"fmt"
	"github.com/christiangelone/bang/lib/system/resource"
	"github.com/christiangelone/bang/lib/system/resource/dir"
	"github.com/christiangelone/bang/lib/system/resource/file"
	"reflect"
)

type LocalFileSystem struct {
	inMemory bool
	memory   inMemoryFileSystem
}

func (fs *LocalFileSystem) GetStructName() string {
	return fmt.Sprintf("%+v", reflect.TypeOf(fs).Elem())
}

func (fs *LocalFileSystem) Exists(filePath string) bool {
	if fs.inMemory {
		return fs.memory.Exists(filePath)
	}
	return false
}

func (fs *LocalFileSystem) Get(filePath string) (resource.Resource, error) {
	if fs.inMemory {
		return fs.memory.Get(filePath)
	}
	panic("not implemented")
}

func (fs *LocalFileSystem) CreateDir(path, dirName string) error {
	d := dir.New(dir.Props{
		Name:  dirName,
		Path:  path,
	})

	if fs.inMemory {
		return fs.memory.CreateDir(d)
	}
	panic("not implemented")
}

func (fs *LocalFileSystem) CreateFile(path, fileName string, data []byte) error {
	f := file.New(file.Props{
		Name: fileName,
		Path: path,
		Data: data,
	})

	if fs.inMemory {
		return fs.memory.CreateFile(f)
	}
	panic("not implemented")
}
