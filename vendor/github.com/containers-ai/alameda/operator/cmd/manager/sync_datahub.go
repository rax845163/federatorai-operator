package main

import (
	"fmt"
	"time"

	datahub_node "github.com/containers-ai/alameda/operator/datahub/client/node"
	datahub_pod "github.com/containers-ai/alameda/operator/datahub/client/pod"
	"github.com/containers-ai/alameda/operator/pkg/utils/resources"
	alamutils "github.com/containers-ai/alameda/pkg/utils"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func registerNodes(client client.Client, retryInterval int64) {
	scope.Info("Start registering nodes to datahub.")
	for {
		time.Sleep(time.Duration(retryInterval) * time.Second)
		if err := startRegisteringNodes(client); err == nil {
			scope.Info("Register nodes to datahub successfully.")
			break
		} else {
			scope.Errorf("Register nodes to datahub failed due to %s.", err.Error())
		}
	}
	scope.Info("Registering nodes to datahub is done.")
}

func startRegisteringNodes(client client.Client) error {
	listResources := resources.NewListResources(client)
	nodes, err := listResources.ListAllNodes()

	if err != nil {
		return fmt.Errorf("register nodes to Datahub failed: %s", err.Error())
	}

	if len(nodes) == 0 {
		return fmt.Errorf("No nodes found to register to datahub")
	}

	scope.Infof(fmt.Sprintf("%v nodes found in cluster.", len(nodes)))
	datahubNodeRepo := datahub_node.NewAlamedaNodeRepository()
	err = datahubNodeRepo.CreateAlamedaNode(nodes)
	if err != nil {
		return err
	}
	nodesToDel := []*corev1.Node{}
	alamNodes, err := datahubNodeRepo.ListAlamedaNodes()
	for _, alamNode := range alamNodes {
		toDel := true
		for _, node := range nodes {
			if node.GetName() == alamNode.GetName() {
				toDel = false
				break
			}
		}
		if !toDel {
			continue
		}
		delNode := &corev1.Node{}
		delNode.SetName(alamNode.GetName())
		nodesToDel = append(nodesToDel, delNode)
	}

	if len(nodesToDel) > 0 {
		scope.Debugf("Nodes removed from datahub. %s", alamutils.InterfaceToString(nodesToDel))
		err := datahubNodeRepo.DeleteAlamedaNodes(nodesToDel)
		if err != nil {
			return err
		}
	}

	return err
}

func syncAlamedaPodsWithDatahub(client client.Client, retryInterval int64) {
	scope.Info("Start syncing alameda pods to datahub.")
	for {
		if err := startSyncingAlamedaPodsWithDatahubSuccess(client); err == nil {
			scope.Info("Sync alameda pod with datahub successfully.")
			break
		} else {
			scope.Errorf("Sync alameda pod with datahub failed due to %s", err.Error())
		}
		time.Sleep(time.Duration(retryInterval) * time.Second)
	}
	scope.Info("Syncing alameda pods to datahub is done.")
}

func startSyncingAlamedaPodsWithDatahubSuccess(client client.Client) error {
	datahubPodRepo := datahub_pod.NewPodRepository()
	alamedaPods, err := datahubPodRepo.ListAlamedaPods()
	if err != nil {
		return fmt.Errorf("Sync alameda pod with datahub failed: %s", err.Error())
	}
	podsNeedToRm := []*datahub_v1alpha1.Pod{}
	getResource := resources.NewGetResource(client)

	for _, alamedaPod := range alamedaPods {
		namespacedName := alamedaPod.GetNamespacedName()
		alamPodNS := namespacedName.GetNamespace()
		alamPodName := namespacedName.GetName()
		_, err := getResource.GetPod(alamPodNS, alamPodName)
		if err != nil && k8sErrors.IsNotFound(err) {
			podsNeedToRm = append(podsNeedToRm, alamedaPod)
			continue
		} else if err != nil {
			return fmt.Errorf("Get pod (%s/%s) failed while sync alameda pod with datahub. (%s)", alamPodNS, alamPodName, err.Error())
		}

		alamedaScaler := alamedaPod.GetAlamedaScaler()
		alamScalerNS := alamedaScaler.GetNamespace()
		alamScalerName := alamedaScaler.GetName()
		_, err = getResource.GetAlamedaScaler(alamScalerNS, alamScalerName)
		if err != nil && k8sErrors.IsNotFound(err) {
			podsNeedToRm = append(podsNeedToRm, alamedaPod)
			continue
		} else if err != nil {
			return fmt.Errorf("Get alameda scaler (%s/%s) failed while sync alameda pod with datahub. (%s)", alamedaScaler, alamScalerName, err.Error())
		}
	}

	if len(podsNeedToRm) > 0 {
		err := datahubPodRepo.DeletePods(podsNeedToRm)
		return err
	}
	return nil
}
