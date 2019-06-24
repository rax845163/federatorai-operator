package metadata

import (
	"fmt"
	"sync"

	openshiftappsv1 "github.com/openshift/api/apps/v1"
	"github.com/openshift/client-go/apps/clientset/versioned"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	sigs_k8s_client_config "sigs.k8s.io/controller-runtime/pkg/client/config"
)

var (
	alamedaRecommendationGVR = schema.GroupVersionResource{}
	resourcesKindMapMutex    = &sync.Mutex{}
)

// OwnerReferenceTracer struct to trace owner references
type OwnerReferenceTracer struct {
	k8sClient          kubernetes.Interface
	k8sDynamicClient   dynamic.Interface
	k8sDiscoveryClient *discovery.DiscoveryClient
	k8sRestMapper      meta.RESTMapper

	openshiftClientset versioned.Interface
}

// NewDefaultOwnerReferenceTracer build OwnerReferenceTracer
func NewDefaultOwnerReferenceTracer() (*OwnerReferenceTracer, error) {

	config, err := sigs_k8s_client_config.GetConfig()
	if err != nil {
		return nil, errors.Errorf("new OwnerReferenceTracer failed: %s", err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Errorf("new OwnerReferenceTracer failed: %s", err.Error())
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, errors.Errorf("new OwnerReferenceTracer failed: %s", err.Error())
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, errors.Errorf("new OwnerReferenceTracer failed: %s", err.Error())
	}

	gr, err := restmapper.GetAPIGroupResources(discoveryClient)
	if err != nil {
		return nil, errors.Errorf("new OwnerReferenceTracer failed: %s", err.Error())
	}
	restMapper := restmapper.NewDiscoveryRESTMapper(gr)

	openshiftClientset, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, errors.Errorf("new OwnerReferenceTracer failed: %s", err.Error())
	}

	o := &OwnerReferenceTracer{
		k8sClient:          client,
		k8sDynamicClient:   dynamicClient,
		k8sDiscoveryClient: discoveryClient,
		k8sRestMapper:      restMapper,
		openshiftClientset: openshiftClientset,
	}

	return o, nil
}

// NewOwnerReferenceTracerWithConfig build OwnerReferenceTracer
func NewOwnerReferenceTracerWithConfig(cfg rest.Config) (*OwnerReferenceTracer, error) {

	copyCfg := cfg

	client, err := kubernetes.NewForConfig(&copyCfg)
	if err != nil {
		return nil, errors.Errorf("new resource recommendator failed: %s", err.Error())
	}

	dynamicClient, err := dynamic.NewForConfig(&copyCfg)
	if err != nil {
		return nil, errors.Errorf("new resource recommendator failed: %s", err.Error())
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(&copyCfg)
	if err != nil {
		return nil, errors.Errorf("new resource recommendator failed: %s", err.Error())
	}

	gr, err := restmapper.GetAPIGroupResources(discoveryClient)
	if err != nil {
		return nil, errors.Errorf("new OwnerReferenceTracer failed: %s", err.Error())
	}
	restMapper := restmapper.NewDiscoveryRESTMapper(gr)

	openshiftClientset, err := versioned.NewForConfig(&copyCfg)
	if err != nil {
		return nil, errors.Errorf("new OwnerReferenceTracer failed: %s", err.Error())
	}

	impl := &OwnerReferenceTracer{
		k8sClient:          client,
		k8sDynamicClient:   dynamicClient,
		k8sDiscoveryClient: discoveryClient,
		k8sRestMapper:      restMapper,
		openshiftClientset: openshiftClientset,
	}

	return impl, nil
}

// GetRootControllerKindAndNameOfOwnerReferences gets root owner references that is Controller
func (ort *OwnerReferenceTracer) GetRootControllerKindAndNameOfOwnerReferences(namespace string, ownerRefs []meta_v1.OwnerReference) (kind, name string, err error) {

	var nextControllerOwnerRef *meta_v1.OwnerReference
	var foundControllerOwnerRef *meta_v1.OwnerReference

	for true {

		nextControllerOwnerRef = nil

		// get owner that is controller
		for _, ownerRef := range ownerRefs {
			if ownerRef.Controller != nil && *ownerRef.Controller {
				nextControllerOwnerRef = &ownerRef
				break
			}
		}

		// no next controller ownerReference find
		if nextControllerOwnerRef == nil {
			break
		} else {
			foundControllerOwnerRef = nextControllerOwnerRef
		}

		gvk := schema.FromAPIVersionAndKind(nextControllerOwnerRef.APIVersion, nextControllerOwnerRef.Kind)
		ownerRefs, err = ort.getOwnerRefsOfResource(namespace, nextControllerOwnerRef.Name, gvk)
		if err != nil {
			return "", "", errors.Wrap(err, "get root controller name from owner references failed")
		}

	}

	if foundControllerOwnerRef == nil {
		err = errors.Errorf("none of ownerReferences is controller")
	} else {
		kind = foundControllerOwnerRef.Kind
		name = foundControllerOwnerRef.Name
	}

	return kind, name, err
}

func (ort *OwnerReferenceTracer) GetDeploymentOrDeploymentConfigOwningPod(pod core_v1.Pod) (*appsv1.Deployment, *openshiftappsv1.DeploymentConfig, error) {

	searchingNamespace := pod.Namespace
	ownerRefs := pod.GetOwnerReferences()

	var controllerOwnerRef *meta_v1.OwnerReference
	finish := false
	for !finish {

		if len(ownerRefs) == 0 {
			finish = true
			break
		}

		// get owner that is controller
		for _, ownerRef := range ownerRefs {
			if ownerRef.Controller != nil && *ownerRef.Controller {
				ownerName := ownerRef.Name
				switch ownerRef.Kind {
				case "Deployment":
					dep, err := ort.k8sClient.AppsV1().Deployments(searchingNamespace).Get(ownerName, meta_v1.GetOptions{})
					if err != nil {
						return nil, nil, errors.Errorf("get deployment owning pod %s/%s failed, %s", err.Error())
					}
					return dep, nil, nil
				case "DeploymentConfig":
					depConfig, err := ort.openshiftClientset.AppsV1().DeploymentConfigs(searchingNamespace).Get(ownerName, meta_v1.GetOptions{})
					if err != nil {
						return nil, nil, errors.Errorf("get deployment owning pod %s/%s failed, %s", err.Error())
					}
					return nil, depConfig, nil
				}
				controllerOwnerRef = &ownerRef
				break
			}
		}

		// there is no ownerReference that is Controller, need no tracing
		if controllerOwnerRef == nil {
			finish = true
			break
		}

		gvk := schema.FromAPIVersionAndKind(controllerOwnerRef.APIVersion, controllerOwnerRef.Kind)
		resOwnerRefs, err := ort.getOwnerRefsOfResource(searchingNamespace, controllerOwnerRef.Name, gvk)
		if err != nil {
			return nil, nil, errors.Wrap(err, "get deployment or deploymentConfig owning pod failed")
		}
		ownerRefs = resOwnerRefs
	}

	return nil, nil, nil
}

func (ort *OwnerReferenceTracer) GetControllerOwnerReferenceLink(objectMeta meta_v1.Object) ([]meta_v1.OwnerReference, error) {

	var link = make([]meta_v1.OwnerReference, 0)
	var nextControllerOwnerRef *meta_v1.OwnerReference
	var namespace = objectMeta.GetNamespace()

	ownerRefs := objectMeta.GetOwnerReferences()
	for true {

		nextControllerOwnerRef = nil

		// get owner that is controller
		for _, ownerRef := range ownerRefs {
			if ownerRef.Controller != nil && *ownerRef.Controller {
				link = append(link, ownerRef)
				nextControllerOwnerRef = &ownerRef
				break
			}
		}

		// no next controller ownerReference find
		if nextControllerOwnerRef == nil {
			break
		}

		var err error
		gvk := schema.FromAPIVersionAndKind(nextControllerOwnerRef.APIVersion, nextControllerOwnerRef.Kind)
		ownerRefs, err = ort.getOwnerRefsOfResource(namespace, nextControllerOwnerRef.Name, gvk)
		if err != nil {
			return link, err
		}
	}

	return link, nil
}

func (ort *OwnerReferenceTracer) getOwnerRefsOfResource(namespace, name string, gvk schema.GroupVersionKind) ([]meta_v1.OwnerReference, error) {

	ownerRefs := make([]meta_v1.OwnerReference, 0)

	restMapping, err := ort.k8sRestMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return ownerRefs, errors.Errorf("get owner references of %s/%s gvk: %s failed: %s", namespace, name, gvk.String(), err.Error())
	}

	gvr := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: fmt.Sprintf("namespaces/%s/%s", namespace, restMapping.Resource.Resource),
	}
	us, err := ort.k8sDynamicClient.Resource(gvr).Get(name, meta_v1.GetOptions{})
	if err != nil {
		return ownerRefs, errors.Errorf("get owner references of resource %s in namespace %s failed: %s", gvr.String(), namespace, err.Error())
	}
	ownerRefs = us.GetOwnerReferences()

	return ownerRefs, nil
}
