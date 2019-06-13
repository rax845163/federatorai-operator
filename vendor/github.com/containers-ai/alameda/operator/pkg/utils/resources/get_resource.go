package resources

import (
	"context"

	autuscaling "github.com/containers-ai/alameda/operator/pkg/apis/autoscaling/v1alpha1"
	logUtil "github.com/containers-ai/alameda/pkg/utils/log"
	appsapi_v1 "github.com/openshift/api/apps/v1"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	scope = logUtil.RegisterScope("resource_utils", "resource utils", 0)
)

// GetResource define resource list functions
type GetResource struct {
	client.Client
}

// NewGetResource return GetResource instance
func NewGetResource(client client.Client) *GetResource {
	return &GetResource{
		client,
	}
}

// GetPod returns pod
func (getResource *GetResource) GetPod(namespace, name string) (*corev1.Pod, error) {
	pod := &corev1.Pod{}
	err := getResource.getResource(pod, namespace, name)
	return pod, err
}

// GetReplicationController returns replicationController
func (getResource *GetResource) GetReplicationController(namespace, name string) (*corev1.ReplicationController, error) {
	replicationController := &corev1.ReplicationController{}
	err := getResource.getResource(replicationController, namespace, name)
	return replicationController, err
}

// GetReplicaSet returns replicaSet
func (getResource *GetResource) GetReplicaSet(namespace, name string) (*appsv1.ReplicaSet, error) {
	replicaSet := &appsv1.ReplicaSet{}
	err := getResource.getResource(replicaSet, namespace, name)
	return replicaSet, err
}

// GetDeploymentConfig returns deploymentconfig
func (getResource *GetResource) GetDeploymentConfig(namespace, name string) (*appsapi_v1.DeploymentConfig, error) {
	deploymentConfig := &appsapi_v1.DeploymentConfig{}
	err := getResource.getResource(deploymentConfig, namespace, name)
	return deploymentConfig, err
}

// GetDeployment returns deployment
func (getResource *GetResource) GetDeployment(namespace, name string) (*appsv1.Deployment, error) {
	deployment := &appsv1.Deployment{}
	err := getResource.getResource(deployment, namespace, name)
	return deployment, err
}

// GetAlamedaScaler return alamedascaler
func (getResource *GetResource) GetAlamedaScaler(namespace, name string) (*autuscaling.AlamedaScaler, error) {
	alamedaScaler := &autuscaling.AlamedaScaler{}
	err := getResource.getResource(alamedaScaler, namespace, name)
	return alamedaScaler, err
}

// GetAlamedaRecommendation return AlamedaRecommendation
func (getResource *GetResource) GetAlamedaRecommendation(namespace, name string) (*autuscaling.AlamedaRecommendation, error) {
	alamedaRecommendation := &autuscaling.AlamedaRecommendation{}
	err := getResource.getResource(alamedaRecommendation, namespace, name)
	return alamedaRecommendation, err
}

func (getResource *GetResource) GetObservingAlamedaScalerOfController(controllerType autuscaling.AlamedaControllerType, controllerNamespace, controllerName string) (*autuscaling.AlamedaScaler, error) {

	listResources := NewListResources(getResource)

	alamedaScalers, _ := listResources.ListNamespaceAlamedaScaler(controllerNamespace)
	for _, alamedaScaler := range alamedaScalers {

		switch controllerType {
		case autuscaling.DeploymentController:

			matchedLblDeployments, err := listResources.ListDeploymentsByNamespaceLabels(controllerNamespace, alamedaScaler.Spec.Selector.MatchLabels)
			if err != nil {
				return nil, errors.Errorf("get observing AlamedaScaler of Deployment %s/%s failed: %s", controllerNamespace, controllerName, err.Error())
			}
			for _, matchedLblDeployment := range matchedLblDeployments {
				// deployment can only join one AlamedaScaler
				if matchedLblDeployment.GetName() == controllerName {
					return &alamedaScaler, nil
				}
			}
		case autuscaling.DeploymentConfigController:

			matchedLblDeploymentConfigs, err := listResources.ListDeploymentConfigsByNamespaceLabels(controllerNamespace, alamedaScaler.Spec.Selector.MatchLabels)
			if err != nil {
				return nil, errors.Errorf("get observing AlamedaScaler of DeploymentConfig %s/%s failed: %s", controllerNamespace, controllerName, err.Error())
			}
			for _, matchedLblDeploymentConfig := range matchedLblDeploymentConfigs {
				// deploymentConfig can only join one AlamedaScaler
				if matchedLblDeploymentConfig.GetName() == controllerName {
					return &alamedaScaler, nil
				}
			}
		default:
			return nil, errors.Errorf("controllerType: %d not support", controllerType)
		}

	}

	return nil, nil
}

func (getResource *GetResource) getResource(resource runtime.Object, namespace, name string) error {
	if namespace == "" || name == "" {
		return errors.Errorf("Namespace: %s or name: %s is empty", namespace, name)
	}
	if err := getResource.Get(context.TODO(),
		types.NamespacedName{
			Namespace: namespace,
			Name:      name,
		},
		resource); err != nil {
		scope.Debug(err.Error())
		return err
	}
	return nil
}
