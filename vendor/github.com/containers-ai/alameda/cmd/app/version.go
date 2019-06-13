package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// VERSION is sofeware version
	VERSION string
	// BUILD_TIME is build time
	BUILD_TIME string
	// GO_VERSION is go version
	GO_VERSION string
	// PRODUCT_NAME is product name
	PRODUCT_NAME string
)

var (
	VersionCmd = &cobra.Command{
		Use:   "version",
		Short: "show " + PRODUCT_NAME + "version",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			PrintSoftwareVer()
		},
	}
)

// PrintSoftwareVer shows software version
func PrintSoftwareVer() {
	fmt.Println(fmt.Sprintf("%s Version: %s", PRODUCT_NAME, VERSION))
	fmt.Println(fmt.Sprintf("%s Build Time: %s", PRODUCT_NAME, BUILD_TIME))
	fmt.Println(fmt.Sprintf("%s GO Version: %s", PRODUCT_NAME, GO_VERSION))
}
