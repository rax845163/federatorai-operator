package recommendation

import (
	"github.com/containers-ai/alameda/datahub/pkg/repository/influxdb"
)

const (
	// Container is container measurement
	Container  influxdb.Measurement = "container"
	Controller influxdb.Measurement = "controller"
)
