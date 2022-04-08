package system_test

import (
	"path/filepath"
	"strings"
	"testing"

	. "github.com/christiangelone/bang/lib/sugar/test"
	"github.com/christiangelone/bang/lib/system"
)

func TestFileSystem(t *testing.T) {

	DescribeFunc(system.NewLocalFileSystem, func() {
		It("should instantiate the fyle system", func() {
			fs := system.NewLocalFileSystem(true)
			AssertThat(fs, ShouldNotBeNil)
			AssertThat(fs, ShouldImplement, (*system.FileSystem)(nil))
		})
	}, t)

	var fs system.FileSystem = system.NewLocalFileSystem(true)
	DescribeStruct(fs, func() {
		DescribeMethod(fs.Get, func() {
			Context("File exists with valid path", func() {
				fileName := "theFile"
				path := "a/path/to"
				pathParts := strings.Split(path, "/")
				for i := 0; i < len(pathParts); i++ {
					partialPath := strings.Join(pathParts[0:i], "/")
					if !fs.Exists(partialPath) {
						//fs.CreateDir(dirPath[:], fileName)
					}
				}

				filePath := filepath.Join(path, fileName)
				data := []byte{0x8, 0x3}
				err := fs.CreateFile(path, fileName, data)
				AssertThat(err, ShouldBeNil)
				It("should get the file", func() {
					file, err := fs.Get(filePath)
					AssertThat(err, ShouldBeNil)
					AssertThat(file.IsDir(), ShouldBeFalse)
					AssertThat(file.Path(), ShouldEqual, path)
					AssertThat(file.Name, ShouldEqual, fileName)
					AssertThat(file.FullPath(), ShouldEqual, filePath)
					fileData, err := file.Read()
					AssertThat(err, ShouldBeNil)
					AssertThat(fileData, ShouldResemble, data)
				})
			})
		})
	}, t)
}
