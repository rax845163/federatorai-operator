package resourceread

import (
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	rbacScheme = runtime.NewScheme()
	rbacCodecs = serializer.NewCodecFactory(rbacScheme)
)

func init() {
	if err := rbacv1.AddToScheme(rbacScheme); err != nil {
		log.Error(err, "Fail AddToScheme")
	}
}

func ReadClusterRoleBindingV1(objBytes []byte) *rbacv1.ClusterRoleBinding {
	requiredObj, err := runtime.Decode(rbacCodecs.UniversalDecoder(rbacv1.SchemeGroupVersion), objBytes)
	if err != nil {
		log.Error(err, "Fail ReadClusterRoleBindingV1")
	}
	return requiredObj.(*rbacv1.ClusterRoleBinding)
}

func ReadRoleBindingV1(objBytes []byte) *rbacv1.RoleBinding {
	requiredObj, err := runtime.Decode(rbacCodecs.UniversalDecoder(rbacv1.SchemeGroupVersion), objBytes)
	if err != nil {
		log.Error(err, "Fail ReadRoleBindingV1")
	}
	return requiredObj.(*rbacv1.RoleBinding)
}

func ReadClusterRoleV1(objBytes []byte) *rbacv1.ClusterRole {
	requiredObj, err := runtime.Decode(rbacCodecs.UniversalDecoder(rbacv1.SchemeGroupVersion), objBytes)
	if err != nil {
		log.Error(err, "Fail ReadClusterRoleV1")
	}
	return requiredObj.(*rbacv1.ClusterRole)
}
