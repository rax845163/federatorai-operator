package pod

import (
	"context"

	datahubutils "github.com/containers-ai/alameda/operator/pkg/utils/datahub"
	logUtil "github.com/containers-ai/alameda/pkg/utils/log"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc"
)

var (
	scope = logUtil.RegisterScope("datahub pod repository", "datahub pod repository", 0)
)

// PodRepository creates predicted pod to datahub
type PodRepository struct{}

// NewPodRepository return PodRepository instance
func NewPodRepository() *PodRepository {
	return &PodRepository{}
}

func (repo *PodRepository) ListAlamedaPods() ([]*datahub_v1alpha1.Pod, error) {
	alamedaPods := []*datahub_v1alpha1.Pod{}
	conn, err := grpc.Dial(datahubutils.GetDatahubAddress(), grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return nil, errors.Wrapf(err, "list Alameda pods from Datahub failed: %s", err.Error())
	}

	req := datahub_v1alpha1.ListAlamedaPodsRequest{
		Kind: datahub_v1alpha1.Kind_POD,
	}
	aiServiceClnt := datahub_v1alpha1.NewDatahubServiceClient(conn)
	if resp, err := aiServiceClnt.ListAlamedaPods(context.Background(), &req); err != nil {
		return alamedaPods, errors.Wrapf(err, "list Alameda pods from Datahub failed: %s", err.Error())
	} else if resp.Status != nil && resp.Status.Code != int32(code.Code_OK) {
		return alamedaPods, errors.Errorf("list Alameda pods from Datahub failed: receive code: %d, message: %s", resp.Status.Code, resp.Status.Message)
	} else {
		alamedaPods = resp.GetPods()
	}
	return alamedaPods, nil
}

// DeletePods delete pods from datahub
func (repo *PodRepository) DeletePods(pods []*datahub_v1alpha1.Pod) error {

	conn, err := grpc.Dial(datahubutils.GetDatahubAddress(), grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return errors.Wrapf(err, "delete pods from Datahub failed: %s", err.Error())
	}

	req := datahub_v1alpha1.DeletePodsRequest{
		Pods: pods,
	}

	aiServiceClnt := datahub_v1alpha1.NewDatahubServiceClient(conn)
	if resp, err := aiServiceClnt.DeletePods(context.Background(), &req); err != nil {
		return errors.Wrapf(err, "delete pods from Datahub failed: %s", err.Error())
	} else if resp.Code != int32(code.Code_OK) {
		return errors.Errorf("delete pods from Datahub failed: receive code: %d, message: %s", resp.Code, resp.Message)
	}
	return nil
}
