package use

import (
	"os"
	"strings"

	. "github.com/christiangelone/bang/lib/sugar"
	"github.com/christiangelone/bang/lib/system"
	"github.com/christiangelone/bang/lib/ux/print"
	"github.com/christiangelone/bang/lib/ux/selector"
	"github.com/spf13/cobra"
)

var (
	versionFlag string
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use <binary>",
		Short: "Sets the binary version to use",
		Long:  "TODO",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
			use := New()
			if err := use.Run(args[0]); err != nil {
				fail := print.Sprint(err.Error(), print.FgRed, "✗")
				print.Bullet(fail)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&versionFlag, "version", "v", "", "version to use")

	return cmd
}

type Use struct{}

func New() *Use {
	return &Use{}
}

func (u *Use) Run(binaryName string) error {
	if versionFlag != "" {
		version := If(!strings.Contains(versionFlag, "v"))("v" + versionFlag).Else(versionFlag).(string)
		return u.setVersion(binaryName, version)
	}
	versions, getVersionsErr := u.getAvailableVersions(binaryName)
	if getVersionsErr != nil {
		return getVersionsErr
	}

	selector := selector.New()
	version, selectVersionErr := u.selectVersion(versions, selector)
	if selectVersionErr != nil {
		return selectVersionErr
	}

	return u.setVersion(binaryName, version)
}

func (u *Use) selectVersion(versions []string, selector *selector.Selector) (string, error) {
	index, err := selector.SelectIndex("Select a version", versions)
	if err != nil {
		selector.SelectFailWith("Quiting while selecting version")
		return "", err
	}
	return versions[index], nil
}

func (u *Use) getAvailableVersions(binaryName string) ([]string, error) {
	return system.GetAvailableVersions(binaryName)
}

func (u *Use) setVersion(binaryName, version string) error {
	return system.SetVersion(binaryName, version)
}
