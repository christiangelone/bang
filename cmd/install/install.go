package install

import (
	"context"
	"fmt"
	"os"

	. "github.com/christiangelone/bang/lib/sugar"
	"github.com/christiangelone/bang/lib/ux/print"
	"github.com/christiangelone/bang/lib/ux/progress"
	"github.com/christiangelone/bang/lib/ux/selector"
	"github.com/christiangelone/bang/source"
	"github.com/spf13/cobra"
)

var (
	versionFlag string
)

type Options struct {
	Sources map[source.Type]source.Source
}

type Install struct {
	Options
}

func Cmd(options Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install <binary>",
		Short: "Installs binary",
		Long:  "TODO",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
			install := New(options)
			if err := install.Run(context.Background(), args[0]); err != nil {
				fail := print.Sprint(err.Error(), print.FgRed, "✗")
				print.Bullet(fail)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&versionFlag, "version", "v", "", "version for the binary")

	return cmd
}

func New(options Options) *Install {
	return &Install{
		Options: options,
	}
}

func (i *Install) Run(ctx context.Context, installStr string) error {
	var download source.Download
	if i.isFromSource(installStr) {
		var urlErr error
		download, urlErr = i.GetDownloadFromSourceUrl(ctx, installStr, versionFlag)
		if urlErr != nil {
			return urlErr
		}
	} else {
		downloads, urlErr := i.GetDownloadsFromBinaryName(ctx, installStr, versionFlag)
		if urlErr != nil {
			return urlErr
		}

		if len(downloads) == 0 {
			version := If(versionFlag != "")(print.Sprint("version", print.FgWhite, versionFlag)).
				Else(print.Sprint("version", print.FgWhite, "latest")).(string)
			return fmt.Errorf(
				"no sources found to download %s for %s",
				version, print.Sprint(print.FgYellow, installStr),
			)
		}

		if len(downloads) == 1 {
			download = downloads[0]
		} else {
			selector := selector.New()
			var selectErr error
			download, selectErr = i.SelectDownload(downloads, selector)
			if selectErr != nil {
				return selectErr
			}
		}
	}

	highScoreDownloadUrls, urlErr := download.ChooseHighScoreUrlFromCandidates()
	if urlErr != nil {
		return urlErr
	}

	var downloadUrl string
	if len(highScoreDownloadUrls) == 1 {
		downloadUrl = highScoreDownloadUrls[0]
	} else {
		selector := selector.New()
		var selectErr error
		downloadUrl, selectErr = i.SelectDownloadUrl(highScoreDownloadUrls, selector)
		if selectErr != nil {
			return selectErr
		}
	}

	start := print.Sprint(
		"Getting",
		If(versionFlag == "")(
			print.Sprint("version", print.FgWhite, "latest"),
		).Else(print.Sprint("version", print.FgWhite, versionFlag)),
	)
	start += print.Sprint(" for", print.FgYellow, download.RepoUrl)
	print.Bullet(start)

	downloadProgress := progress.NewProgress()
	fileName, downloadErr := download.Source.Download(ctx, downloadUrl, downloadProgress)
	if downloadErr != nil {
		return downloadErr
	}

	installProgress := progress.NewProgress()
	_, installErr := download.Source.Install(fileName, download.Binary, installProgress)
	if installErr != nil {
		return installErr
	}

	return nil
}

func (i *Install) isFromSource(str string) bool {
	for _, s := range i.Sources {
		if s.IsFromSource(str) {
			return true
		}
	}
	return false
}

func (i *Install) GetDownloadFromSourceUrl(ctx context.Context, url, version string) (source.Download, error) {
	var download source.Download
	for _, s := range i.Sources {
		if s.IsFromSource(url) {
			var err error
			download, err = s.GetDownloadsFromRepoUrl(ctx, url, version)
			if err != nil {
				return source.Download{}, err
			}
			return download, nil
		}
	}

	return download, fmt.Errorf("unable to get download, no source can handle %s", print.Sprint(print.FgYellow, url))
}

func (i *Install) GetDownloadsFromBinaryName(ctx context.Context, binaryName, version string) ([]source.Download, error) {
	downloads := []source.Download{}
	for _, s := range i.Sources {
		someDownloads, err := s.GetDownloadsFromBinaryName(ctx, binaryName, version)
		if err != nil {
			return nil, err
		}
		downloads = append(downloads, someDownloads...)
	}
	return downloads, nil
}

func (i *Install) SelectDownload(downloads []source.Download, selector *selector.Selector) (source.Download, error) {
	index, err := selector.SelectSourceIndex("Select a source", downloads)
	if err != nil {
		selector.SelectFailWith("Quiting while selecting source")
		return source.Download{}, err
	}
	return downloads[index], nil
}

func (i *Install) SelectDownloadUrl(urls []string, selector *selector.Selector) (string, error) {
	index, err := selector.SelectIndex("Select a release url", urls)
	if err != nil {
		selector.SelectFailWith("Quiting while selecting a release url")
		return "", err
	}
	return urls[index], nil
}
