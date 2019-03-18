package component

import (
	"github.com/containers-ai/federatorai-operator/pkg/assets"
	"github.com/containers-ai/federatorai-operator/pkg/lib/resourceread"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
)

func NewComponentADeployment() *appsv1.Deployment {
	deploymentBytes, err := assets.Asset("../../manifests/deployment.yaml")
	if err != nil {
		logrus.Fatalf("Failed to Test create deployment: %v", err)
	}
	d := resourceread.ReadDeploymentV1OrDie(deploymentBytes)
	return d
}
func NewComponentAServiceAccount() *corev1.ServiceAccount {
	saByte, err := assets.Asset("../../manifests/serviceaccount.yaml")
	if err != nil {
		logrus.Fatalf("Failed to Test create deployment: %v", err)
	}
	sa := resourceread.ReadServiceAccountV1OrDie(saByte)
	return sa
}
func NewComponentAConfigMap() *corev1.ConfigMap {
	cmByte, err := assets.Asset("../../manifests/configmap.yaml")
	if err != nil {
		logrus.Fatalf("Failed to Test create ConfigMap: %v", err)
	}
	cm := resourceread.ReadConfigMapV1OrDie(cmByte)
	return cm
}
func NewComponentAService() *corev1.Service {
	svByte, err := assets.Asset("../../manifests/service.yaml")
	if err != nil {
		logrus.Fatalf("Failed to Test create Service: %v", err)
	}
	sv := resourceread.ReadServiceV1OrDie(svByte)
	return sv
}
func NewComponentAClusterRole() *rbacv1.ClusterRole {
	crByte, err := assets.Asset("../../manifests/clusterrole.yaml")
	if err != nil {
		logrus.Fatalf("Failed to Test create clusterrole: %v", err)
	}
	cr := resourceread.ReadClusterRoleV1OrDie(crByte)
	return cr
}
func NewComponentAClusterRoleBinding() *rbacv1.ClusterRoleBinding {
	crbByte, err := assets.Asset("../../manifests/clusterrolebinding.yaml")
	if err != nil {
		logrus.Fatalf("Failed to Test create clusterrolebinding: %v", err)
	}
	crb := resourceread.ReadClusterRoleBindingV1OrDie(crbByte)
	return crb
}
func NewComponentAPersistentVolumeClaim() *corev1.PersistentVolumeClaim {
	pvcByte, err := assets.Asset("../../manifests/persistentvolumeclaim.yaml")
	if err != nil {
		logrus.Fatalf("Failed to Test create persistentvolumeclaim: %v", err)
	}
	pvc := resourceread.ReadPersistentVolumeClaimV1OrDie(pvcByte)
	return pvc
}

func NewClusterRoleBinding(str string) *rbacv1.ClusterRoleBinding {
	crbByte, err := assets.Asset(str)
	if err != nil {
		logrus.Fatalf("Failed to Test create clusterrolebinding: %v", err)
	}
	crb := resourceread.ReadClusterRoleBindingV1OrDie(crbByte)
	return crb
}
func NewClusterRole(str string) *rbacv1.ClusterRole {
	crByte, err := assets.Asset(str)
	if err != nil {
		logrus.Fatalf("Failed to Test create clusterrole: %v", err)
	}
	cr := resourceread.ReadClusterRoleV1OrDie(crByte)
	return cr
}
func NewServiceAccount(str string) *corev1.ServiceAccount {
	saByte, err := assets.Asset(str)
	if err != nil {
		logrus.Fatalf("Failed to Test create deployment: %v", err)
	}
	sa := resourceread.ReadServiceAccountV1OrDie(saByte)
	return sa
}
func NewConfigMap(str string) *corev1.ConfigMap {
	cmByte, err := assets.Asset(str)
	if err != nil {
		logrus.Fatalf("Failed to Test create ConfigMap: %v", err)
	}
	cm := resourceread.ReadConfigMapV1OrDie(cmByte)
	return cm
}
func NewPersistentVolumeClaim(str string) *corev1.PersistentVolumeClaim {
	pvcByte, err := assets.Asset(str)
	if err != nil {
		logrus.Fatalf("Failed to Test create persistentvolumeclaim: %v", err)
	}
	pvc := resourceread.ReadPersistentVolumeClaimV1OrDie(pvcByte)
	return pvc
}
func NewService(str string) *corev1.Service {
	svByte, err := assets.Asset(str)
	if err != nil {
		logrus.Fatalf("Failed to Test create deployment: %v", err)
	}
	sv := resourceread.ReadServiceV1OrDie(svByte)
	return sv
}
func NewDeployment(str string) *appsv1.Deployment {
	deploymentBytes, err := assets.Asset(str)
	if err != nil {
		logrus.Fatalf("Failed to Test create deployment: %v", err)
	}
	d := resourceread.ReadDeploymentV1OrDie(deploymentBytes)
	return d
}

//NewClusterRoleBinding
//ReadConfigMapV1OrDie
