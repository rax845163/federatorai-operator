package alamedaservicekeycode

import (
	"github.com/golang/protobuf/ptypes"

	datahubv1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	"k8s.io/apimachinery/pkg/util/uuid"
)

const (
	componentName = "federatorai-operator"
)

func newLicenseEvent(namespace, message, clusterID string, level datahubv1alpha1.EventLevel) datahubv1alpha1.Event {

	now := ptypes.TimestampNow()
	id := uuid.NewUUID()
	source := datahubv1alpha1.EventSource{
		Host:      "",
		Component: componentName,
	}
	eventType := datahubv1alpha1.EventType_EVENT_TYPE_LICENSE
	version := datahubv1alpha1.EventVersion_EVENT_VERSION_V1
	subject := datahubv1alpha1.K8SObjectReference{
		Kind:       "Pod",
		Namespace:  namespace,
		Name:       "Federator.ai",
		ApiVersion: "v1",
	}
	data := ""

	event := datahubv1alpha1.Event{
		Time:      now,
		Id:        string(id),
		ClusterId: clusterID,
		Source:    &source,
		Type:      eventType,
		Version:   version,
		Level:     level,
		Subject:   &subject,
		Message:   message,
		Data:      data,
	}

	return event
}
