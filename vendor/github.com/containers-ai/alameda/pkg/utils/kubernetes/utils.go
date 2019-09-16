package kubernetes

import (
	"bufio"
	"context"
	"fmt"
	Consts "github.com/containers-ai/alameda/pkg/consts"
	Log "github.com/containers-ai/alameda/pkg/utils/log"
	"github.com/pkg/errors"
	Corev1 "k8s.io/api/core/v1"
	Metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	K8SErrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"strings"
)

var scope = Log.RegisterScope("kubernetes_utils", "Kubernetes utils.", 0)

func NewK8SClient() (client.Client, error){
	k8sClientConfig, err := config.GetConfig()
	if err != nil {
		return nil, errors.New("Failed to get kubernetes configuration: " + err.Error())
	}

	k8sClient, err := client.New(k8sClientConfig, client.Options{})
	if err != nil {
		return nil, errors.New("Failed to create kubernetes client: " + err.Error())
	}

	return k8sClient, nil
}

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
			if resource.Kind == Consts.K8S_KIND_DEPLOYMENTCONFIG {
				return true, nil
			}
		}
	}

	return false, nil
}

func GetClusterUID(k8sClient client.Client) (string, error) {
	possibleNSList := []string{"default"}
	errorList := make([]string, 0)
	clusterId := ""

	for _, possibleNS := range possibleNSList {
		clusterInfoCM := &Corev1.ConfigMap{}
		err := k8sClient.Get(context.Background(), client.ObjectKey{
			Name:      "cluster-info",
			Namespace: possibleNS,
		}, clusterInfoCM)
		if err == nil {
			return string(clusterInfoCM.GetUID()), nil
		} else if !K8SErrors.IsNotFound(err) {
			errorList = append(errorList, err.Error())
		}
	}

	if len(errorList) == 0 {
		return clusterId, fmt.Errorf("no cluster info found")
	}

	return clusterId, errors.New(strings.Join(errorList, ","))
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

func GetPodName() string {
	ns := ""
	nsFile, err := os.Open("/etc/hostname")
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
func NewOwnerReference(objType Metav1.TypeMeta, objMeta Metav1.ObjectMeta, isController bool) Metav1.OwnerReference {
	return Metav1.OwnerReference{
		APIVersion: objType.APIVersion,
		Kind:       objType.Kind,
		Name:       objMeta.Name,
		UID:        objMeta.UID,
		Controller: &isController,
	}
}

// AddOwnerRefToObject appends the desired OwnerReference to the object
func AddOwnerRefToObject(obj Metav1.Object, ownerRef Metav1.OwnerReference) {
	for _, objOwnerReference := range obj.GetOwnerReferences() {
		if objOwnerReference.UID == ownerRef.UID {
			return
		}
	}

	obj.SetOwnerReferences(append(obj.GetOwnerReferences(), ownerRef))
}

// GetPodByNamespaceNameWithConfig fetchs pod resource by namespace and name with k8s rest configuration
func GetPodByNamespaceNameWithConfig(namespace, name string, config rest.Config) (Corev1.Pod, error) {
	pod := Corev1.Pod{}

	if namespace == "" || name == "" {
		return pod, errors.New("get pod by namespace and name failed, cannot get pod if namespace and name do not provide")
	}

	k8sClient, err := kubernetes.NewForConfig(&config)
	if err != nil {
		return pod, errors.Errorf("get pod by namespace and name failed, create k8s api client failed: %s", err.Error())
	}

	p, err := k8sClient.CoreV1().Pods(namespace).Get(name, Metav1.GetOptions{})
	if err != nil {
		return pod, errors.Errorf("get pod by namespace and name failed, %s", err.Error())
	}
	pod = *p

	return pod, nil
}
