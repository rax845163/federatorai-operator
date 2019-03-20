package component

import (
	"github.com/containers-ai/federatorai-operator/pkg/assets"
	"github.com/containers-ai/federatorai-operator/pkg/lib/resourceread"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("controller_alamedaservice")

func NewComponentADeployment() *appsv1.Deployment {
	deploymentBytes, err := assets.Asset("../../manifests/deployment.yaml")
	if err != nil {
		log.Error(err, "Failed to Test create deployment")

	}
	d := resourceread.ReadDeploymentV1OrDie(deploymentBytes)
	return d
}
func NewComponentAServiceAccount() *corev1.ServiceAccount {
	saByte, err := assets.Asset("../../manifests/serviceaccount.yaml")
	if err != nil {
		log.Error(err, "Failed to Test create ServiceAccount")
	}
	sa := resourceread.ReadServiceAccountV1OrDie(saByte)
	return sa
}
func NewComponentAConfigMap() *corev1.ConfigMap {
	cmByte, err := assets.Asset("../../manifests/configmap.yaml")
	if err != nil {
		log.Error(err, "Failed to Test create ConfigMap")
	}
	cm := resourceread.ReadConfigMapV1OrDie(cmByte)
	return cm
}
func NewComponentAService() *corev1.Service {
	svByte, err := assets.Asset("../../manifests/service.yaml")
	if err != nil {
		log.Error(err, "Failed to Test create Service")

	}
	sv := resourceread.ReadServiceV1OrDie(svByte)
	return sv
}
func NewComponentAClusterRole() *rbacv1.ClusterRole {
	crByte, err := assets.Asset("../../manifests/clusterrole.yaml")
	if err != nil {
		log.Error(err, "Failed to Test create clusterrole")

	}
	cr := resourceread.ReadClusterRoleV1OrDie(crByte)
	return cr
}
func NewComponentAClusterRoleBinding() *rbacv1.ClusterRoleBinding {
	crbByte, err := assets.Asset("../../manifests/clusterrolebinding.yaml")
	if err != nil {
		log.Error(err, "Failed to Test create clusterrolebinding")

	}
	crb := resourceread.ReadClusterRoleBindingV1OrDie(crbByte)
	return crb
}
func NewComponentAPersistentVolumeClaim() *corev1.PersistentVolumeClaim {
	pvcByte, err := assets.Asset("../../manifests/persistentvolumeclaim.yaml")
	if err != nil {
		log.Error(err, "Failed to Test create persistentvolumeclaim")

	}
	pvc := resourceread.ReadPersistentVolumeClaimV1OrDie(pvcByte)
	return pvc
}

func NewClusterRoleBinding(str string) *rbacv1.ClusterRoleBinding {
	crbByte, err := assets.Asset(str)
	if err != nil {
		log.Error(err, "Failed to Test create clusterrolebinding")

	}
	crb := resourceread.ReadClusterRoleBindingV1OrDie(crbByte)
	return crb
}
func NewClusterRole(str string) *rbacv1.ClusterRole {
	crByte, err := assets.Asset(str)
	if err != nil {
		log.Error(err, "Failed to Test create clusterrole")

	}
	cr := resourceread.ReadClusterRoleV1OrDie(crByte)
	return cr
}
func NewServiceAccount(str string) *corev1.ServiceAccount {
	saByte, err := assets.Asset(str)
	if err != nil {
		log.Error(err, "Failed to Test create serviceaccount")

	}
	sa := resourceread.ReadServiceAccountV1OrDie(saByte)
	return sa
}
func NewConfigMap(str string) *corev1.ConfigMap {
	cmByte, err := assets.Asset(str)
	if err != nil {
		log.Error(err, "Failed to Test create configmap")

	}
	cm := resourceread.ReadConfigMapV1OrDie(cmByte)
	return cm
}
func NewPersistentVolumeClaim(str string) *corev1.PersistentVolumeClaim {
	pvcByte, err := assets.Asset(str)
	if err != nil {
		log.Error(err, "Failed to Test create persistentvolumeclaim")

	}
	pvc := resourceread.ReadPersistentVolumeClaimV1OrDie(pvcByte)
	return pvc
}
func NewService(str string) *corev1.Service {
	svByte, err := assets.Asset(str)
	if err != nil {
		log.Error(err, "Failed to Test create service")

	}
	sv := resourceread.ReadServiceV1OrDie(svByte)
	return sv
}
func NewDeployment(str string) *appsv1.Deployment {
	deploymentBytes, err := assets.Asset(str)
	if err != nil {
		log.Error(err, "Failed to Test create deployment")

	}
	d := resourceread.ReadDeploymentV1OrDie(deploymentBytes)
	return d
}

func RegistryCustomResourceDefinition(str string) *apiextv1beta1.CustomResourceDefinition {
	crdBytes, err := assets.Asset(str)
	if err != nil {
		log.Error(err, "Failed to Test create testcrd")
	}
	crd := resourceread.ReadCustomResourceDefinitionV1Beta1OrDie(crdBytes)
	return crd
}

