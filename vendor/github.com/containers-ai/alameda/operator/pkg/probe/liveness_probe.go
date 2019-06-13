package probe

import (
	"fmt"
	"os/exec"
)

type LivenessProbeConfig struct {
	ValidationSvc *ValidationSvc
}

type ValidationSvc struct {
	SvcName string
	SvcNS   string
	SvcPort int32
}

func queryWebhookSvc(svcName string, svcNS string, port int32) error {
	svcURL := fmt.Sprintf("https://%s.%s:%s", svcName, svcNS, fmt.Sprint(port))
	curlCmd := exec.Command("curl", "-k", svcURL)

	_, err := curlCmd.CombinedOutput()
	if err != nil {
		return err
	}

	return err
}
