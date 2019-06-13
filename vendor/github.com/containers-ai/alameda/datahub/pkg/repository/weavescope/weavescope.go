package weavescope

import (
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"

	"io/ioutil"
	"net/http"

	"crypto/tls"
	"fmt"
)

func NewWeaveScopeRepositoryWithConfig(config *Config) *WeaveScopeRepository {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	return &WeaveScopeRepository{
		URL: config.URL,
	}
}

type WeaveScopeRepository struct {
	URL string
}

func (w *WeaveScopeRepository) ListWeaveScopeHosts(in *datahub_v1alpha1.ListWeaveScopeHostsRequest) (string, error) {
	url := fmt.Sprintf("%s%s", w.URL, "/api/topology/hosts")

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	readBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(readBody), nil
}

func (w *WeaveScopeRepository) GetWeaveScopeHostDetails(in *datahub_v1alpha1.ListWeaveScopeHostsRequest) (string, error) {
	url := fmt.Sprintf("%s%s%s;<host>", w.URL, "/api/topology/hosts/", in.GetHostId())

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	readBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(readBody), nil
}

func (w *WeaveScopeRepository) ListWeaveScopePods(in *datahub_v1alpha1.ListWeaveScopePodsRequest) (string, error) {
	url := fmt.Sprintf("%s%s", w.URL, "/api/topology/pods")

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	readBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(readBody), nil
}

func (w *WeaveScopeRepository) GetWeaveScopePodDetails(in *datahub_v1alpha1.ListWeaveScopePodsRequest) (string, error) {
	url := fmt.Sprintf("%s%s%s;<pod>", w.URL, "/api/topology/pods/", in.GetPodId())

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	readBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(readBody), nil
}

func (w *WeaveScopeRepository) ListWeaveScopeContainers(in *datahub_v1alpha1.ListWeaveScopeContainersRequest) (string, error) {
	url := fmt.Sprintf("%s%s", w.URL, "/api/topology/containers")

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	readBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(readBody), nil
}

func (w *WeaveScopeRepository) ListWeaveScopeContainersByHostname(in *datahub_v1alpha1.ListWeaveScopeContainersRequest) (string, error) {
	url := fmt.Sprintf("%s%s", w.URL, "/api/topology/containers-by-hostname")

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	readBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(readBody), nil
}

func (w *WeaveScopeRepository) ListWeaveScopeContainersByImage(in *datahub_v1alpha1.ListWeaveScopeContainersRequest) (string, error) {
	url := fmt.Sprintf("%s%s", w.URL, "/api/topology/containers-by-image")

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	readBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(readBody), nil
}

func (w *WeaveScopeRepository) GetWeaveScopeContainerDetails(in *datahub_v1alpha1.ListWeaveScopeContainersRequest) (string, error) {
	url := fmt.Sprintf("%s%s%s;<container>", w.URL, "/api/topology/containers/", in.GetContainerId())

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	readBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(readBody), nil
}
