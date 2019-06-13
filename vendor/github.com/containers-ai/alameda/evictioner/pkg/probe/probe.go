package probe

import (
	"os"

	"github.com/containers-ai/alameda/pkg/utils/log"
)

var scope = log.RegisterScope("probe", "evictioner health probe", 0)

func LivenessProbe(cfg *LivenessProbeConfig) {
	os.Exit(0)
}

func ReadinessProbe(cfg *ReadinessProbeConfig) {
	datahubAddr := cfg.DatahubAddr
	err := queryDatahub(datahubAddr)
	if err != nil {
		scope.Errorf("Readiness probe: query datahub failed due to %s", err.Error())
		os.Exit(1)
	}

	admCtlSvcName := cfg.AdmissionController.SvcName
	admCtlSvcNS := cfg.AdmissionController.Namespace
	admCtlPort := cfg.AdmissionController.Port
	err = queryAdmissionControlerSvc(admCtlSvcName, admCtlSvcNS, admCtlPort)
	if err != nil {
		scope.Errorf("Readiness probe: query admission controller service failed due to %s", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
