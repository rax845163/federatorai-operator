package app

import (
	"github.com/spf13/cobra"
)

var (
	configurationFilePath string
)

var RootCmd = &cobra.Command{
	Use:              "admission-controller",
	Short:            "alameda admission-controller",
	Long:             "",
	TraverseChildren: true,
}

func init() {
	RootCmd.AddCommand(RunCmd)
	RootCmd.AddCommand(VersionCmd)
	RootCmd.AddCommand(ProbeCmd)

	RootCmd.PersistentFlags().StringVar(&configurationFilePath, "config", "/etc/alameda/admission-controller/admission-controller.yml", "File path to admission-controller coniguration")
}
