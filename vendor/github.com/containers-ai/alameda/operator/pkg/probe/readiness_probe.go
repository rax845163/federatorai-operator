package probe

import (
	"context"
	"fmt"
	"os/exec"

	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	"google.golang.org/grpc"
)

type ReadinessProbeConfig struct {
	ValidationSrvPort int32
	DatahubAddr       string
}

func queryDatahub(datahubAddr string) error {
	conn, err := grpc.Dial(datahubAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	defer conn.Close()
	datahubServiceClnt := datahub_v1alpha1.NewDatahubServiceClient(conn)
	_, err = datahubServiceClnt.ListAlamedaNodes(context.Background(), &datahub_v1alpha1.ListAlamedaNodesRequest{})
	if err != nil {
		return err
	}

	return err
}

func queryWebhookSrv(port int32) error {

	svcURL := fmt.Sprintf("https://localhost:%s", fmt.Sprint(port))
	curlCmd := exec.Command("curl", "-k", svcURL)

	_, err := curlCmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}
