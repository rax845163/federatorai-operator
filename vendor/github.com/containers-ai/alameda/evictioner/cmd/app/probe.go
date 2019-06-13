package app

import (
	"os"

	"github.com/containers-ai/alameda/evictioner/pkg/probe"
	"github.com/containers-ai/alameda/pkg/utils/kubernetes"
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
		Short: "probe alameda evictioner",
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
	ProbeCmd.Flags().StringVar(&probeType, "type", PROBE_TYPE_READINESS, "The probe type for evicionter.")
}

func startProbing() {
	if probeType == PROBE_TYPE_LIVENESS {
		probe.LivenessProbe(&probe.LivenessProbeConfig{})
	} else if probeType == PROBE_TYPE_READINESS {
		runngingNS := kubernetes.GetRunningNamespace()
		probe.ReadinessProbe(&probe.ReadinessProbeConfig{
			DatahubAddr: config.Datahub.Address,
			AdmissionController: &probe.AdmissionController{
				SvcName:   config.AdmCtr.SvcName,
				Namespace: runngingNS,
				Port:      config.AdmCtr.SvcPort,
			},
		})
	} else {
		scope.Errorf("Probe type does not supports %s, please try %s or %s.", probeType, PROBE_TYPE_LIVENESS, PROBE_TYPE_READINESS)
		os.Exit(1)
	}
}
