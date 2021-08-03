package system

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func GetAvailableVersions(binaryName string) ([]string, error) {
	bangPath, err := BangFolderPath()
	if err != nil {
		return nil, err
	}

	binPath := filepath.Join(bangPath, BinFolderName)

	binTargetPath := filepath.Join(binPath, binaryName)
	if _, err := os.Stat(binTargetPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s has no versions installed", binaryName)
	}

	files, filesErr := ioutil.ReadDir(binTargetPath)
	if filesErr != nil {
		return nil, filesErr
	}

	versions := []string{}

	for _, f := range files {
		if f.IsDir() {
			versions = append(versions, f.Name())
		}
	}

	return versions, nil
}

func SetVersion(binaryName, version string) error {
	bangPath, err := BangFolderPath()
	if err != nil {
		return nil
	}

	binPath := filepath.Join(bangPath, BinFolderName)
	linkPath := filepath.Join(binPath, LinkFolderName, binaryName)

	binTargetPath := filepath.Join(binPath, binaryName)
	if _, err := os.Stat(binTargetPath); os.IsNotExist(err) {
		return fmt.Errorf("%s has no versions installed", binaryName)
	}

	versionPath := filepath.Join(binTargetPath, version)
	if _, err := os.Stat(versionPath); os.IsNotExist(err) {
		return fmt.Errorf("version %s for %s is not installed", version, binaryName)
	}

	filePath := filepath.Join(versionPath, binaryName)

	if _, linkErr := os.Lstat(linkPath); linkErr == nil {
		if unlinkErr := os.Remove(linkPath); unlinkErr != nil {
			return unlinkErr
		}
	}

	linkErr := os.Symlink(filePath, linkPath)
	if linkErr != nil {
		return linkErr
	}

	return nil
}
