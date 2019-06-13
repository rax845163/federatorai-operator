package probe

import (
	"os"

	"github.com/containers-ai/alameda/pkg/utils/log"
)

var scope = log.RegisterScope("probe", "datahub health probe", 0)

func LivenessProbe(cfg *LivenessProbeConfig) {
	bindAddr := cfg.BindAddr
	err := queryDatahub(bindAddr)
	if err != nil {
		scope.Errorf("Liveess probe failed due to %s", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func ReadinessProbe(cfg *ReadinessProbeConfig) {
	influxdbAddr := cfg.InfluxdbAddr
	prometheusCfg := cfg.PrometheusCfg

	err := pingInfluxdb(influxdbAddr)
	if err != nil {
		scope.Errorf("Readiness probe: ping influxdb failed due to %s", err.Error())
		os.Exit(1)
	}

	err = queryPrometheus(prometheusCfg)
	if err != nil {
		scope.Errorf("Readiness probe: query failed due to %s", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
