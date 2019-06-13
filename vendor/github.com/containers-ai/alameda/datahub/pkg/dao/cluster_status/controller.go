package clusterstatus

import (
	datahub_api "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
)

type ControllerOperation interface {
	CreateControllers([]*datahub_api.Controller) error
	ListControllers(*datahub_api.ListControllersRequest) ([]*datahub_api.Controller, error)
	DeleteControllers(*datahub_api.DeleteControllersRequest) error
}
