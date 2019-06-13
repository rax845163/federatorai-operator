package utils

import (
	"encoding/json"
	"os"
)

// InterfaceToString encodes interface to string
func InterfaceToString(data interface{}) string {
	if configBin, err := json.Marshal(data); err != nil {
		return ""
	} else {
		return string(configBin)
	}
}

// GetRunningNamespace retrieves value from env NAMESPACE_NAME
func GetRunningNamespace() string {
	return os.Getenv("NAMESPACE_NAME")
}

// GetRunnningPodName retrieves value from env POD_NAME
func GetRunnningPodName() string {
	return os.Getenv("POD_NAME")
}
