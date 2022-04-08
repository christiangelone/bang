package source

import (
	"context"
	"errors"
	"fmt"
	"github.com/christiangelone/bang/lib/system/resource/file"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"

	"github.com/christiangelone/bang/lib/config"
	. "github.com/christiangelone/bang/lib/sugar"
	"github.com/christiangelone/bang/lib/ux/progress"
)

const GithubSourceType = Type("github")

type GithubOptions struct {
	HttpClient interface {
		Do(req *http.Request) (*http.Response, error)
	}
	FileHandler interface {
		OpenFile(name string, flag int, perm os.FileMode) (*os.File, error)
		Copy(dst io.Writer, src io.Reader) (written int64, err error)
	}
}

type Github struct {
	texts
	BaseDownloader
	BaseInstaller
	GithubOptions
	client *github.Client
}

func NewGithub(options GithubOptions) *Github {
	cleanOptions := completeDefaults(options)
	gitHubClient := getGitHubClient()
	return &Github{
		texts:          texts{},
		BaseDownloader: BaseDownloader{},
		BaseInstaller:  BaseInstaller{},
		GithubOptions:  cleanOptions,
		client:         gitHubClient,
	}
}

func getGitHubClient() *github.Client {
	config := config.GetGitHubConfig()
	var tc *http.Client
	if config.GitHubAuthToken != "" {
		tc = oauth2.NewClient(
			context.Background(),
			oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: config.GitHubAuthToken},
			),
		)
	}
	return github.NewClient(tc)
}

func completeDefaults(options GithubOptions) GithubOptions {
	if options.FileHandler == nil {
		options.FileHandler = file.NewHandler()
	}
	if options.HttpClient == nil {
		options.HttpClient = http.DefaultClient
	}
	return options
}

func (g *Github) IsFromSource(str string) bool {
	_, _, ok := g.getOwnerAndRepoNameFrom(str)
	return ok
}

func (g *Github) GetDownloadsFromRepoUrl(ctx context.Context, repoUrl, version string) (Download, error) {
	owner, repoName, ok := g.getOwnerAndRepoNameFrom(repoUrl)

	if !ok {
		return Download{}, errors.New("malformed github repo url")
	}

	repo := &github.Repository{
		Owner: &github.User{
			Login: &owner,
		},
		Name:     &repoName,
		CloneURL: &repoUrl,
	}

	return g.getDownloadFromRepo(ctx, repo, version)
}

func (g *Github) GetDownloadsFromBinaryName(ctx context.Context, binaryName, version string) ([]Download, error) {

	options := &github.SearchOptions{
		Sort: "full_name",
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 5,
		},
	}

	result, _, err := g.client.Search.Repositories(ctx, binaryName, options)
	if err != nil {
		return nil, err
	}

	downloads := []Download{}
	for _, repo := range result.Repositories {
		download, err := g.getDownloadFromRepo(ctx, repo, version)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				continue
			}
			fmt.Printf("Error while fetching latest release => %s\n", err.Error())
			return nil, err
		}
		downloads = append(downloads, download)
	}

	return downloads, nil
}

func (g *Github) getOwnerAndRepoNameFrom(str string) (string, string, bool) {
	re := regexp.MustCompile("^(https?://)?github.com/([A-Za-z0-9]{1,})/([A-Za-z0-9-_]{1,})$")
	parts := re.FindStringSubmatch(str)

	if parts == nil {
		return "", "", false
	}

	return parts[2], parts[3], true
}

func (g *Github) getDownloadFromRepo(ctx context.Context, repo *github.Repository, version string) (Download, error) {
	owner := repo.GetOwner().GetLogin()
	name := repo.GetName()

	var release *github.RepositoryRelease
	var err error
	if version == "" {
		release, _, err = g.client.Repositories.GetLatestRelease(ctx, owner, name)
	} else {
		tag := If(!strings.Contains(version, "v"))("v" + version).Else(version).(string)
		release, _, err = g.client.Repositories.GetReleaseByTag(ctx, owner, name, tag)
	}
	if err != nil {
		return Download{}, err
	}

	url := strings.Replace(repo.GetCloneURL(), "https://", "", 1)
	url = strings.Replace(url, ".git", "", 1)
	swu := Download{
		RepoUrl: url,
		Binary: Binary{
			Name:    name,
			Version: release.GetTagName(),
		},
		Source: g,
	}

	if len(release.Assets) == 0 {
		swu.CandidateUrls = append(swu.CandidateUrls, release.GetTarballURL())
	} else {
		for _, asset := range release.Assets {
			swu.CandidateUrls = append(swu.CandidateUrls, asset.GetBrowserDownloadURL())
		}
	}

	return swu, nil
}

func (g *Github) Download(ctx context.Context, url string, progress *progress.Progress) (string, error) {
	progress.Start(g.texts.DownloadText(url))

	path, pathErr := g.BaseDownloader.Download(ctx, url, progress)
	if pathErr != nil {
		progress.StopFailWith("Error while creating tmp dir => " + pathErr.Error())
		return "", pathErr
	}

	req, reqErr := http.NewRequestWithContext(ctx, "GET", url, nil)
	if reqErr != nil {
		progress.StopFailWith("Error while building request => " + reqErr.Error())
		return "", reqErr
	}

	resp, respErr := g.HttpClient.Do(req)
	if respErr != nil {
		progress.StopFailWith("Error while requesting url => " + respErr.Error())
		return "", respErr
	}
	defer resp.Body.Close()

	contentDisposition := resp.Header.Get("Content-Disposition")
	_, params, dispErr := mime.ParseMediaType(contentDisposition)
	if dispErr != nil {
		progress.StopFailWith("Error while extracting filename from Content-Disposition => " + dispErr.Error())
		return "", dispErr
	}
	fileName := params["filename"]

	file, fileErr := g.FileHandler.OpenFile(filepath.Join(path, fileName), os.O_CREATE|os.O_WRONLY, 0644)
	if fileErr != nil {
		progress.StopFailWith("Error while creating binary file => " + fileErr.Error())
		return "", fileErr
	}
	defer file.Close()

	_, copyErr := g.FileHandler.Copy(io.MultiWriter(file, progress.GetBar(resp.ContentLength)), resp.Body)
	if copyErr != nil {
		progress.StopFailWith("Error while copying bytes to binary file => " + copyErr.Error())
		return "", copyErr
	}

	progress.StopSuccess()
	return fileName, nil
}

func (g *Github) Install(fileName string, binary Binary, progress *progress.Progress) (string, error) {
	progress.Start(g.texts.InstallText(binary.Name))

	_, pathErr := g.BaseInstaller.Install(fileName, binary, progress)
	if pathErr != nil {
		progress.StopFailWith("Error while installing => " + pathErr.Error())
		return "", pathErr
	}

	progress.StopSuccess()
	return fileName, nil
}
