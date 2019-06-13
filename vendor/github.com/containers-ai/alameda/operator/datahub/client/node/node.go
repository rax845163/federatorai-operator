package node

import (
	"context"

	datahubutils "github.com/containers-ai/alameda/operator/pkg/utils/datahub"
	logUtil "github.com/containers-ai/alameda/pkg/utils/log"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc"
	corev1 "k8s.io/api/core/v1"
)

var (
	scope = logUtil.RegisterScope("datahub node repository", "datahub node repository", 0)
)

// AlamedaNodeRepository creates predicted node to datahub
type AlamedaNodeRepository struct{}

// NewAlamedaNodeRepository return AlamedaNodeRepository instance
func NewAlamedaNodeRepository() *AlamedaNodeRepository {
	return &AlamedaNodeRepository{}
}

// CreateAlamedaNode creates predicted node to datahub
func (repo *AlamedaNodeRepository) CreateAlamedaNode(nodes []*corev1.Node) error {
	retries := 3
	for retry := 1; retry <= retries; retry++ {
		err := repo.createAlamedaNode(nodes)
		if err == nil {
			break
		}
		scope.Debugf("Create Alameda node failed. (%v try)", retry)
		if retry == retries {
			return err
		}
	}
	return nil
}

func (repo *AlamedaNodeRepository) createAlamedaNode(nodes []*corev1.Node) error {
	alamedaNodes := []*datahub_v1alpha1.Node{}
	for _, node := range nodes {

		cpuCores, ok := node.Status.Capacity.Cpu().AsInt64()
		if !ok {
			// TODO: use node.Status.Capacity.Cpu().AsDec()
			scope.Errorf("create alameda node %s to Datahub failed, cannot convert cpu capacity, skip this node", node.GetName())
			return nil
		}

		memoryBytes, ok := node.Status.Capacity.Memory().AsInt64()
		if !ok {
			// TODO: use node.Status.Capacity.Cpu().AsDec()
			scope.Errorf("create alameda node %s to Datahub failed, cannot convert memory capacity, skip this node", node.GetName())
			return nil
		}

		alamedaNodes = append(alamedaNodes, &datahub_v1alpha1.Node{
			Name: node.GetName(),
			Capacity: &datahub_v1alpha1.Capacity{
				CpuCores:    cpuCores,
				MemoryBytes: memoryBytes,
			},
			StartTime: &timestamp.Timestamp{
				Seconds: node.ObjectMeta.GetCreationTimestamp().Unix(),
			},
		})
	}
	req := datahub_v1alpha1.CreateAlamedaNodesRequest{
		AlamedaNodes: alamedaNodes,
	}
	conn, err := grpc.Dial(datahubutils.GetDatahubAddress(), grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		scope.Error(err.Error())
		return err
	}

	datahubServiceClnt := datahub_v1alpha1.NewDatahubServiceClient(conn)
	if reqRes, err := datahubServiceClnt.CreateAlamedaNodes(context.Background(), &req); err != nil {
		scope.Error(reqRes.GetMessage())
		return err
	}
	return nil
}

// DeleteAlamedaNodes delete predicted node from datahub
func (repo *AlamedaNodeRepository) DeleteAlamedaNodes(nodes []*corev1.Node) error {

	conn, err := grpc.Dial(datahubutils.GetDatahubAddress(), grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return errors.Wrapf(err, "delete node from Datahub failed: %s", err.Error())
	}

	alamedaNodes := []*datahub_v1alpha1.Node{}
	for _, node := range nodes {
		alamedaNodes = append(alamedaNodes, &datahub_v1alpha1.Node{
			Name: node.GetName(),
		})
	}
	req := datahub_v1alpha1.DeleteAlamedaNodesRequest{
		AlamedaNodes: alamedaNodes,
	}

	aiServiceClnt := datahub_v1alpha1.NewDatahubServiceClient(conn)
	if resp, err := aiServiceClnt.DeleteAlamedaNodes(context.Background(), &req); err != nil {
		return errors.Wrapf(err, "delete node from Datahub failed: %s", err.Error())
	} else if resp.Code != int32(code.Code_OK) {
		return errors.Errorf("delete node from Datahub failed: receive code: %d, message: %s", resp.Code, resp.Message)
	}
	return nil
}

// ListAlamedaNodes lists nodes to datahub
func (repo *AlamedaNodeRepository) ListAlamedaNodes() ([]*datahub_v1alpha1.Node, error) {
	retries := 3
	alamNodes := []*datahub_v1alpha1.Node{}
	for retry := 1; retry <= retries; retry++ {
		nodes, err := repo.listAlamedaNodes()
		if err == nil {
			alamNodes = nodes
			break
		}
		scope.Debugf("List alameda nodes failed. (%v try)", retry)
		if retry == retries {
			return nil, err
		}
	}
	return alamNodes, nil
}

func (repo *AlamedaNodeRepository) listAlamedaNodes() ([]*datahub_v1alpha1.Node, error) {
	alamNodes := []*datahub_v1alpha1.Node{}
	req := datahub_v1alpha1.ListAlamedaNodesRequest{}
	conn, err := grpc.Dial(datahubutils.GetDatahubAddress(), grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		scope.Error(err.Error())
		return nil, err
	}

	datahubServiceClnt := datahub_v1alpha1.NewDatahubServiceClient(conn)
	if reqRes, err := datahubServiceClnt.ListAlamedaNodes(context.Background(), &req); err != nil {
		if reqRes.Status != nil {
			scope.Error(reqRes.Status.GetMessage())
		}
		return alamNodes, err
	} else {
		alamNodes = reqRes.GetNodes()
	}
	return alamNodes, nil
}
