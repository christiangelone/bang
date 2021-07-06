package system

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/christiangelone/bang/lib/file/perm"
)

const bangFolderName = ".bang"

func BangFolderPath() (string, error) {
	dirname, errDir := os.UserHomeDir()
	if errDir != nil {
		return "", errDir
	}

	bangPath := filepath.Join(dirname, bangFolderName)
	if _, err := os.Stat(bangPath); os.IsNotExist(err) {
		bangDirErr := os.Mkdir(bangPath, perm.OS_USER_RWX|perm.OS_GROUP_X|perm.OS_OTH_X)
		if bangDirErr != nil {
			return "", bangDirErr
		}
	}

	return bangPath, nil
}

func Arch() []string {
	arch := runtime.GOARCH
	archValues := []string{arch, strings.ToUpper(arch)}
	switch arch {
	case "amd64":
		archValues = append(archValues, "x86_64", "X86_64")
	case "386":
		archValues = append(archValues, "x86_32", "X86_32")
	default:
	}
	return archValues
}

func Os() string {
	return runtime.GOOS
}
