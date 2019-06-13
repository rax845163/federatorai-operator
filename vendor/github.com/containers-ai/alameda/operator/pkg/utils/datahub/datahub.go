package datahub

import (
	"os"
)

func GetDatahubAddress() string {
	datahubServer := os.Getenv("ALAMEDA_OPERATOR_DATAHUB_ADDRESS")
	if len(datahubServer) == 0 {
		return "datahub.alameda.svc.cluster.local:50050"
	}
	return datahubServer
}
