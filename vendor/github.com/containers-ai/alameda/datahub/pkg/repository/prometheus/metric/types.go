package metric

import (
	"github.com/containers-ai/alameda/datahub/pkg/repository/prometheus"
)

type nodeMetricsFetchingFunction func(nodeName string, options ...Option) ([]prometheus.Entity, error)
