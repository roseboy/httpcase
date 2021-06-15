package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime/debug"
)

var (
	GitVersion = "dev"
	GitCommit  = ""
	GoVersion  = ""
	BuildDate  = ""
	Platform   = ""
	BuiltBy    = ""
)

type versionCmd struct {
	cmd  *cobra.Command
	opts versionOpts
}

type versionOpts struct {
}

func newVersionCmd() *versionCmd {
	root := &versionCmd{}

	cmd := &cobra.Command{
		Use:   "version",
		Short: "print HttpCase version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(buildVersion())
		},
	}

	root.cmd = cmd
	return root
}
func buildVersion() string {
	result := GitVersion
	if GitCommit != "" {
		result = fmt.Sprintf("%s\nGitCommit: %s", result, GitCommit)
	}
	if GoVersion != "" {
		result = fmt.Sprintf("%s\nGoVersion: %s", result, GoVersion)
	}
	if Platform != "" {
		result = fmt.Sprintf("%s\nPlatform: %s", result, Platform)
	}
	if BuildDate != "" {
		result = fmt.Sprintf("%s\nBuildDate: %s", result, BuildDate)
	}
	if BuiltBy != "" {
		result = fmt.Sprintf("%s\nBuiltBy: %s", result, BuiltBy)
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Sum != "" {
		result = fmt.Sprintf("%s\nmodule version: %s, Checksum: %s", result, info.Main.Version, info.Main.Sum)
	}
	return result
}
