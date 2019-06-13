package app

import (
	"flag"
	"os"

	"github.com/containers-ai/alameda/admission-controller/pkg/probe"
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
		Short: "probe alameda admission-controller server",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			flag.Parse()

			initConfig()
			initLog()

			if !config.Enable {
				os.Exit(0)
			}
			startProbing()
			return nil
		},
	}
)

func init() {
	ProbeCmd.Flags().StringVar(&probeType, "type", PROBE_TYPE_READINESS, "The probe type for admission controller.")
}

func startProbing() {
	if probeType == PROBE_TYPE_LIVENESS {
		runngingNS := kubernetes.GetRunningNamespace()
		probe.LivenessProbe(&probe.LivenessProbeConfig{
			AdmissionController: &probe.AdmissionController{
				SvcName:   config.Service.Name,
				Namespace: runngingNS,
				Port:      config.Service.Port,
			},
		})
	} else if probeType == PROBE_TYPE_READINESS {
		probe.ReadinessProbe(&probe.ReadinessProbeConfig{
			DatahubAddr:   config.Datahub.Address,
			AdmCtrSrvPort: config.Port,
		})
	} else {
		scope.Errorf("Probe type does not supports %s, please try %s or %s.", probeType, PROBE_TYPE_LIVENESS, PROBE_TYPE_READINESS)
		os.Exit(1)
	}
}
