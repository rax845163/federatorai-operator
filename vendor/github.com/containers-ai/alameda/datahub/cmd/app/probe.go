package app

import (
	"os"

	"github.com/containers-ai/alameda/datahub/pkg/probe"
	"github.com/spf13/cobra"
)

const (
	PROBE_TYPE_READINESS = "readiness"
	PROBE_TYPE_LIVENESS  = "liveness"
)

var (
	probeType string

	ProbeCmd = &cobra.Command{
		Use:   "probe",
		Short: "probe alameda datahub server",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			initConfig()
			initLogger()
			setLoggerScopesWithConfig(*config.Log)
			startProbing()
		},
	}
)

func init() {
	parseProbeFlag()
}

func parseProbeFlag() {
	ProbeCmd.Flags().StringVar(&probeType, "type", PROBE_TYPE_READINESS, "The probe type for datahub.")
}

func startProbing() {
	if probeType == PROBE_TYPE_LIVENESS {
		probe.LivenessProbe(&probe.LivenessProbeConfig{
			BindAddr: config.BindAddress,
		})
	} else if probeType == PROBE_TYPE_READINESS {
		probe.ReadinessProbe(&probe.ReadinessProbeConfig{
			InfluxdbAddr:  config.InfluxDB.Address,
			PrometheusCfg: config.Prometheus,
		})
	} else {
		scope.Errorf("Probe type does not supports %s, please try %s or %s.", probeType, PROBE_TYPE_LIVENESS, PROBE_TYPE_READINESS)
		os.Exit(1)
	}
}
