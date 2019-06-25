package influxdb

import (
	"fmt"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	"strings"
	"time"
)

type StatementNew struct {
	Measurement    Measurement
	SelectedFields []string
	GroupByTags    []string
	WhereClause    string

	OrderClause string
	LimitClause string
	StepClause  string

	QueryCondition *datahub_v1alpha1.QueryCondition
}

func (s *StatementNew) AppendWhereCondition(key string, operator string, value string) {
	if value == "" {
		return
	}

	if s.WhereClause == "" {
		s.WhereClause += fmt.Sprintf("WHERE \"%s\"%s'%s' ", key, operator, value)
	} else {
		s.WhereClause += fmt.Sprintf("AND \"%s\"%s'%s' ", key, operator, value)
	}
}

func (s *StatementNew) AppendWhereConditionDirect(condition string) {
	if condition == "" {
		return
	}

	if s.WhereClause == "" {
		s.WhereClause += fmt.Sprintf("WHERE %s ", condition)
	} else {
		s.WhereClause += fmt.Sprintf("AND %s ", condition)
	}
}

func (s *StatementNew) AppendTimeCondition(operator string, value int64) {
	if value == 0 {
		return
	}

	tm := time.Unix(int64(value), 0)

	if s.WhereClause == "" {
		s.WhereClause += fmt.Sprintf("WHERE time%s'%s' ", operator, tm.UTC().Format(time.RFC3339))
	} else {
		s.WhereClause += fmt.Sprintf("AND time%s'%s' ", operator, tm.UTC().Format(time.RFC3339))
	}
}

func (s *StatementNew) AppendTimeConditionFromQueryCondition() {
	s.AppendTimeCondition(">=", s.QueryCondition.GetTimeRange().GetStartTime().GetSeconds())
	s.AppendTimeCondition("<=", s.QueryCondition.GetTimeRange().GetEndTime().GetSeconds())
}

func (s *StatementNew) AppendOrderClauseFromQueryCondition() {
	switch s.QueryCondition.GetOrder() {
	case datahub_v1alpha1.QueryCondition_ASC:
		s.OrderClause = "ORDER BY time ASC"
	case datahub_v1alpha1.QueryCondition_DESC:
		s.OrderClause = "ORDER BY time DESC"
	default:
		s.OrderClause = "ORDER BY time ASC"
	}
}

func (s *StatementNew) AppendLimitClauseFromQueryCondition() {
	limit := s.QueryCondition.GetLimit()
	if limit > 0 {
		s.LimitClause = fmt.Sprintf("LIMIT %v", limit)
	}
}

func (s StatementNew) BuildQueryCmd() string {
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
