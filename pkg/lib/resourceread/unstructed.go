package resourceread

import (
	unstructuredv1 "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	unstructuredsScheme = runtime.NewScheme()
	unstructuredsCodecs = serializer.NewCodecFactory(unstructuredsScheme)
)

// ReadJSONBytes converts the received JSON bytes to Unstructed type.
func ReadJSONBytes(objBytes []byte) (*unstructuredv1.Unstructured, error) {
	obj := unstructuredv1.Unstructured{}
	_, _, err := unstructuredv1.UnstructuredJSONScheme.Decode(objBytes, nil, &obj)
	return &obj, err
}
