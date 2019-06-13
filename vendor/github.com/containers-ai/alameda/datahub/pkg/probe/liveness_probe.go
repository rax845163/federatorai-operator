package probe

import (
	"context"
	"fmt"

	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	"google.golang.org/grpc"
)

type LivenessProbeConfig struct {
	BindAddr string
}

func queryDatahub(bindAddr string) error {
	conn, err := grpc.Dial(fmt.Sprintf("localhost%s", bindAddr), grpc.WithInsecure())
	if err != nil {
		return err
	}

	defer conn.Close()
	datahubServiceClnt := datahub_v1alpha1.NewDatahubServiceClient(conn)
	_, err = datahubServiceClnt.ListAlamedaNodes(context.Background(), &datahub_v1alpha1.ListAlamedaNodesRequest{})

	return err
}
