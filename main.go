package main

import (
	"fmt"
	"os"

	"github.com/christiangelone/bang/cmd"
	"github.com/christiangelone/bang/lib/config"
)

var Version string

func main() {
	config.NewIfNotExist()
	if err := cmd.RootCmd(Version).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
