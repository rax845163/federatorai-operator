package utils

import (
	"encoding/json"
	"strings"

	"github.com/containers-ai/alameda/admission-controller/pkg/recommendator/resource"
	"github.com/mattbaird/jsonpatch"
	"github.com/pkg/errors"
	core_v1 "k8s.io/api/core/v1"
)

func GetPatchesFromPodResourceRecommendation(pod *core_v1.Pod, recommendation *resource.PodResourceRecommendation) ([]jsonpatch.JsonPatchOperation, error) {

	patches := make([]jsonpatch.JsonPatchOperation, 0)

	originPod := pod.DeepCopy()
	mutatedPod := pod.DeepCopy()

	containerNameMap := make(map[string]int)
	for i, container := range mutatedPod.Spec.Containers {
		containerNameMap[container.Name] = i
	}

	for _, containerResourceRecommendation := range recommendation.ContainerResourceRecommendations {

		containerName := containerResourceRecommendation.Name
		containerIndex, exist := containerNameMap[containerName]
		if !exist {
			continue
		}

		mutatedPod.Spec.Containers[containerIndex].Resources.Limits = containerResourceRecommendation.Limits
		mutatedPod.Spec.Containers[containerIndex].Resources.Requests = containerResourceRecommendation.Requests
	}

	originPodbytes, err := json.Marshal(originPod)
	if err != nil {
		return patches, errors.Errorf("get patch bytes failed: %s", err.Error())
	}
	mutatedPodbytes, err := json.Marshal(mutatedPod)
	if err != nil {
		return patches, errors.Errorf("get patch bytes failed: %s", err.Error())
	}

	patches, err = jsonpatch.CreatePatch([]byte(originPodbytes), []byte(mutatedPodbytes))
	if err != nil {
		return patches, errors.Errorf("Error creating JSON patch: %s", err.Error())
	}

	return patches, nil
}

type ValidatePatchFunc func(patch jsonpatch.JsonPatchOperation) error

func ValidatePatches(patches []jsonpatch.JsonPatchOperation, validateFunc ValidatePatchFunc) error {

	for _, patch := range patches {
		if err := validateFunc(patch); err != nil {
			return errors.Wrapf(err, "json patch validate failed")
		}
	}
	return nil
}

func GetK8SPatchesString(patches []jsonpatch.JsonPatchOperation) string {

	patch := ""
	for _, operation := range patches {
		patch += operation.Json() + ","
	}
	patch = strings.TrimSuffix(patch, ",")
	patch = "[" + patch + "]"
	return patch
}
