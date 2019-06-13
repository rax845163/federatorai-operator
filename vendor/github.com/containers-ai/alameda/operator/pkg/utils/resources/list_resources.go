package resources

import (
	"context"
	"strings"

	openshift_apps_v1 "github.com/openshift/api/apps/v1"

	autuscaling "github.com/containers-ai/alameda/operator/pkg/apis/autoscaling/v1alpha1"
	appsapi_v1 "github.com/openshift/api/apps/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ListResources define resource list functions
type ListResources struct {
	client client.Client
}

// NewListResources return ListResources instance
func NewListResources(client client.Client) *ListResources {
	return &ListResources{
		client: client,
	}
}

// ListAllNodes return all nodes in cluster
func (listResources *ListResources) ListAllNodes() ([]*corev1.Node, error) {

	nodes := make([]*corev1.Node, 0)
	nodeList := &corev1.NodeList{}

	if err := listResources.listAllResources(nodeList); err != nil {
		return nodes, err
	}

	for _, node := range nodeList.Items {
		copyNode := node
		nodes = append(nodes, &copyNode)
	}

	return nodes, nil
}

// ListPodsByLabels return pods by labels
func (listResources *ListResources) ListPodsByLabels(labels map[string]string) ([]corev1.Pod, error) {
	podList := &corev1.PodList{}
	if err := listResources.listResourcesByLabels(podList, labels); err != nil {
		return []corev1.Pod{}, err
	}

	return podList.Items, nil
}

// ListDeploymentsByLabels return deployments by labels
func (listResources *ListResources) ListDeploymentsByLabels(labels map[string]string) ([]appsv1.Deployment, error) {
	deploymentList := &appsv1.DeploymentList{}
	if err := listResources.listResourcesByLabels(deploymentList, labels); err != nil {
		return []appsv1.Deployment{}, err
	}

	return deploymentList.Items, nil
}

// ListDeploymentsByNamespaceLabels return deployments by namespace and labels
func (listResources *ListResources) ListDeploymentsByNamespaceLabels(namespace string, labels map[string]string) ([]appsv1.Deployment, error) {
	deploymentList := &appsv1.DeploymentList{}

	if err := listResources.listResourcesByNamespaceLabels(deploymentList, namespace, labels); err != nil {
		return []appsv1.Deployment{}, err
	}

	return deploymentList.Items, nil
}

// ListDeploymentConfigsByNamespaceLabels return deploymentconfigs by namespace and labels
func (listResources *ListResources) ListDeploymentConfigsByNamespaceLabels(namespace string, labels map[string]string) ([]appsapi_v1.DeploymentConfig, error) {
	deploymentConfigList := &appsapi_v1.DeploymentConfigList{}

	if err := listResources.listResourcesByNamespaceLabels(deploymentConfigList, namespace, labels); err != nil {
		return []appsapi_v1.DeploymentConfig{}, err
	}

	return deploymentConfigList.Items, nil
}

// ListDeploymentConfigsByLabels return DeploymentConfigs by labels
func (listResources *ListResources) ListDeploymentConfigsByLabels(labels map[string]string) ([]appsapi_v1.DeploymentConfig, error) {
	deploymentConfigList := &appsapi_v1.DeploymentConfigList{}
	if err := listResources.listResourcesByLabels(deploymentConfigList, labels); err != nil {
		return []appsapi_v1.DeploymentConfig{}, err
	}

	return deploymentConfigList.Items, nil
}

// ListPodsByDeployment return pods by deployment namespace and name
func (listResources *ListResources) ListPodsByDeployment(deployNS, deployName string) ([]corev1.Pod, error) {
	pods := []corev1.Pod{}
	deploymentIns := &appsv1.Deployment{}

	err := listResources.client.Get(context.TODO(), types.NamespacedName{
		Namespace: deployNS,
		Name:      deployName,
	}, deploymentIns)
	if err != nil {
		return pods, err
	}

	replicasetListIns := &appsv1.ReplicaSetList{}
	err = listResources.client.List(context.TODO(),
		client.InNamespace(deployNS),
		replicasetListIns)
	if err != nil {
		return pods, err
	}

	for _, replicasetIns := range replicasetListIns.Items {
		for _, or := range replicasetIns.GetOwnerReferences() {
			if or.Controller != nil && *or.Controller && strings.ToLower(or.Kind) == "deployment" && or.Name == deployName {
				podListIns := &corev1.PodList{}
				err = listResources.client.List(context.TODO(),
					client.InNamespace(deployNS).MatchingLabels(replicasetIns.Spec.Selector.MatchLabels),
					podListIns)
				if err != nil {
					scope.Error(err.Error())
					continue
				}
				pods = append(pods, podListIns.Items...)
			}
		}
	}

	return pods, nil
}

// ListPodsByDeploymentConfig return pods by deployment namespace and name
func (listResources *ListResources) ListPodsByDeploymentConfig(deployConfigNS, deployConfigName string) ([]corev1.Pod, error) {
	pods := []corev1.Pod{}
	deploymentConfigIns := &openshift_apps_v1.DeploymentConfig{}

	err := listResources.client.Get(context.TODO(), types.NamespacedName{
		Namespace: deployConfigNS,
		Name:      deployConfigName,
	}, deploymentConfigIns)
	if err != nil {
		return pods, err
	}

	replicationControllerListIns := &corev1.ReplicationControllerList{}
	err = listResources.client.List(context.TODO(),
		client.InNamespace(deployConfigNS),
		replicationControllerListIns)
	if err != nil {
		return pods, err
	}

	for _, replicationControllerIns := range replicationControllerListIns.Items {
		for _, or := range replicationControllerIns.GetOwnerReferences() {
			if or.Controller != nil && *or.Controller && strings.ToLower(or.Kind) == "deploymentconfig" && or.Name == deployConfigName {
				podListIns := &corev1.PodList{}
				err = listResources.client.List(context.TODO(),
					client.InNamespace(deployConfigNS).MatchingLabels(replicationControllerIns.Spec.Selector),
					podListIns)
				if err != nil {
					scope.Error(err.Error())
					continue
				}
				pods = append(pods, podListIns.Items...)
			}
		}
	}

	return pods, nil
}

// ListAllAlamedaScaler return all AlamedaScaler in cluster
func (listResources *ListResources) ListAllAlamedaScaler() ([]autuscaling.AlamedaScaler, error) {
	alamedaScalerList := &autuscaling.AlamedaScalerList{}
	if err := listResources.listAllResources(alamedaScalerList); err != nil {
		return []autuscaling.AlamedaScaler{}, err
	}
	return alamedaScalerList.Items, nil
}

// ListNamespaceAlamedaScaler return all AlamedaScaler in specific namespace
func (listResources *ListResources) ListNamespaceAlamedaScaler(namespace string) ([]autuscaling.AlamedaScaler, error) {
	alamedaScalerList := &autuscaling.AlamedaScalerList{}
	if err := listResources.listResourcesByNamespace(alamedaScalerList, namespace); err != nil {
		return []autuscaling.AlamedaScaler{}, err
	}
	return alamedaScalerList.Items, nil
}

// ListAlamedaRecommendationOwnedByAlamedaScaler return all AlamedaRecommendation created by input AlamedaScaler
func (listResources *ListResources) ListAlamedaRecommendationOwnedByAlamedaScaler(alamedaScaler *autuscaling.AlamedaScaler) ([]autuscaling.AlamedaRecommendation, error) {

	alamedaRecommendationList := &autuscaling.AlamedaRecommendationList{}

	lbls := make(map[string]string)
	for k, v := range alamedaScaler.GetLabelMapToSetToAlamedaRecommendationLabel() {
		lbls[k] = v
	}

	if err := listResources.listResourcesByNamespaceLabels(alamedaRecommendationList, alamedaScaler.Namespace, lbls); err != nil {
		return []autuscaling.AlamedaRecommendation{}, err
	}

	return alamedaRecommendationList.Items, nil
}

func (listResources *ListResources) listAllResources(resourceList runtime.Object) error {
	if err := listResources.client.List(context.TODO(),
		&client.ListOptions{},
		resourceList); err != nil {
		scope.Error(err.Error())
		return err
	}
	return nil
}

func (listResources *ListResources) listResourcesByNamespace(resourceList runtime.Object, namespace string) error {
	if err := listResources.client.List(context.TODO(),
		&client.ListOptions{
			Namespace: namespace,
		}, resourceList); err != nil {
		scope.Error(err.Error())
		return err
	}
	return nil
}

func (listResources *ListResources) listResourcesByLabels(resourceList runtime.Object, lbls map[string]string) error {
	if err := listResources.client.List(context.TODO(),
		client.MatchingLabels(lbls),
		resourceList); err != nil {
		scope.Error(err.Error())
		return err
	}
	return nil
}

func (listResources *ListResources) listResourcesByNamespaceLabels(resourceList runtime.Object, namespace string, lbls map[string]string) error {
	if err := listResources.client.List(context.TODO(),
		client.InNamespace(namespace).MatchingLabels(lbls),
		resourceList); err != nil {
		scope.Debug(err.Error())
		return err
	}
	return nil
}
