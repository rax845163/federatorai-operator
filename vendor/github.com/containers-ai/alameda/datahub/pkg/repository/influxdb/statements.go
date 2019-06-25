package influxdb

import (
	"fmt"
	Common "github.com/containers-ai/api/common"
	"strings"
	"time"
)

type InfluxStatement struct {
	QueryCondition *QueryCondition
	Database       Database
	Measurement    Measurement
	SelectedFields []string
	GroupByTags    []string
	WhereClause    string
	OrderClause    string
	LimitClause    string
}

func NewInfluxStatement(query *Common.Query) *InfluxStatement {
	queryCondition := QueryCondition{}

	if query.GetCondition().GetTimeRange().GetStartTime() != nil {
		startTime := time.Unix(query.GetCondition().GetTimeRange().GetStartTime().GetSeconds(), 0)
		queryCondition.StartTime = &startTime
	}
	if query.GetCondition().GetTimeRange().GetEndTime() != nil {
		endTime := time.Unix(query.GetCondition().GetTimeRange().GetEndTime().GetSeconds(), 0)
		queryCondition.EndTime = &endTime
	}
	if query.GetCondition().GetTimeRange().GetStep() != nil {
		stepTime := time.Duration(query.GetCondition().GetTimeRange().GetStep().GetSeconds())
		queryCondition.StepTime = &stepTime
	}
	queryCondition.TimestampOrder = Order(query.GetCondition().GetOrder())
	queryCondition.Limit = int(query.GetCondition().GetLimit())

	statement := InfluxStatement{
		QueryCondition: &queryCondition,
		Database:       Database(query.GetDatabase()),
		Measurement:    Measurement(query.GetTable()),
		SelectedFields: query.GetCondition().GetSelects(),
		GroupByTags:    query.GetCondition().GetGroups(),
		WhereClause:    query.GetCondition().GetWhereClause(),
	}

	return &statement
}

func (s *InfluxStatement) AppendTimeConditionIntoWhereClause() {
	timeConditionStr := ""
	if s.QueryCondition.StartTime != nil && s.QueryCondition.EndTime != nil {
		timeConditionStr = fmt.Sprintf("time >= %v AND time <= %v", s.QueryCondition.StartTime.UnixNano(), s.QueryCondition.EndTime.UnixNano())
	} else if s.QueryCondition.StartTime != nil && s.QueryCondition.EndTime == nil {
		timeConditionStr = fmt.Sprintf("time >= %v", s.QueryCondition.StartTime.UnixNano())
	} else if s.QueryCondition.StartTime == nil && s.QueryCondition.EndTime != nil {
		timeConditionStr = fmt.Sprintf("time <= %v", s.QueryCondition.EndTime.UnixNano())
	}

	if s.WhereClause == "" && timeConditionStr != "" {
		s.WhereClause = fmt.Sprintf("WHERE %s", timeConditionStr)
	} else if s.WhereClause != "" && timeConditionStr != "" {
		s.WhereClause = fmt.Sprintf("%s AND %s", s.WhereClause, timeConditionStr)
	}

}

func (s *InfluxStatement) SetOrderClauseFromQueryCondition() {
	switch s.QueryCondition.TimestampOrder {
	case Asc:
		s.OrderClause = "ORDER BY time ASC"
	case Desc:
		s.OrderClause = "ORDER BY time DESC"
	default:
		s.OrderClause = "ORDER BY time ASC"
	}
}

func (s *InfluxStatement) SetLimitClauseFromQueryCondition() {
	limit := s.QueryCondition.Limit
	if limit > 0 {
		s.LimitClause = fmt.Sprintf("LIMIT %v", limit)
	}
}

func (s InfluxStatement) BuildQueryCmd() string {
	var (
		cmd        = ""
		fieldsStr  = "*"
		groupByStr = ""
	)

	if len(s.SelectedFields) > 0 {
		fieldsStr = ""
		for _, field := range s.SelectedFields {
			fieldsStr += fmt.Sprintf(`"%s",`, field)
		}
		fieldsStr = strings.TrimSuffix(fieldsStr, ",")
	}

	if len(s.GroupByTags) > 0 {
		groupByStr = "GROUP BY "
		for _, field := range s.GroupByTags {
			groupByStr += fmt.Sprintf(`"%s",`, field)
		}
		groupByStr = strings.TrimSuffix(groupByStr, ",")
	}

	cmd = fmt.Sprintf("SELECT %s FROM %s %s %s %s %s", fieldsStr, s.Measurement, s.WhereClause, groupByStr, s.OrderClause, s.LimitClause)

	return cmd
}
