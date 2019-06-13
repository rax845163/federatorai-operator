package weavescope

import (
	weaveScope_repository "github.com/containers-ai/alameda/datahub/pkg/repository/weavescope"

	//influxdb_repository_recommendation "github.com/containers-ai/alameda/datahub/pkg/repository/influxdb/recommendation"
	"github.com/containers-ai/alameda/pkg/utils/log"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
)

var (
	scope = log.RegisterScope("recommendation_dao_implement", "recommended dao implement", 0)
)

// Container Implements ContainerOperation interface
type WeaveScope struct {
	WeaveScopeConfig *weaveScope_repository.Config
}

func (w *WeaveScope) ListWeaveScopeHosts(in *datahub_v1alpha1.ListWeaveScopeHostsRequest) (string, error) {
	weaveScopeRepository := weaveScope_repository.NewWeaveScopeRepositoryWithConfig(w.WeaveScopeConfig)
	return weaveScopeRepository.ListWeaveScopeHosts(in)
}

func (w *WeaveScope) GetWeaveScopeHostDetails(in *datahub_v1alpha1.ListWeaveScopeHostsRequest) (string, error) {
	weaveScopeRepository := weaveScope_repository.NewWeaveScopeRepositoryWithConfig(w.WeaveScopeConfig)
	return weaveScopeRepository.GetWeaveScopeHostDetails(in)
}

func (w *WeaveScope) ListWeaveScopePods(in *datahub_v1alpha1.ListWeaveScopePodsRequest) (string, error) {
	weaveScopeRepository := weaveScope_repository.NewWeaveScopeRepositoryWithConfig(w.WeaveScopeConfig)
	return weaveScopeRepository.ListWeaveScopePods(in)
}

func (w *WeaveScope) GetWeaveScopePodDetails(in *datahub_v1alpha1.ListWeaveScopePodsRequest) (string, error) {
	weaveScopeRepository := weaveScope_repository.NewWeaveScopeRepositoryWithConfig(w.WeaveScopeConfig)
	return weaveScopeRepository.GetWeaveScopePodDetails(in)
}

func (w *WeaveScope) ListWeaveScopeContainers(in *datahub_v1alpha1.ListWeaveScopeContainersRequest) (string, error) {
	weaveScopeRepository := weaveScope_repository.NewWeaveScopeRepositoryWithConfig(w.WeaveScopeConfig)
	return weaveScopeRepository.ListWeaveScopeContainers(in)
}

func (w *WeaveScope) ListWeaveScopeContainersByHostname(in *datahub_v1alpha1.ListWeaveScopeContainersRequest) (string, error) {
	weaveScopeRepository := weaveScope_repository.NewWeaveScopeRepositoryWithConfig(w.WeaveScopeConfig)
	return weaveScopeRepository.ListWeaveScopeContainersByHostname(in)
}

func (w *WeaveScope) ListWeaveScopeContainersByImage(in *datahub_v1alpha1.ListWeaveScopeContainersRequest) (string, error) {
	weaveScopeRepository := weaveScope_repository.NewWeaveScopeRepositoryWithConfig(w.WeaveScopeConfig)
	return weaveScopeRepository.ListWeaveScopeContainersByImage(in)
}

func (w *WeaveScope) GetWeaveScopeContainerDetails(in *datahub_v1alpha1.ListWeaveScopeContainersRequest) (string, error) {
	weaveScopeRepository := weaveScope_repository.NewWeaveScopeRepositoryWithConfig(w.WeaveScopeConfig)
	return weaveScopeRepository.GetWeaveScopeContainerDetails(in)
}
