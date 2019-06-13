package probe

import (
	"fmt"
	"os/exec"
)

type LivenessProbeConfig struct {
	AdmissionController *AdmissionController
}

type AdmissionController struct {
	SvcName   string
	Namespace string
	Port      int32
}

func queryWebhookSvc(admCtlSvcName string, admCtlSvcNS string, admCtlPort int32) error {
	svcURL := fmt.Sprintf("https://%s.%s:%s", admCtlSvcName, admCtlSvcNS, fmt.Sprint(admCtlPort))
	curlCmd := exec.Command("curl", "-k", svcURL)

	_, err := curlCmd.CombinedOutput()
	if err != nil {
		return err
	}

	return err
}
