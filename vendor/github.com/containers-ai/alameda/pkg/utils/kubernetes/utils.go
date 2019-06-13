package kubernetes

import (
	"bufio"
	"os"
	"strings"

	"github.com/containers-ai/alameda/pkg/consts"
	logUtil "github.com/containers-ai/alameda/pkg/utils/log"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var scope = logUtil.RegisterScope("kubernetes_utils", "Kubernetes utils.", 0)

func IsOKDCluster() (bool, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return false, err
	}
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return false, err
	}
	apiResourceLists, err := discoveryClient.ServerResources()
	if err != nil {
		return false, err
	}

	for _, apiResourceList := range apiResourceLists {
		for _, resource := range apiResourceList.APIResources {
			if resource.Kind == consts.K8S_KIND_DEPLOYMENTCONFIG {
				return true, nil
			}
		}
	}
	return false, nil
}

func GetRunningNamespace() string {
	ns := ""
	nsFile, err := os.Open("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		scope.Errorf(err.Error())
	}
	defer nsFile.Close()

	scanner := bufio.NewScanner(nsFile)
	for scanner.Scan() {
		ns = ns + scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		scope.Errorf(err.Error())
	}
	return strings.TrimSpace(ns)
}

// NewOwnerReference provides ownerReference reference to input object
func NewOwnerReference(objType metav1.TypeMeta, objMeta metav1.ObjectMeta, isController bool) metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion: objType.APIVersion,
		Kind:       objType.Kind,
		Name:       objMeta.Name,
		UID:        objMeta.UID,
		Controller: &isController,
	}
}

// AddOwnerRefToObject appends the desired OwnerReference to the object
func AddOwnerRefToObject(obj metav1.Object, ownerRef metav1.OwnerReference) {

	for _, objOwnerReference := range obj.GetOwnerReferences() {
		if objOwnerReference.UID == ownerRef.UID {
			return
		}
	}

	obj.SetOwnerReferences(append(obj.GetOwnerReferences(), ownerRef))
}

// GetPodByNamespaceNameWithConfig fetchs pod resource by namespace and name with k8s rest configuration
func GetPodByNamespaceNameWithConfig(namespace, name string, config rest.Config) (corev1.Pod, error) {

	pod := corev1.Pod{}

	if namespace == "" || name == "" {
		return pod, errors.New("get pod by namespace and name failed, cannot get pod if namespace and name do not provide")
	}

	k8sClient, err := kubernetes.NewForConfig(&config)
	if err != nil {
		return pod, errors.Errorf("get pod by namespace and name failed, create k8s api client failed: %s", err.Error())
	}

	p, err := k8sClient.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return pod, errors.Errorf("get pod by namespace and name failed, %s", err.Error())
	}
	pod = *p

	return pod, nil
}
