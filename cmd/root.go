package cmd

import (
	"os"

	"github.com/christiangelone/bang/cmd/completion"
	"github.com/christiangelone/bang/cmd/install"
	. "github.com/christiangelone/bang/lib/sugar"
	"github.com/christiangelone/bang/source"
	"github.com/spf13/cobra"
)

func RootCmd(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Version: If(version != "")(version).Else("N/A").(string),
		Use:     "bang",
		Short:   "Bang a binary manager that really shoots",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}
	return cmdList(rootCmd)(
		completion.Cmd(),
		install.Cmd(install.Options{
			Sources: map[source.Type]source.Source{
				source.GithubSourceType: source.NewGithub(source.GithubOptions{}),
			},
		}),
	)
}

func cmdList(rootCmd *cobra.Command) func(cmds ...*cobra.Command) *cobra.Command {
	return func(cmds ...*cobra.Command) *cobra.Command {
		for _, cmd := range cmds {
			rootCmd.AddCommand(cmd)
		}
		return rootCmd
	}
}
