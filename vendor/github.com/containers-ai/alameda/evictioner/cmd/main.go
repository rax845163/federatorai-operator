package main

import (
	"github.com/containers-ai/alameda/cmd/app"
	evictioner_app "github.com/containers-ai/alameda/evictioner/cmd/app"
)

var (
	// VERSION is sofeware version
	VERSION string
	// BUILD_TIME is build time
	BUILD_TIME string
	// GO_VERSION is go version
	GO_VERSION string
)

func init() {
	setSoftwareInfo()
}

func main() {
	evictioner_app.RootCmd.Execute()
}

func setSoftwareInfo() {
	app.VERSION = VERSION
	app.BUILD_TIME = BUILD_TIME
	app.GO_VERSION = GO_VERSION
	app.PRODUCT_NAME = "evictioner"
}
