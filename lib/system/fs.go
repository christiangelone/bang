package system

import (
	"github.com/christiangelone/bang/lib/meta"
	"github.com/christiangelone/bang/lib/system/resource"
)

type FileSystem interface {
	meta.MetaStruct
	Exists(path string) bool
	Get(filePath string) (resource.Resource, error)
	CreateFile(path, fileName string, data []byte) error
	CreateDir(path, dirName string) error
}

func NewLocalFileSystem(inMemory bool) *LocalFileSystem {
	return &LocalFileSystem{
		inMemory: inMemory,
	}
}
