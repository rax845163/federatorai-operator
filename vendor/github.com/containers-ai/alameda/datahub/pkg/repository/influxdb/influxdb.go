package influxdb

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/containers-ai/alameda/pkg/utils/log"
	client "github.com/influxdata/influxdb/client/v2"
)

type Database string
type Measurement string

var (
	// ZeroTime is used as a constant of timestamp
	ZeroTime time.Time = time.Unix(0, 0)
)

type InfluxDBEntity struct {
	Time   time.Time
	Tags   map[string]string
	Fields map[string]interface{}
}

type InfluxDBRow struct {
	Name    string
	Tags    map[string]string
	Data    []map[string]string
	Partial bool
}

var (
	scope = log.RegisterScope("influxdb_client", "influxdb client", 0)
)

// InfluxDBRepository interacts with database
type InfluxDBRepository struct {
	Address                string
	Username               string
	Password               string
	RetentionDuration      string
	RetentionShardDuration string
}

// New returns InfluxDBRepository instance
func New(influxCfg *Config) *InfluxDBRepository {
	return &InfluxDBRepository{
		Address:                influxCfg.Address,
		Username:               influxCfg.Username,
		Password:               influxCfg.Password,
		RetentionDuration:      influxCfg.RetentionDuration,
		RetentionShardDuration: influxCfg.RetentionShardDuration,
	}
}

// CreateDatabase creates database
func (influxDBRepository *InfluxDBRepository) CreateDatabase(db string) error {
	_, err := influxDBRepository.QueryDB(fmt.Sprintf("CREATE DATABASE %s", db), db)
	return err
}

// Modify default retention policy
func (influxDBRepository *InfluxDBRepository) ModifyDefaultRetentionPolicy(db string) error {
	duration := influxDBRepository.RetentionDuration
	shardGroupDuration := influxDBRepository.RetentionShardDuration
	retentionCmd := fmt.Sprintf("ALTER RETENTION POLICY \"autogen\" on \"%s\" DURATION %s SHARD DURATION %s", db, duration, shardGroupDuration)
	_, err := influxDBRepository.QueryDB(retentionCmd, db)
	return err
}

func (influxDBRepository *InfluxDBRepository) newHttpClient() client.Client {
	clnt, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:               influxDBRepository.Address,
		Username:           influxDBRepository.Username,
		Password:           influxDBRepository.Password,
		InsecureSkipVerify: true,
	})
	if err != nil {
		scope.Error(err.Error())
	}
	return clnt
}

// WritePoints writes points to database
func (influxDBRepository *InfluxDBRepository) WritePoints(points []*client.Point, bpCfg client.BatchPointsConfig) error {
	clnt := influxDBRepository.newHttpClient()
	defer clnt.Close()

	bp, err := client.NewBatchPoints(bpCfg)
	if err != nil {
		scope.Error(err.Error())
	}

	for _, point := range points {
		bp.AddPoint(point)
	}

	if err := clnt.Write(bp); err != nil {
		if strings.Contains(err.Error(), "database not found") {
			if err = influxDBRepository.CreateDatabase(bpCfg.Database); err != nil {
				scope.Error(err.Error())
				return err
			} else {
				err = influxDBRepository.WritePoints(points, bpCfg)
			}
		}
		if err != nil {
			scope.Error(err.Error())
			fmt.Print(err.Error())
			return err
		}
	}

	return nil
}

// QueryDB queries database
func (influxDBRepository *InfluxDBRepository) QueryDB(cmd, database string) (res []client.Result, err error) {
	clnt := influxDBRepository.newHttpClient()
	defer clnt.Close()
	q := client.Query{
		Command:  cmd,
		Database: database,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func PackMap(results []client.Result) []*InfluxDBRow {
	var rows []*InfluxDBRow

	if len(results[0].Series) == 0 {
		return rows
	}

	for _, result := range results {
		for _, r := range result.Series {
			row := InfluxDBRow{Name: r.Name, Partial: r.Partial}
			row.Tags = r.Tags
			for _, v := range r.Values {
				data := make(map[string]string)
				// Pack tag
				for key, value := range r.Tags {
					data[key] = value
				}
				// Pack values
				for j, col := range r.Columns {
					value := v[j]
					if value != nil {
						switch value.(type) {
						case bool:
							data[col] = strconv.FormatBool(value.(bool))
						case string:
							data[col] = value.(string)
						case json.Number:
							data[col] = value.(json.Number).String()
						case nil:
							data[col] = ""
						default:
							data[col] = value.(string)
						}
					} else {
						data[col] = ""
					}
				}
				row.Data = append(row.Data, data)
			}
			rows = append(rows, &row)
		}
	}

	return NormalizeResult(rows)
}

func NormalizeResult(rows []*InfluxDBRow) []*InfluxDBRow {
	var rowList []*InfluxDBRow

	for _, r := range rows {
		row := InfluxDBRow{Name: r.Name, Partial: r.Partial}
		row.Tags = r.Tags
		for _, d := range r.Data {
			data := make(map[string]string)
			for key, value := range d {
				if strings.HasSuffix(key, "_1") {
					found := false
					newKey := strings.TrimSuffix(key, "_1")
					for k := range data {
						if k == newKey {
							found = true
							if value != "" {
								data[k] = value
							}
							break
						}
					}
					if !found {
						data[key] = value
					}
				} else {
					found := false
					newKey := fmt.Sprintf("%s_1", key)
					for k, v := range data {
						if k == newKey {
							found = true
							if v != "" {
								delete(data, newKey)
								data[key] = v
							} else {
								delete(data, newKey)
								data[key] = value
							}
							break
						}
					}
					if !found {
						data[key] = value
					}
				}
			}
			row.Data = append(row.Data, data)
		}
		rowList = append(rowList, &row)
	}

	return rowList
}
