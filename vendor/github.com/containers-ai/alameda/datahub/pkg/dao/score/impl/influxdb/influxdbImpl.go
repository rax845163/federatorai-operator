package influxdb

import (
	score_dao "github.com/containers-ai/alameda/datahub/pkg/dao/score"
	influxdb_entity_score "github.com/containers-ai/alameda/datahub/pkg/entity/influxdb/score"
	"github.com/containers-ai/alameda/datahub/pkg/repository/influxdb"
	influxdb_repository_score "github.com/containers-ai/alameda/datahub/pkg/repository/influxdb/score"
	"github.com/pkg/errors"
)

type influxdbDAO struct {
	config influxdb.Config
}

// NewWithConfig New influxdb score dao implement
func NewWithConfig(config influxdb.Config) score_dao.DAO {
	return influxdbDAO{
		config: config,
	}
}

// ListSimulatedScheduingScores Function implementation of score dao
func (dao influxdbDAO) ListSimulatedScheduingScores(request score_dao.ListRequest) ([]*score_dao.SimulatedSchedulingScore, error) {

	var (
		err error

		scoreRepository       influxdb_repository_score.SimulatedSchedulingScoreRepository
		influxdbScoreEntities []*influxdb_entity_score.SimulatedSchedulingScoreEntity
		scores                = make([]*score_dao.SimulatedSchedulingScore, 0)
	)

	scoreRepository = influxdb_repository_score.NewRepositoryWithConfig(dao.config)
	influxdbScoreEntities, err = scoreRepository.ListScoresByRequest(request)
	if err != nil {
		return scores, errors.Wrap(err, "list simulated scheduing scores failed")
	}

	for _, influxdbScoreEntity := range influxdbScoreEntities {

		score := score_dao.SimulatedSchedulingScore{
			Timestamp: influxdbScoreEntity.Time,
		}

		if scoreBefore := influxdbScoreEntity.ScoreBefore; scoreBefore != nil {
			score.ScoreBefore = *scoreBefore
		}

		if scoreAfter := influxdbScoreEntity.ScoreAfter; scoreAfter != nil {
			score.ScoreAfter = *scoreAfter
		}

		scores = append(scores, &score)
	}

	return scores, nil
}

// CreateSimulatedScheduingScores Function implementation of score dao
func (dao influxdbDAO) CreateSimulatedScheduingScores(scores []*score_dao.SimulatedSchedulingScore) error {

	var (
		err error

		scoreRepository influxdb_repository_score.SimulatedSchedulingScoreRepository
	)

	scoreRepository = influxdb_repository_score.NewRepositoryWithConfig(dao.config)
	err = scoreRepository.CreateScores(scores)
	if err != nil {
		return errors.Wrap(err, "create simulated scheduing scores failed")
	}

	return nil
}
