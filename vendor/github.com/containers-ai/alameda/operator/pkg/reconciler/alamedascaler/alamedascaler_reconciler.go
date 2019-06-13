package alamedascaler

import (
	"fmt"

	autoscaling_v1alpha1 "github.com/containers-ai/alameda/operator/pkg/apis/autoscaling/v1alpha1"
	utils "github.com/containers-ai/alameda/operator/pkg/utils"
	utilsresource "github.com/containers-ai/alameda/operator/pkg/utils/resources"
	logUtil "github.com/containers-ai/alameda/pkg/utils/log"
	appsapi_v1 "github.com/openshift/api/apps/v1"
	appsv1 "k8s.io/api/apps/v1"
	core_v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	alamedascalerReconcilerScope = logUtil.RegisterScope("alamedascaler_reconciler", "alamedascaler_reconciler", 0)
	podPhaseNeedsMonitoring      = map[core_v1.PodPhase]bool{
		core_v1.PodPending:   true,
		core_v1.PodRunning:   true,
		core_v1.PodSucceeded: false,
		core_v1.PodFailed:    false,
		core_v1.PodUnknown:   true,
	}
)

// Reconciler reconciles AlamedaScaler object
type Reconciler struct {
	client        client.Client
	alamedascaler *autoscaling_v1alpha1.AlamedaScaler
}

// NewReconciler creates Reconciler object
func NewReconciler(client client.Client, alamedascaler *autoscaling_v1alpha1.AlamedaScaler) *Reconciler {
	return &Reconciler{
		client:        client,
		alamedascaler: alamedascaler,
	}
}

// HasAlamedaDeployment checks the AlamedaScaler has the deployment or not
func (reconciler *Reconciler) HasAlamedaDeployment(deploymentNS, deploymentName string) bool {
	key := utils.GetNamespacedNameKey(deploymentNS, deploymentName)
	_, ok := reconciler.alamedascaler.Status.AlamedaController.Deployments[autoscaling_v1alpha1.NamespacedName(key)]
	return ok
}

// HasAlamedaDeploymentConfig checks the AlamedaScaler has the deploymentconfig or not
func (reconciler *Reconciler) HasAlamedaDeploymentConfig(deploymentConfigNS, deploymentConfigName string) bool {
	key := utils.GetNamespacedNameKey(deploymentConfigNS, deploymentConfigName)
	_, ok := reconciler.alamedascaler.Status.AlamedaController.DeploymentConfigs[autoscaling_v1alpha1.NamespacedName(key)]
	return ok
}

// HasAlamedaPod checks the AlamedaScaler has the AlamedaPod or not
func (reconciler *Reconciler) HasAlamedaPod(podNS, podName string) bool {
	for _, deployment := range reconciler.alamedascaler.Status.AlamedaController.Deployments {
		deploymentNS := deployment.Namespace
		for _, pod := range deployment.Pods {
			if deploymentNS == podNS && pod.Name == podName {
				return true
			}
		}
	}
	for _, deploymentConfig := range reconciler.alamedascaler.Status.AlamedaController.DeploymentConfigs {
		deploymentConfigNS := deploymentConfig.Namespace
		for _, pod := range deploymentConfig.Pods {
			if deploymentConfigNS == podNS && pod.Name == podName {
				return true
			}
		}
	}
	return false
}

// RemoveAlamedaDeployment removes deployment from alamedaController of AlamedaScaler
func (reconciler *Reconciler) RemoveAlamedaDeployment(deploymentNS, deploymentName string) *autoscaling_v1alpha1.AlamedaScaler {
	key := utils.GetNamespacedNameKey(deploymentNS, deploymentName)

	if _, ok := reconciler.alamedascaler.Status.AlamedaController.Deployments[autoscaling_v1alpha1.NamespacedName(key)]; ok {
		delete(reconciler.alamedascaler.Status.AlamedaController.Deployments, autoscaling_v1alpha1.NamespacedName(key))
	}

	return reconciler.alamedascaler
}

// RemoveAlamedaDeploymentConfig removes deployment from alamedaController of AlamedaScaler
func (reconciler *Reconciler) RemoveAlamedaDeploymentConfig(deploymentConfigNS, deploymentConfigName string) *autoscaling_v1alpha1.AlamedaScaler {
	key := utils.GetNamespacedNameKey(deploymentConfigNS, deploymentConfigName)

	if _, ok := reconciler.alamedascaler.Status.AlamedaController.DeploymentConfigs[autoscaling_v1alpha1.NamespacedName(key)]; ok {
		delete(reconciler.alamedascaler.Status.AlamedaController.DeploymentConfigs, autoscaling_v1alpha1.NamespacedName(key))
		return reconciler.alamedascaler
	}
	return reconciler.alamedascaler
}

// InitAlamedaController try to initialize alamedaController field of AlamedaScaler
func (reconciler *Reconciler) InitAlamedaController() (*autoscaling_v1alpha1.AlamedaScaler, bool) {
	needUpdate := false
	if reconciler.alamedascaler.Status.AlamedaController.Deployments == nil {
		reconciler.alamedascaler.Status.AlamedaController.Deployments = map[autoscaling_v1alpha1.NamespacedName]autoscaling_v1alpha1.AlamedaResource{}
		needUpdate = true
	}
	if reconciler.alamedascaler.Status.AlamedaController.DeploymentConfigs == nil {
		reconciler.alamedascaler.Status.AlamedaController.DeploymentConfigs = map[autoscaling_v1alpha1.NamespacedName]autoscaling_v1alpha1.AlamedaResource{}
		needUpdate = true
	}
	return reconciler.alamedascaler, needUpdate
}

// UpdateStatusByDeployment updates status by deployment
func (reconciler *Reconciler) UpdateStatusByDeployment(deployment *appsv1.Deployment) *autoscaling_v1alpha1.AlamedaScaler {
	alamedaScalerNS := reconciler.alamedascaler.GetNamespace()
	alamedaScalerName := reconciler.alamedascaler.GetName()

	listResources := utilsresource.NewListResources(reconciler.client)
	alamedaDeploymentNS := deployment.GetNamespace()
	alamedaDeploymentName := deployment.GetName()
	alamedaDeploymentUID := deployment.GetUID()
	alamedaPodsMap := map[autoscaling_v1alpha1.NamespacedName]autoscaling_v1alpha1.AlamedaPod{}
	alamedaDeploymentsMap := reconciler.alamedascaler.Status.AlamedaController.Deployments
	if alamedaDeploymentsMap == nil {
		alamedaDeploymentsMap = map[autoscaling_v1alpha1.NamespacedName]autoscaling_v1alpha1.AlamedaResource{}
	}
	if alamedaPods, err := listResources.ListPodsByDeployment(alamedaDeploymentNS, alamedaDeploymentName); err == nil && len(alamedaPods) > 0 {
		for _, alamedaPod := range alamedaPods {
			if !PodIsMonitoredByAlameda(&alamedaPod) {
				continue
			}
			alamedaPodNamespace := alamedaPod.GetNamespace()
			alamedaPodName := alamedaPod.GetName()
			alamedaPodUID := alamedaPod.GetUID()
			alamedascalerReconcilerScope.Debug(fmt.Sprintf("Pod (%s/%s) belongs to AlamedaScaler (%s/%s).", alamedaDeploymentNS, alamedaPodName, alamedaScalerNS, alamedaScalerName))
			alamedaContainers := []autoscaling_v1alpha1.AlamedaContainer{}
			for _, alamedaContainer := range alamedaPod.Spec.Containers {
				alamedaContainers = append(alamedaContainers, autoscaling_v1alpha1.AlamedaContainer{
					Name: alamedaContainer.Name,
				})
			}
			alamedaPodsMap[autoscaling_v1alpha1.NamespacedName(utils.GetNamespacedNameKey(alamedaPod.GetNamespace(), alamedaPodName))] = autoscaling_v1alpha1.AlamedaPod{
				Namespace:  alamedaPodNamespace,
				Name:       alamedaPodName,
				UID:        string(alamedaPodUID),
				Containers: alamedaContainers,
			}
		}
	}

	alamedaDeploymentsMap[autoscaling_v1alpha1.NamespacedName(utils.GetNamespacedNameKey(deployment.GetNamespace(), deployment.GetName()))] = autoscaling_v1alpha1.AlamedaResource{
		Namespace: alamedaDeploymentNS,
		Name:      alamedaDeploymentName,
		UID:       string(alamedaDeploymentUID),
		Pods:      alamedaPodsMap,
	}
	reconciler.alamedascaler.Status.AlamedaController.Deployments = alamedaDeploymentsMap
	return reconciler.alamedascaler
}

// UpdateStatusByDeploymentConfig updates status by DeploymentConfig
func (reconciler *Reconciler) UpdateStatusByDeploymentConfig(deploymentconfig *appsapi_v1.DeploymentConfig) *autoscaling_v1alpha1.AlamedaScaler {
	scalerNS := reconciler.alamedascaler.GetNamespace()
	scalerName := reconciler.alamedascaler.GetName()

	listResources := utilsresource.NewListResources(reconciler.client)
	deploymentConfigNS := deploymentconfig.GetNamespace()
	deploymentConfigName := deploymentconfig.GetName()
	deploymentConfigUID := deploymentconfig.GetUID()
	podsMap := map[autoscaling_v1alpha1.NamespacedName]autoscaling_v1alpha1.AlamedaPod{}
	deploymentConfigsMap := reconciler.alamedascaler.Status.AlamedaController.DeploymentConfigs
	if deploymentConfigsMap == nil {
		deploymentConfigsMap = map[autoscaling_v1alpha1.NamespacedName]autoscaling_v1alpha1.AlamedaResource{}
	}
	if alamedaPods, err := listResources.ListPodsByDeploymentConfig(deploymentConfigNS, deploymentConfigName); err == nil && len(alamedaPods) > 0 {
		for _, alamedaPod := range alamedaPods {
			if !PodIsMonitoredByAlameda(&alamedaPod) {
				continue
			}
			alamedaPodNamespace := alamedaPod.GetNamespace()
			alamedaPodName := alamedaPod.GetName()
			alamedaPodUID := alamedaPod.GetUID()
			alamedascalerReconcilerScope.Debug(fmt.Sprintf("Pod (%s/%s) belongs to AlamedaScaler (%s/%s).", deploymentConfigNS, alamedaPodName, scalerNS, scalerName))
			alamedaContainers := []autoscaling_v1alpha1.AlamedaContainer{}
			for _, alamedaContainer := range alamedaPod.Spec.Containers {
				alamedaContainers = append(alamedaContainers, autoscaling_v1alpha1.AlamedaContainer{
					Name: alamedaContainer.Name,
				})
			}
			podsMap[autoscaling_v1alpha1.NamespacedName(utils.GetNamespacedNameKey(alamedaPod.GetNamespace(), alamedaPodName))] = autoscaling_v1alpha1.AlamedaPod{
				Namespace:  alamedaPodNamespace,
				Name:       alamedaPodName,
				UID:        string(alamedaPodUID),
				Containers: alamedaContainers,
			}
		}
	}

	deploymentConfigsMap[autoscaling_v1alpha1.NamespacedName(utils.GetNamespacedNameKey(deploymentconfig.GetNamespace(), deploymentconfig.GetName()))] = autoscaling_v1alpha1.AlamedaResource{
		Namespace: deploymentConfigNS,
		Name:      deploymentConfigName,
		UID:       string(deploymentConfigUID),
		Pods:      podsMap,
	}
	reconciler.alamedascaler.Status.AlamedaController.DeploymentConfigs = deploymentConfigsMap
	return reconciler.alamedascaler
}

func PodIsMonitoredByAlameda(pod *core_v1.Pod) bool {
	if !podPhaseIsMonitoredByAlameda(pod.Status.Phase) || pod.ObjectMeta.DeletionTimestamp != nil {
		return false
	}
	return true
}

func podPhaseIsMonitoredByAlameda(podPhase core_v1.PodPhase) bool {
	if isMonitored, exist := podPhaseNeedsMonitoring[podPhase]; exist {
		return isMonitored
	}
	return false
}
