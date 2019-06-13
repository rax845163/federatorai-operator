package influxdb

import (
	"fmt"
	"strings"
	"time"
)

// Order Order enumerator
type Order = int

const (
	// Asc Represent ascending order
	Asc Order = 0
	// Desc Represent descending order
	Desc Order = 1
)

// QueryCondition Others query condition
type QueryCondition struct {
	StartTime      *time.Time
	EndTime        *time.Time
	StepTime       *time.Duration
	TimestampOrder Order
	Limit          int
}

type Statement struct {
	Measurement    Measurement
	SelectedFields []string
	GroupByTags    []string
	WhereClause    string
	OrderClause    string
	LimitClause    string
}

func (s *Statement) AppendTimeConditionIntoWhereClause(queryCondition QueryCondition) {

	var (
		reqStartTime = queryCondition.StartTime
		reqEndTime   = queryCondition.EndTime
	)

	timeConditionStr := ""
	if reqStartTime != nil && reqEndTime != nil {
		timeConditionStr = fmt.Sprintf("time >= %v AND time <= %v", reqStartTime.UnixNano(), reqEndTime.UnixNano())
	} else if reqStartTime != nil && reqEndTime == nil {
		timeConditionStr = fmt.Sprintf("time >= %v", reqStartTime.UnixNano())
	} else if reqStartTime == nil && reqEndTime != nil {
		timeConditionStr = fmt.Sprintf("time <= %v", reqEndTime.UnixNano())
	}

	if s.WhereClause == "" && timeConditionStr != "" {
		s.WhereClause = fmt.Sprintf("WHERE %s", timeConditionStr)
	} else if s.WhereClause != "" && timeConditionStr != "" {
		s.WhereClause = fmt.Sprintf("%s AND %s", s.WhereClause, timeConditionStr)
	}

}

func (s *Statement) SetOrderClauseFromQueryCondition(queryCondition QueryCondition) {

	switch queryCondition.TimestampOrder {
	case Asc:
		s.OrderClause = "ORDER BY time ASC"
	case Desc:
		s.OrderClause = "ORDER BY time DESC"
	default:
		s.OrderClause = "ORDER BY time ASC"
	}
}

func (s *Statement) SetLimitClauseFromQueryCondition(queryCondition QueryCondition) {

	limit := queryCondition.Limit
	if limit > 0 {
		s.LimitClause = fmt.Sprintf("LIMIT %v", limit)
	}
}

func (s Statement) BuildQueryCmd() string {

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

	cmd = fmt.Sprintf("SELECT %s FROM %s %s %s %s %s",
		fieldsStr, s.Measurement, s.WhereClause,
		groupByStr, s.OrderClause, s.LimitClause)

	return cmd
}
