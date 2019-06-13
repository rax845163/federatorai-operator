package score

import (
	score_dao "github.com/containers-ai/alameda/datahub/pkg/dao/score"
	influxdb_entity_score "github.com/containers-ai/alameda/datahub/pkg/entity/influxdb/score"
	"github.com/containers-ai/alameda/datahub/pkg/repository/influxdb"
	influxdb_client "github.com/influxdata/influxdb/client/v2"
	"github.com/pkg/errors"
)

// SimulatedSchedulingScoreRepository Repository of simulated_scheduling_score data
type SimulatedSchedulingScoreRepository struct {
	influxDB *influxdb.InfluxDBRepository
}

// NewRepositoryWithConfig New SimulatedSchedulingScoreRepository with influxdb configuration
func NewRepositoryWithConfig(cfg influxdb.Config) SimulatedSchedulingScoreRepository {
	return SimulatedSchedulingScoreRepository{
		influxDB: influxdb.New(&cfg),
	}
}

// ListScoresByRequest List scores from influxDB
func (r SimulatedSchedulingScoreRepository) ListScoresByRequest(request score_dao.ListRequest) ([]*influxdb_entity_score.SimulatedSchedulingScoreEntity, error) {

	var (
		err error

		results      []influxdb_client.Result
		influxdbRows []*influxdb.InfluxDBRow
		scores       = make([]*influxdb_entity_score.SimulatedSchedulingScoreEntity, 0)
	)

	influxdbStatement := influxdb.Statement{
		Measurement: SimulatedSchedulingScore,
	}

	queryCondition := influxdb.QueryCondition{
		StartTime:      request.QueryCondition.StartTime,
		EndTime:        request.QueryCondition.EndTime,
		StepTime:       request.QueryCondition.StepTime,
		TimestampOrder: request.QueryCondition.TimestampOrder,
		Limit:          request.QueryCondition.Limit,
	}
	influxdbStatement.AppendTimeConditionIntoWhereClause(queryCondition)
	influxdbStatement.SetLimitClauseFromQueryCondition(queryCondition)
	influxdbStatement.SetOrderClauseFromQueryCondition(queryCondition)
	cmd := influxdbStatement.BuildQueryCmd()

	results, err = r.influxDB.QueryDB(cmd, string(influxdb.Score))
	if err != nil {
		return scores, errors.Wrap(err, "list scores failed")
	}

	influxdbRows = influxdb.PackMap(results)
	for _, influxdbRow := range influxdbRows {
		for _, data := range influxdbRow.Data {
			scoreEntity := influxdb_entity_score.NewSimulatedSchedulingScoreEntityFromMap(data)
			scores = append(scores, &scoreEntity)
		}
	}

	return scores, nil

}

// CreateScores Create simulated_scheduling_score data points into influxdb
func (r SimulatedSchedulingScoreRepository) CreateScores(scores []*score_dao.SimulatedSchedulingScore) error {

	var (
		err error

		points = make([]*influxdb_client.Point, 0)
	)

	for _, score := range scores {

		time := score.Timestamp
		scoreBefore := score.ScoreBefore
		scoreAfter := score.ScoreAfter
		entity := influxdb_entity_score.SimulatedSchedulingScoreEntity{
			Time:        time,
			ScoreBefore: &scoreBefore,
			ScoreAfter:  &scoreAfter,
		}

		point, err := entity.InfluxDBPoint(string(SimulatedSchedulingScore))
		if err != nil {
			return errors.Wrap(err, "create scores failed")
		}
		points = append(points, point)
	}

	err = r.influxDB.WritePoints(points, influxdb_client.BatchPointsConfig{
		Database: string(influxdb.Score),
	})
	if err != nil {
		return errors.Wrap(err, "create scores failed")
	}

	return nil
}
