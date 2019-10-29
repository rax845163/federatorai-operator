package alamedaservicekeycode

import (
	datahubv1alpha1_event "github.com/containers-ai/api/alameda_api/v1alpha1/datahub/events"
	"github.com/golang/protobuf/ptypes"
	"k8s.io/apimachinery/pkg/util/uuid"
)

const (
	componentName = "federatorai-operator"
)

func newLicenseEvent(namespace, message, clusterID string, level datahubv1alpha1_event.EventLevel) datahubv1alpha1_event.Event {

	now := ptypes.TimestampNow()
	id := uuid.NewUUID()
	source := datahubv1alpha1_event.EventSource{
		Host:      "",
		Component: componentName,
	}
	eventType := datahubv1alpha1_event.EventType_EVENT_TYPE_LICENSE
	version := datahubv1alpha1_event.EventVersion_EVENT_VERSION_V1
	subject := datahubv1alpha1_event.K8SObjectReference{
		Kind:       "Pod",
		Namespace:  namespace,
		Name:       "Federator.ai",
		ApiVersion: "v1",
	}
	data := ""

	event := datahubv1alpha1_event.Event{
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
