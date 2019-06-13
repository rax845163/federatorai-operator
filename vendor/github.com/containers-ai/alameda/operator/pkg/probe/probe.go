package probe

import (
	"os"

	"github.com/containers-ai/alameda/pkg/utils/log"
)

var scope = log.RegisterScope("probe", "datahub health probe", 0)

func LivenessProbe(cfg *LivenessProbeConfig) {
	svcName := cfg.ValidationSvc.SvcName
	svcNS := cfg.ValidationSvc.SvcNS
	svcPort := cfg.ValidationSvc.SvcPort
	err := queryWebhookSvc(svcName, svcNS, svcPort)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func ReadinessProbe(cfg *ReadinessProbeConfig) {
	datahubAddr := cfg.DatahubAddr
	err := queryDatahub(datahubAddr)
	if err != nil {
		scope.Errorf("Readiness probe: query datahub failed due to %s", err.Error())
		os.Exit(1)
	}

	err = queryWebhookSrv(cfg.ValidationSrvPort)
	if err != nil {
		scope.Errorf("Readiness probe: query validation webhook server failed due to %s", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
