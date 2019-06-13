package kubernetes

import (
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var Codecs = serializer.NewCodecFactory(Scheme)