package app

import (
	"github.com/containers-ai/alameda/cmd/app"
	"github.com/spf13/cobra"
)

var (
	VersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Display the alameda admission-controller version",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			app.PrintSoftwareVer()
		},
	}
)
