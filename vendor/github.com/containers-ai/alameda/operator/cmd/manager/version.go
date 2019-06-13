package main

import (
	"fmt"
)

func printSoftwareInfo() {
	scope.Infof(fmt.Sprintf("Alameda Version: %s", VERSION))
	scope.Infof(fmt.Sprintf("Alameda Build Time: %s", BUILD_TIME))
	scope.Infof(fmt.Sprintf("Alameda GO Version: %s", GO_VERSION))
}
