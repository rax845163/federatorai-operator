package prometheus

import (
	"encoding/json"
	"errors"
	"fmt"
	DAO "github.com/containers-ai/alameda/datahub/pkg/dao"
	Common "github.com/containers-ai/api/common"
	"github.com/golang/protobuf/ptypes"
	"time"
)

var (
	daoAggregateFunction_PrometheusAggregationOverTimeFunction = map[DAO.AggregateFunction]string{
		DAO.Max: "max_over_time",
	}
)

var (
	DatahubAggregateFunction_DAOAggregateFunction = map[Common.TimeRange_AggregateFunction]DAO.AggregateFunction{
		Common.TimeRange_NONE: DAO.None,
		Common.TimeRange_MAX:  DAO.Max,
	}
)

func ReadRawdata(config *Config, queries []*Common.Query) ([]*Common.ReadRawdata, error) {
	rawdata := make([]*Common.ReadRawdata, 0)

	prometheusClient, err := New(*config)
	if err != nil {
		scope.Errorf("failed to read rawdata from Prometheus: %v", err)
		return make([]*Common.ReadRawdata, 0), errors.New("failed to instance prometheus client")
	}

	for _, query := range queries {
		response := Response{}
		err := errors.New("")

		queryExpression := ""
		queryCondition := BuildQueryCondition(query)

		options := []Option{
			StartTime(queryCondition.StartTime),
			EndTime(queryCondition.EndTime),
			Timeout(queryCondition.Timeout),
			StepDuration(queryCondition.StepTime),
			AggregateOverTimeFunction(queryCondition.AggregateOverTimeFunction),
		}

		opt := BuildDefaultOptions()
		for _, option := range options {
			option(&opt)
		}

		if query.GetCondition().GetWhereClause() != "" {
			queryExpression = fmt.Sprintf("%s{%s}", query.GetTable(), query.GetCondition().GetWhereClause())
		} else {
			queryExpression = fmt.Sprintf("%s", query.GetTable())
		}

		if query.GetCondition().GetTimeRange().GetStep() != nil {
			stepTimeInSeconds := int64(opt.stepTime.Nanoseconds() / int64(time.Second))
			queryExpression, err = wrapQueryExpressionWithAggregationOverTimeFunction(queryExpression, opt.aggregateOverTimeFunc, stepTimeInSeconds)
			if err != nil {
				return make([]*Common.ReadRawdata, 0), errors.New(err.Error())
			}
		}

		if query.GetExpression() == "query" {
			response, err = prometheusClient.Query(queryExpression, opt.startTime, opt.timeout)
		} else if query.GetExpression() == "query_range" {
			response, err = prometheusClient.QueryRange(queryExpression, opt.startTime, opt.endTime, opt.stepTime)
		} else {
			response, err = prometheusClient.QueryRange(queryExpression, opt.startTime, opt.endTime, opt.stepTime)
		}

		if err != nil {
			return make([]*Common.ReadRawdata, 0), errors.New(err.Error())
		} else if response.Status != StatusSuccess {
			scope.Errorf("receive error response from prometheus: %s", response.Error)
			return make([]*Common.ReadRawdata, 0), errors.New(response.Error)
		} else {
			readRawdata := PrometheusResponseToReadRawdata(&response, query)
			rawdata = append(rawdata, readRawdata)
		}
	}

	return rawdata, nil
}

func BuildQueryCondition(query *Common.Query) DAO.QueryCondition {
	var (
		queryStartTime      *time.Time
		queryEndTime        *time.Time
		queryTimeout        *time.Time
		queryStepTime       *time.Duration
		queryTimestampOrder int
		queryLimit          int
		queryCondition      = DAO.QueryCondition{}
		aggregateFunc       = DAO.None
	)

	if query.GetCondition() == nil {
		return queryCondition
	}

	if query.GetCondition().GetTimeRange() != nil {
		timeRange := query.GetCondition().GetTimeRange()

		if timeRange.GetStartTime() != nil {
			tmpTime, _ := ptypes.Timestamp(timeRange.GetStartTime())
			queryStartTime = &tmpTime
		}

		if timeRange.GetEndTime() != nil {
			tmpTime, _ := ptypes.Timestamp(timeRange.GetEndTime())
			queryEndTime = &tmpTime
		}

		if timeRange.GetTimeout() != nil {
			tmpTime, _ := ptypes.Timestamp(timeRange.GetTimeout())
			queryTimeout = &tmpTime
		}

		if timeRange.GetStep() != nil {
			tmpTime, _ := ptypes.Duration(timeRange.GetStep())
			queryStepTime = &tmpTime
		}

		switch query.GetCondition().GetOrder() {
		case Common.QueryCondition_ASC:
			queryTimestampOrder = DAO.Asc
		case Common.QueryCondition_DESC:
			queryTimestampOrder = DAO.Desc
		default:
			queryTimestampOrder = DAO.Asc
		}

		queryLimit = int(query.GetCondition().GetLimit())
	}

	queryTimestampOrder = int(query.GetCondition().GetOrder())
	queryLimit = int(query.GetCondition().GetLimit())

	if aggFunc, exist := DatahubAggregateFunction_DAOAggregateFunction[Common.TimeRange_AggregateFunction(query.GetCondition().GetTimeRange().GetAggregateFunction())]; exist {
		aggregateFunc = aggFunc
	}

	queryCondition = DAO.QueryCondition{
		StartTime:      queryStartTime,
		EndTime:        queryEndTime,
		Timeout:        queryTimeout,
		StepTime:       queryStepTime,
		TimestampOrder: queryTimestampOrder,
		Limit:          queryLimit,
		AggregateOverTimeFunction: aggregateFunc,
	}

	return queryCondition
}

func PrometheusResponseToReadRawdata(response *Response, query *Common.Query) *Common.ReadRawdata {
	readRawdata := Common.ReadRawdata{Query: query}

	if len(response.Data.Result) == 0 {
		return &readRawdata
	}
	jsonStr, err := json.Marshal(response.Data)
	if err != nil {
		scope.Errorf("failed to Marshal response from Prometheus: %v", err)
	}

	readRawdata.Rawdata = string(jsonStr)

	return &readRawdata
}

func wrapQueryExpressionWithAggregationOverTimeFunction(queryExpression string, aggregateFunc DAO.AggregateFunction, aggregateOverSeconds int64) (string, error) {

	if aggregateFunc == DAO.None {
		return queryExpression, nil
	}

	if funcName, exist := daoAggregateFunction_PrometheusAggregationOverTimeFunction[aggregateFunc]; !exist {
		return queryExpression, errors.New(fmt.Sprintf("wrap prometheus query expression with function failed: no mapping function for function: %d", aggregateFunc))
	} else {
		queryExpression = fmt.Sprintf("%s(%s[%ds])", funcName, queryExpression, aggregateOverSeconds)
	}

	return queryExpression, nil
}
