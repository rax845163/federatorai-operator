package score

import (
	"strconv"
	"time"

	"github.com/containers-ai/alameda/datahub/pkg/utils"
	influxdb_client "github.com/influxdata/influxdb/client/v2"
)

type simulatedSchedulingScoreField = string
type simulatedSchedulingScoreTag = string

const (
	// SimulatedSchedulingScoreTime is the time simulatedSchedulingScore information is inserted to databse
	SimulatedSchedulingScoreTime simulatedSchedulingScoreTag = "time"

	// SimulatedSchedulingScoreScoreBefore Represents the field name in influxdb
	SimulatedSchedulingScoreScoreBefore simulatedSchedulingScoreField = "score_before"
	// SimulatedSchedulingScoreScoreAfter Represents the field name in influxdb
	SimulatedSchedulingScoreScoreAfter simulatedSchedulingScoreField = "score_after"
)

// SimulatedSchedulingScoreEntity Represents a record in influxdb
type SimulatedSchedulingScoreEntity struct {
	Time        time.Time
	ScoreBefore *float64
	ScoreAfter  *float64
}

// NewSimulatedSchedulingScoreEntityFromMap Build entity from map
func NewSimulatedSchedulingScoreEntityFromMap(data map[string]string) SimulatedSchedulingScoreEntity {

	tempTimestamp, _ := utils.ParseTime(data[SimulatedSchedulingScoreTime])

	entity := SimulatedSchedulingScoreEntity{
		Time: tempTimestamp,
	}

	if scoreBefore, exist := data[SimulatedSchedulingScoreScoreBefore]; exist {
		value, _ := strconv.ParseFloat(scoreBefore, 64)
		entity.ScoreBefore = &value
	}

	if scoreAfter, exist := data[SimulatedSchedulingScoreScoreAfter]; exist {
		value, _ := strconv.ParseFloat(scoreAfter, 64)
		entity.ScoreAfter = &value
	}

	return entity
}

// InfluxDBPoint Build influxdb point base on current entity's properties
func (e SimulatedSchedulingScoreEntity) InfluxDBPoint(measurementName string) (*influxdb_client.Point, error) {

	tags := map[string]string{}

	fields := map[string]interface{}{}
	if e.ScoreBefore != nil {
		fields[SimulatedSchedulingScoreScoreBefore] = *e.ScoreBefore
	}
	if e.ScoreAfter != nil {
		fields[SimulatedSchedulingScoreScoreAfter] = *e.ScoreAfter
	}

	return influxdb_client.NewPoint(measurementName, tags, fields, e.Time)
}
