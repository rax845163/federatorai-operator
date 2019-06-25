package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
	"time"

	admission_controller_kubernetes "github.com/containers-ai/alameda/admission-controller/pkg/kubernetes"
	"github.com/containers-ai/alameda/admission-controller/pkg/recommendator/resource"
	admission_controller_utils "github.com/containers-ai/alameda/admission-controller/pkg/utils"
	controller_validator "github.com/containers-ai/alameda/admission-controller/pkg/validator/controller"
	autoscalingv1alpha1 "github.com/containers-ai/alameda/operator/pkg/apis/autoscaling/v1alpha1"
	alamedascaler_reconciler "github.com/containers-ai/alameda/operator/pkg/reconciler/alamedascaler"
	"github.com/containers-ai/alameda/operator/pkg/utils/resources"
	metadata_utils "github.com/containers-ai/alameda/pkg/utils/kubernetes/metadata"
	"github.com/containers-ai/alameda/pkg/utils/log"

	"github.com/mattbaird/jsonpatch"
	"github.com/pkg/errors"
	admission_v1beta1 "k8s.io/api/admission/v1beta1"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	DefaultPodMutatePatchValdationFunction = admission_controller_utils.ValidatePatchFunc(func(patch jsonpatch.JsonPatchOperation) error {
		return nil
	})
	OKD3_9tPodMutatePatchValdationFunction = admission_controller_utils.ValidatePatchFunc(func(patch jsonpatch.JsonPatchOperation) error {

		allowedPatchOperations := map[string]bool{
			"add": true,
		}
		if allowed, exist := allowedPatchOperations[patch.Operation]; !exist || !allowed {
			return errors.Errorf("cannot patch with operation %s", patch.Operation)
		}

		return nil
	})

	patchType                = admission_v1beta1.PatchTypeJSONPatch
	scope                    = log.RegisterScope("admission-controller", "admission-controller", 0)
	defaultAdmissionResponse = admission_v1beta1.AdmissionResponse{
		Allowed: true,
	}
)

type admitFunc func(*admission_v1beta1.AdmissionReview) (admission_v1beta1.AdmissionResponse, error)

type admissionController struct {
	config *Config

	lock                                  *sync.Mutex
	controllerRecommendationMap           map[namespaceKindName]*controllerRecommendation
	controllerLockMap                     map[namespaceKindName]*sync.Mutex
	resourceRecommendatorSyncTimeout      time.Duration
	resourceRecommendatorSyncRetryTime    int
	resourceRecommendatorSyncWaitInterval time.Duration

	sigsK8SClient         client.Client
	k8sDeserializer       runtime.Decoder
	ownerReferenceTracer  *metadata_utils.OwnerReferenceTracer
	resourceRecommendator resource.ResourceRecommendator
	controllerValidator   controller_validator.Validator

	podMutatePatchValdationFunction admission_controller_utils.ValidatePatchFunc
}

// NewAdmissionControllerWithConfig creates AdmissionController with configuration and dependencies
func NewAdmissionControllerWithConfig(cfg Config, sigsK8SClient client.Client, resourceRecommendator resource.ResourceRecommendator, controllerValidator controller_validator.Validator, podMutatePatchValdationFunction admission_controller_utils.ValidatePatchFunc) (AdmissionController, error) {

	defaultOwnerReferenceTracer, err := metadata_utils.NewDefaultOwnerReferenceTracer()
	if err != nil {
		return nil, errors.Wrap(err, "new AdmissionController failed")
	}

	ac := &admissionController{
		config: &cfg,

		lock:                                  &sync.Mutex{},
		controllerRecommendationMap:           make(map[namespaceKindName]*controllerRecommendation),
		controllerLockMap:                     make(map[namespaceKindName]*sync.Mutex),
		resourceRecommendatorSyncTimeout:      10 * time.Second,
		resourceRecommendatorSyncRetryTime:    3,
		resourceRecommendatorSyncWaitInterval: 5 * time.Second,

		sigsK8SClient:         sigsK8SClient,
		k8sDeserializer:       admission_controller_kubernetes.Codecs.UniversalDecoder(),
		ownerReferenceTracer:  defaultOwnerReferenceTracer,
		resourceRecommendator: resourceRecommendator,
		controllerValidator:   controllerValidator,

		podMutatePatchValdationFunction: podMutatePatchValdationFunction,
	}

	return ac, nil
}

func (ac *admissionController) MutatePod(w http.ResponseWriter, r *http.Request) {
	ac.serve(w, r, ac.mutatePod)
}

func (ac *admissionController) serve(w http.ResponseWriter, r *http.Request, admit admitFunc) {

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		scope.Warnf("serve failed, skip serving: receive contentType=%s, expect application/json", contentType)
		ac.writeDefaultAdmissionReview(w)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		scope.Warnf("serve failed, skip serving: read http request failed: %s", err.Error())
		ac.writeDefaultAdmissionReview(w)
		return
	}

	admissionReview := &admission_v1beta1.AdmissionReview{}
	if err := json.Unmarshal(body, admissionReview); err != nil {
		scope.Warnf("serve failed, skip serving: unmarshal AdmissionReview failed: %s", err.Error())
		ac.writeDefaultAdmissionReview(w)
		return
	}
	admissionResponse, err := admit(admissionReview)
	if err != nil {
		scope.Warnf("admit with error: %s, skip serving AdmissionReview: %+v", err.Error(), admissionReview)
		ac.writeDefaultAdmissionReview(w)
		return
	}
	admissionResponse.UID = admissionReview.Request.UID

	err = ac.writeAdmissionReview(w, admissionResponse)
	if err != nil {
		scope.Warnf("")
	}
}

func (ac *admissionController) writeAdmissionReview(w http.ResponseWriter, admissionResponse admission_v1beta1.AdmissionResponse) error {

	admissionReview := admission_v1beta1.AdmissionReview{
		Response: admissionResponse.DeepCopy(),
	}
	admissionReviewBytes, err := json.Marshal(admissionReview)
	if err != nil {
		return errors.Errorf("marshal AdmissionReview failed: %s", err.Error())
	}

	_, err = w.Write(admissionReviewBytes)
	if err != nil {
		return errors.Errorf("write AdmissionReview failed: %s", err.Error())
	}

	return nil
}

func (ac *admissionController) writeDefaultAdmissionReview(w http.ResponseWriter) {

	err := ac.writeAdmissionReview(w, defaultAdmissionResponse)
	if err != nil {
		scope.Warnf("write default AdmissionReview failed: %s", err.Error())
	}
}

func (ac *admissionController) mutatePod(ar *admission_v1beta1.AdmissionReview) (admission_v1beta1.AdmissionResponse, error) {

	admissionResponse := admission_v1beta1.AdmissionResponse{Allowed: true}

	podResource := meta_v1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	if ar.Request.Resource != podResource {
		err := errors.Errorf("mutating pod failed: expect resource to be %s, get %s, skip mutating pod", podResource.String(), ar.Request.Resource.String())
		return admissionResponse, err
	}

	raw := ar.Request.Object.Raw
	pod := core_v1.Pod{}
	if _, _, err := ac.k8sDeserializer.Decode(raw, nil, &pod); err != nil {
		return admissionResponse, errors.Errorf("mutating pod failed: deserialize AdmissionRequest.Raw to Pod failed, skip mutating pod: %s", err.Error())
	}
	pod.SetNamespace(ar.Request.Namespace)

	scope.Infof("mutating pod: %+v", pod.ObjectMeta)

	controllerID, err := ac.getControllerIDToQueryPodRecommendations(&pod)
	if err != nil {
		return admissionResponse, errors.Wrapf(err, "mutating pod failed: get controller information of Pod failed, skip mutating pod: %s", err.Error())
	}

	executionEnabeld, err := ac.isControllerExecutionEnabled(controllerID)
	if err != nil {
		return admissionResponse, errors.Wrapf(err, "check if pod needs mutating faield, skip mutating pod: Pod: %+v", pod.ObjectMeta)
	} else if !executionEnabeld {
		return admissionResponse, errors.Errorf("execution of AlamedaScaler monitoring this pod is not enabled, skip mutating pod: Pod: %+v", pod.ObjectMeta)
	}
	recommendation, err := ac.getPodResourceRecommendationByControllerID(controllerID)
	if err != nil {
		return admissionResponse, errors.Errorf("get pod resource recommendations failed, controllerID: %s, skip mutating pod: Pod: %+v, errMsg: %s", controllerID.String(), pod.ObjectMeta, err.Error())
	} else if recommendation == nil {
		return admissionResponse, errors.Errorf("fetch empty recommendations of controller, controllerID: %s, skip mutating pod: Pod: %+v", controllerID.String(), pod.ObjectMeta)
	}

	patches, err := admission_controller_utils.GetPatchesFromPodResourceRecommendation(&pod, recommendation)
	if err != nil {
		return admissionResponse, errors.Wrapf(err, "get patches to mutate pod resource failed, skip mutating pod: Pod: %+v", pod.ObjectMeta)
	}
	err = admission_controller_utils.ValidatePatches(patches, ac.podMutatePatchValdationFunction)
	if err != nil {
		return admissionResponse, errors.Wrapf(err, "validate patches to mutate pod resource failed, skip mutating pod: Pod: %+v", pod.ObjectMeta)
	}
	patchString := admission_controller_utils.GetK8SPatchesString(patches)
	scope.Infof("patch %s to pod %+v ", patchString, pod.ObjectMeta)

	admissionResponse.Patch = []byte(patchString)
	admissionResponse.PatchType = &patchType

	return admissionResponse, nil
}

func (ac *admissionController) getPodResourceRecommendationByControllerID(controllerID namespaceKindName) (*resource.PodResourceRecommendation, error) {

	var recommendation *resource.PodResourceRecommendation

	controllerRecommendation := ac.getControllerRecommendation(controllerID)
	controllerLock := ac.getControllerLock(controllerID)

	retryTime := ac.resourceRecommendatorSyncRetryTime
	controllerLock.Lock()
	defer controllerLock.Unlock()
	for recommendation == nil && retryTime > 0 {
		if newRecommendations, err := ac.fetchNewPodRecommendations(controllerID); err != nil {
			scope.Warnf("fetch new recommendation failed, retry fetching, errMsg: %s", err.Error())
		} else {
			controllerRecommendation.setPodRecommendations(newRecommendations)
			break
		}
		retryTime--
	}
	validRecommedations, err := ac.listValidPodRecommendations(controllerID, controllerRecommendation.getPodRecommendations())
	if err != nil {
		return nil, err
	}
	controllerRecommendation.setPodRecommendations(validRecommedations)
	recommendation = controllerRecommendation.dispatchOneValidPodRecommendation(time.Now())

	return recommendation, nil
}

func (ac *admissionController) getControllerRecommendation(controllerID namespaceKindName) *controllerRecommendation {

	ac.lock.Lock()
	controllerRecommendation, exist := ac.controllerRecommendationMap[controllerID]
	if !exist {
		scope.Debugf("controllerID: %s, controller recommendation not exist, create new recommendation.", controllerID)
		ac.controllerRecommendationMap[controllerID] = NewControllerPodResourceRecommendation()
		controllerRecommendation = ac.controllerRecommendationMap[controllerID]
	}
	ac.lock.Unlock()

	return controllerRecommendation
}

func (ac *admissionController) getControllerLock(controllerID namespaceKindName) *sync.Mutex {

	ac.lock.Lock()
	lock, exist := ac.controllerLockMap[controllerID]
	if !exist {
		ac.controllerLockMap[controllerID] = &sync.Mutex{}
		lock = ac.controllerLockMap[controllerID]
	}
	ac.lock.Unlock()

	return lock
}

func (ac *admissionController) listValidPodRecommendations(controllerID namespaceKindName, recommendations []*resource.PodResourceRecommendation) ([]*resource.PodResourceRecommendation, error) {

	validRecommendations := make([]*resource.PodResourceRecommendation, 0)

	podRecommendationNumberMap := buildPodRecommendationNumberMap(recommendations)
	scope.Debugf("list valid pod recommdendations: controllerID: %s, podRecommendationNumberMap %+v", controllerID.String(), podRecommendationNumberMap)

	pods, err := ac.listPodByController(controllerID)
	if err != nil {
		return validRecommendations, errors.Wrapf(err, "list valid recommendations failed, controllerID: %s", controllerID.String())
	}
	removeApplyingPodRecommendations(podRecommendationNumberMap, pods)

	scope.Debugf("valid podRecommendationNumberMap for controllerID: %s, podRecommendationNumberMap %+v", controllerID.String(), podRecommendationNumberMap)
	validRecommendations = listValidPodRecommedationsFromRecommendationNumberMap(podRecommendationNumberMap, recommendations)

	return validRecommendations, nil
}

func (ac *admissionController) listPodByController(controllerID namespaceKindName) ([]*core_v1.Pod, error) {
	pods := make([]*core_v1.Pod, 0)

	var err error
	podsInCluster := make([]core_v1.Pod, 0)
	listResource := resources.NewListResources(ac.sigsK8SClient)
	switch controllerID.getKind() {
	case "Deployment":
		podsInCluster, err = listResource.ListPodsByDeployment(controllerID.getNamespace(), controllerID.getName())
		if err != nil {
			return pods, errors.Wrapf(err, "list pods controlled by controllerID: %s failed", controllerID.String())
		}
	case "DeploymentConfig":
		podsInCluster, err = listResource.ListPodsByDeploymentConfig(controllerID.getNamespace(), controllerID.getName())
		if err != nil {
			return pods, errors.Wrapf(err, "list pods controlled by controllerID: %s failed", controllerID.String())
		}
	default:
		return pods, errors.Errorf("no matching resource lister for controller kind: %s", controllerID.getKind())
	}

	for _, pod := range podsInCluster {
		copyPod := pod
		pods = append(pods, &copyPod)
	}

	return pods, nil
}

func (ac *admissionController) fetchNewPodRecommendations(controllerID namespaceKindName) ([]*resource.PodResourceRecommendation, error) {

	scope.Debugf("fetching new recommendations from recommendator, controllerID: %s", controllerID.String())

	var err error
	recommendations := make([]*resource.PodResourceRecommendation, 0)
	done := make(chan bool)

	go func(chan bool) {
		queryTime := time.Now()
		recommendations, err = ac.resourceRecommendator.ListControllerPodResourceRecommendations(resource.ListControllerPodResourceRecommendationsRequest{
			Namespace: controllerID.getNamespace(),
			Name:      controllerID.getName(),
			Kind:      controllerID.getKind(),
			Time:      &queryTime,
		})
		done <- true
	}(done)

	select {
	case _ = <-done:
	case _ = <-time.After(ac.resourceRecommendatorSyncTimeout):
		err = errors.Errorf("fetch recommendations failed: controllerID: %s, timeout after %f seconds", controllerID.String(), ac.resourceRecommendatorSyncTimeout.Seconds())
	}

	return recommendations, err
}

func (ac *admissionController) isControllerExecutionEnabled(controllerID namespaceKindName) (bool, error) {

	return ac.controllerValidator.IsControllerEnabledExecution(controllerID.namespace, controllerID.name, controllerID.kind)
}

func (ac *admissionController) getControllerIDToQueryPodRecommendations(pod *core_v1.Pod) (namespaceKindName, error) {

	var controllerID = namespaceKindName{}

	link, err := ac.ownerReferenceTracer.GetControllerOwnerReferenceLink(pod)
	if err != nil {
		return controllerID, err
	}

	for i := len(link) - 1; i >= 0; i-- {
		ownerRef := link[i]
		if _, exist := autoscalingv1alpha1.K8SKindToAlamedaControllerType[ownerRef.Kind]; exist {
			controllerID.namespace = pod.Namespace
			controllerID.name = ownerRef.Name
			controllerID.kind = ownerRef.Kind
			break
		}
	}

	return controllerID, nil
}

func buildPodRecommendationNumberMap(recommendations []*resource.PodResourceRecommendation) map[string]int {
	currentTime := time.Now()
	recommendationNumberMap := make(map[string]int)
	for _, recommendation := range recommendations {
		if !(recommendation.ValidStartTime.Unix() < currentTime.Unix() && currentTime.Unix() < recommendation.ValidEndTime.Unix()) {
			continue
		}
		recommendationID := buildPodResourceIDFromPodRecommendation(recommendation)
		recommendationNumberMap[recommendationID]++
	}
	return recommendationNumberMap
}

func removeApplyingPodRecommendations(recommendationNumberMap map[string]int, pods []*core_v1.Pod) {
	for _, pod := range pods {
		scope.Debugf("try to decrease recommendation from pod: %s/%s", pod.Namespace, pod.Name)
		if !alamedascaler_reconciler.PodIsMonitoredByAlameda(pod) {
			scope.Debugf("skip decreasing recommendation cause pod's %s/%s phase: %s is not monitored by Alameda", pod.Namespace, pod.Name, pod.Status.Phase)
			continue
		}
		recommendationID := buildPodResourceIDFromPod(pod)
		if _, exist := recommendationNumberMap[recommendationID]; exist {
			scope.Debugf("decrease recommendation for pod %s/%s", pod.Namespace, pod.Name)
			recommendationNumberMap[recommendationID]--
		}
	}
}

func listValidPodRecommedationsFromRecommendationNumberMap(recommendationNumberMap map[string]int, recommendations []*resource.PodResourceRecommendation) []*resource.PodResourceRecommendation {

	validRecommendations := make([]*resource.PodResourceRecommendation, 0)
	for _, recommendation := range recommendations {
		copyRecommendation := recommendation
		recommendationID := buildPodResourceIDFromPodRecommendation(recommendation)
		if remainRecommendationsNum := recommendationNumberMap[recommendationID]; remainRecommendationsNum > 0 {
			recommendationNumberMap[recommendationID]--
			validRecommendations = append(validRecommendations, copyRecommendation)
		}
	}
	return validRecommendations
}

func buildPodResourceIDFromPod(pod *core_v1.Pod) string {

	containers := pod.Spec.Containers

	sort.SliceStable(containers, func(i, j int) bool {
		return containers[i].Name < containers[j].Name
	})

	id := ""
	for _, container := range containers {
		requestCPU := container.Resources.Requests.Cpu().MilliValue()
		requestMem := container.Resources.Requests.Memory().Value()
		limitsCPU := container.Resources.Limits.Cpu().MilliValue()
		limitsMem := container.Resources.Limits.Memory().Value()
		id += fmt.Sprintf("container-name-%s/requset-cpu-%d-mem-%d/limit-cpu-%d-mem-%d/", container.Name,
			requestCPU, requestMem,
			limitsCPU, limitsMem,
		)
	}

	return id
}

func buildPodResourceIDFromPodRecommendation(recommendation *resource.PodResourceRecommendation) string {

	containerRecommendations := recommendation.ContainerResourceRecommendations
	sort.SliceStable(containerRecommendations, func(i, j int) bool {
		return containerRecommendations[i].Name < containerRecommendations[j].Name
	})

	id := ""
	for _, containerRecommendation := range containerRecommendations {
		requestCPU := containerRecommendation.Requests.Cpu().MilliValue()
		requestMem := containerRecommendation.Requests.Memory().Value()
		limitsCPU := containerRecommendation.Limits.Cpu().MilliValue()
		limitsMem := containerRecommendation.Limits.Memory().Value()
		id += fmt.Sprintf("container-name-%s/requset-cpu-%d-mem-%d/limit-cpu-%d-mem-%d/", containerRecommendation.Name,
			requestCPU, requestMem,
			limitsCPU, limitsMem,
		)
	}
	return id
}
