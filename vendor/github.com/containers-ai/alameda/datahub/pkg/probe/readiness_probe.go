package probe

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/containers-ai/alameda/datahub/pkg/repository/prometheus"
	promRepository "github.com/containers-ai/alameda/datahub/pkg/repository/prometheus/metric"
	"github.com/pkg/errors"
)

type ReadinessProbeConfig struct {
	InfluxdbAddr  string
	PrometheusCfg *prometheus.Config
}

func pingInfluxdb(influxdbAddr string) error {
	pingURL := fmt.Sprintf("%s/ping", influxdbAddr)
	curlCmd := exec.Command("curl", "-sl", "-I", pingURL)
	if strings.Contains(pingURL, "https") {
		curlCmd = exec.Command("curl", "-sl", "-I", "-k", pingURL)
	}
	_, err := curlCmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

func queryPrometheus(prometheusConfig *prometheus.Config) error {
	emr := "exceeded maximum resolution"
	options := []promRepository.Option{}

	podContainerCPURepo := promRepository.NewPodContainerCPUUsagePercentageRepositoryWithConfig(*prometheusConfig)
	containerCPUEntities, err := podContainerCPURepo.ListMetricsByPodNamespacedName("", "", options...)
	if err != nil && !strings.Contains(err.Error(), emr) {
		return errors.Wrap(err, "list pod metrics failed")
	}

	if err == nil && len(containerCPUEntities) == 0 {
		return fmt.Errorf("No container CPU metric found")
	}

	podContainerMemoryRepo := promRepository.NewPodContainerMemoryUsageBytesRepositoryWithConfig(*prometheusConfig)
	containerMemoryEntities, err := podContainerMemoryRepo.ListMetricsByPodNamespacedName("", "", options...)
	if err != nil && !strings.Contains(err.Error(), emr) {
		return errors.Wrap(err, "list pod metrics failed")
	}

	if err == nil && len(containerMemoryEntities) == 0 {
		return fmt.Errorf("No container memory metric found")
	}

	nodeCPUUsageRepo := promRepository.NewNodeCPUUsagePercentageRepositoryWithConfig(*prometheusConfig)
	nodeCPUUsageEntities, err := nodeCPUUsageRepo.ListMetricsByNodeName("", options...)
	if err != nil && !strings.Contains(err.Error(), emr) {
		return errors.Wrap(err, "list node cpu usage metrics failed")
	}

	if err == nil && len(nodeCPUUsageEntities) == 0 {
		return fmt.Errorf("No node CPU metric found")
	}

	nodeMemoryUsageRepo := promRepository.NewNodeMemoryUsageBytesRepositoryWithConfig(*prometheusConfig)
	nodeMemoryUsageEntities, err := nodeMemoryUsageRepo.ListMetricsByNodeName("", options...)
	if err != nil && !strings.Contains(err.Error(), emr) {
		return errors.Wrap(err, "list node memory usage metrics failed")
	}

	if err == nil && len(nodeMemoryUsageEntities) == 0 {
		return fmt.Errorf("No node memory metric found")
	}

	return nil
}
