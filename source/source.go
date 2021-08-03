package source

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/christiangelone/bang/lib/file/perm"
	"github.com/christiangelone/bang/lib/system"
	"github.com/christiangelone/bang/lib/ux/print"
	"github.com/christiangelone/bang/lib/ux/progress"
	"github.com/fatih/color"
)

type Type string

type Downloadable interface {
	Download(ctx context.Context, url string, progress *progress.Progress) (string, error)
}

type Installable interface {
	Install(fileName string, binary Binary, progress *progress.Progress) (string, error)
}

type Source interface {
	Downloadable
	Installable
	IsFromSource(str string) bool
	GetDownloadsFromBinaryName(ctx context.Context, binaryName, version string) ([]Download, error)
	GetDownloadsFromRepoUrl(ctx context.Context, repoUrl, version string) (Download, error)
}

type BaseDownloader struct{}

func (b *BaseDownloader) Download(_ context.Context, _ string, _ *progress.Progress) (string, error) {
	bangPath, err := system.BangFolderPath()
	if err != nil {
		return "", err
	}

	binPath := filepath.Join(bangPath, system.BinFolderName)
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		binDirErr := os.Mkdir(binPath, perm.OS_USER_RWX|perm.OS_GROUP_X|perm.OS_OTH_X)
		if binDirErr != nil {
			return "", binDirErr
		}
	}

	linkPath := filepath.Join(binPath, system.LinkFolderName)
	if _, err := os.Stat(linkPath); os.IsNotExist(err) {
		linkDirErr := os.Mkdir(linkPath, perm.OS_USER_RWX|perm.OS_GROUP_X|perm.OS_OTH_X)
		if linkDirErr != nil {
			return "", linkDirErr
		}
	}

	tmpPath := filepath.Join(bangPath, system.TmpFolderName)
	if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
		tmpDirErr := os.Mkdir(tmpPath, perm.OS_USER_RWX|perm.OS_GROUP_X|perm.OS_OTH_X)
		if tmpDirErr != nil {
			return "", tmpDirErr
		}
	}

	return tmpPath, nil
}

type BaseInstaller struct{}

func (b *BaseInstaller) Install(fileName string, binary Binary, progress *progress.Progress) (string, error) {
	bangPath, err := system.BangFolderPath()
	if err != nil {
		return "", err
	}

	tmpPath := filepath.Join(bangPath, system.TmpFolderName)
	binPath := filepath.Join(bangPath, system.BinFolderName)
	linkPath := filepath.Join(binPath, system.LinkFolderName, binary.Name)

	binTargetPath := filepath.Join(binPath, binary.Name)
	if _, err := os.Stat(binTargetPath); os.IsNotExist(err) {
		targetDirErr := os.Mkdir(binTargetPath, perm.OS_USER_RWX|perm.OS_GROUP_X|perm.OS_OTH_X)
		if targetDirErr != nil {
			return "", targetDirErr
		}
	}

	versionPath := filepath.Join(binTargetPath, binary.Version)
	if _, err := os.Stat(versionPath); os.IsNotExist(err) {
		versionDirErr := os.Mkdir(versionPath, perm.OS_USER_RWX|perm.OS_GROUP_X|perm.OS_OTH_X)
		if versionDirErr != nil {
			return "", versionDirErr
		}
	}

	inputFile, err := ioutil.ReadFile(filepath.Join(tmpPath, fileName))
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(versionPath, binary.Name)
	writeErr := ioutil.WriteFile(filePath, inputFile, perm.OS_USER_RWX|perm.OS_GROUP_X|perm.OS_ALL_X)
	if writeErr != nil {
		return "", writeErr
	}

	if _, linkErr := os.Lstat(linkPath); linkErr == nil {
		if unlinkErr := os.Remove(linkPath); unlinkErr != nil {
			return "", unlinkErr
		}
	}

	linkErr := os.Symlink(filePath, linkPath)
	if linkErr != nil {
		return "", linkErr
	}

	return filepath.Join(binPath, binary.Name, binary.Version), nil
}

type texts struct{}

func (t texts) DownloadText(text string) string {
	return print.Sprint(color.NoColor, "Downloading from", color.FgYellow, text)
}

func (t texts) InstallText(text string) string {
	return print.Sprint(color.NoColor, "Installing", color.FgYellow, text)
}

type Binary struct {
	Name    string
	Version string
}
type Download struct {
	Binary        Binary
	RepoUrl       string
	CandidateUrls []string
	Source        Source
}

type scoreUrl struct {
	score float64
	url   string
}

func (swu Download) ChooseHighScoreUrlFromCandidates() ([]string, error) {
	if len(swu.CandidateUrls) == 0 {
		return nil, errors.New("no candidate url to download binary")
	}

	scoreUrls := []scoreUrl{}
	for _, url := range swu.CandidateUrls {
		scoreUrl := swu.scoreUrl(url)
		scoreUrls = append(scoreUrls, scoreUrl)
	}

	maxScoreUrl := scoreUrls[0]
	for _, scoreUrl := range scoreUrls {
		if scoreUrl.score > maxScoreUrl.score {
			maxScoreUrl = scoreUrl
		}
	}

	allMaxScoreUrls := []string{}
	for _, scoreUrl := range scoreUrls {
		if scoreUrl.score == maxScoreUrl.score {
			allMaxScoreUrls = append(allMaxScoreUrls, scoreUrl.url)
		}
	}

	return allMaxScoreUrls, nil
}

func (swu Download) scoreUrl(url string) scoreUrl {
	toLowerUrl := strings.ToLower(url)
	score := 0.0

	if strings.Contains(toLowerUrl, "tar") || strings.Contains(toLowerUrl, "gz") || strings.Contains(toLowerUrl, "zip") {
		score -= 0.5
	}

	if strings.Contains(toLowerUrl, swu.Binary.Name) {
		score++
	}

	os := system.Os()
	if strings.Contains(toLowerUrl, os) {
		score++
	}

	archValues := system.Arch()
	for _, arch := range archValues {
		if strings.Contains(toLowerUrl, arch) {
			score++
			break
		}
	}

	return scoreUrl{
		score: score,
		url:   url,
	}
}
