package file

import (
	"io"
	"os"
)

type FileHandler struct{}

func NewHandler() *FileHandler {
	return &FileHandler{}
}

func (fh *FileHandler) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (fh *FileHandler) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return io.Copy(dst, src)
}
