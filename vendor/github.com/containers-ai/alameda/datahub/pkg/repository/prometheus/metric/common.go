package metric

import (
	"fmt"
	"time"

	"github.com/containers-ai/alameda/datahub/pkg/dao"
	"github.com/pkg/errors"
)

var (
	defaultStepTime, _ = time.ParseDuration("30s")
)

var (
	daoAggregateFunction_PrometheusAggregationOverTimeFunction = map[dao.AggregateFunction]string{
		dao.Max: "max_over_time",
	}
)

func wrapQueryExpressionWithAggregationOverTimeFunction(queryExpression string, aggregateFunc dao.AggregateFunction, aggregateOverSeconds int64) (string, error) {

	if aggregateFunc == dao.None {
		return queryExpression, nil
	}

	if funcName, exist := daoAggregateFunction_PrometheusAggregationOverTimeFunction[aggregateFunc]; !exist {
		return queryExpression, errors.Errorf("wrap prometheus query expression with function failed: no mapping function for function: %d", aggregateFunc)
	} else {
		queryExpression = fmt.Sprintf("%s(%s[%ds])", funcName, queryExpression, aggregateOverSeconds)
	}

	return queryExpression, nil
}
