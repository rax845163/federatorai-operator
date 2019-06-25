package influxdb

import (
	"encoding/json"
	"errors"
	"fmt"
	Common "github.com/containers-ai/api/common"
	InfluxDBClient "github.com/influxdata/influxdb/client/v2"
	"strconv"
	"time"
)

func ReadRawdata(config *Config, queries []*Common.Query) ([]*Common.ReadRawdata, error) {
	influxClient := New(config)
	rawdata := make([]*Common.ReadRawdata, 0)

	for _, query := range queries {
		statement := NewInfluxStatement(query)
		statement.AppendTimeConditionIntoWhereClause()
		statement.SetLimitClauseFromQueryCondition()
		statement.SetOrderClauseFromQueryCondition()
		cmd := statement.BuildQueryCmd()

		results, err := influxClient.QueryDB(cmd, query.Database)
		if err != nil {
			scope.Errorf("failed to read rawdata from InfluxDB: %v", err)
			return make([]*Common.ReadRawdata, 0), err
		} else {
			readRawdata := InfluxResultToReadRawdata(results, query)
			rawdata = append(rawdata, readRawdata)
			// RepoInfluxDB.CompareRawdataWithInfluxResults(readRawdata, results) // For debug purpose
		}
	}

	return rawdata, nil
}

func WriteRawdata(config *Config, writeRawdata []*Common.WriteRawdata) error {
	influxClient := New(config)

	for _, rawdata := range writeRawdata {
		points := make([]*InfluxDBClient.Point, 0)

		for _, row := range rawdata.GetRows() {
			index := 0
			tags := make(map[string]string)
			fields := make(map[string]interface{})

			for _, value := range row.GetValues() {
				switch rawdata.GetColumnTypes()[index] {
				case Common.ColumnType_COLUMNTYPE_TAG:
					tags[rawdata.GetColumns()[index]] = value
				case Common.ColumnType_COLUMNTYPE_FIELD:
					fields[rawdata.GetColumns()[index]] = changeFormat(value, rawdata.GetDataTypes()[index])
				default:
					fmt.Println("not support")
				}
				index = index + 1
			}

			// Add time field depends on request
			if row.GetTime() == nil {
				pt, err := InfluxDBClient.NewPoint(rawdata.GetTable(), tags, fields, time.Unix(0, 0))
				if err == nil {
					points = append(points, pt)
				} else {
					fmt.Println(err.Error())
				}
			} else {
				pt, err := InfluxDBClient.NewPoint(rawdata.GetTable(), tags, fields, time.Unix(row.GetTime().GetSeconds(), 0))
				if err == nil {
					points = append(points, pt)
				} else {
					fmt.Println(err.Error())
				}
			}
		}

		err := influxClient.WritePoints(points, InfluxDBClient.BatchPointsConfig{Database: rawdata.GetDatabase()})
		if err != nil {
			scope.Error(err.Error())
		}

		return err
	}

	return nil
}

func InfluxResultToReadRawdata(results []InfluxDBClient.Result, query *Common.Query) *Common.ReadRawdata {
	readRawdata := Common.ReadRawdata{Query: query}

	if len(results[0].Series) == 0 {
		return &readRawdata
	}

	for _, result := range results {
		tagsLen := 0
		valuesLen := 0

		// Build columns
		for k := range results[0].Series[0].Tags {
			readRawdata.Columns = append(readRawdata.Columns, string(k))
			tagsLen = tagsLen + 1
		}
		for _, column := range results[0].Series[0].Columns {
			readRawdata.Columns = append(readRawdata.Columns, column)
			valuesLen = valuesLen + 1
		}

		// One series is one group
		for _, row := range result.Series {
			group := Common.Group{}

			// Build values
			for _, value := range row.Values {
				r := Common.Row{}

				// Tags
				for k, v := range readRawdata.Columns {
					r.Values = append(r.Values, row.Tags[v])
					if k >= (tagsLen - 1) {
						break
					}
				}

				// Fields
				for _, v := range value {
					switch v.(type) {
					case bool:
						r.Values = append(r.Values, strconv.FormatBool(v.(bool)))
					case string:
						r.Values = append(r.Values, v.(string))
					case json.Number:
						r.Values = append(r.Values, v.(json.Number).String())
					case nil:
						r.Values = append(r.Values, "")
					default:
						fmt.Println("Error, not support")
						r.Values = append(r.Values, v.(string))
					}
				}
				group.Rows = append(group.Rows, &r)
			}
			readRawdata.Groups = append(readRawdata.Groups, &group)
		}
	}

	return &readRawdata
}

func ReadRawdataToInfluxDBRow(readRawdata *Common.ReadRawdata) []*InfluxDBRow {
	influxDBRows := make([]*InfluxDBRow, 0)

	tagIndex := make([]int, 0)

	// locate tags index
	for _, tag := range readRawdata.GetQuery().GetCondition().GetGroups() {
		for index, column := range readRawdata.GetColumns() {
			if tag == column {
				tagIndex = append(tagIndex, index)
			}
		}
	}

	for _, group := range readRawdata.GetGroups() {
		influxDBRow := InfluxDBRow{
			Name: readRawdata.GetQuery().GetTable(),
			Tags: make(map[string]string),
		}

		for _, row := range group.GetRows() {
			// Pack tags
			for _, v := range tagIndex {
				for _, row := range group.GetRows() {
					influxDBRow.Tags[readRawdata.GetColumns()[v]] = row.GetValues()[v]
				}
			}

			// Pack data
			data := make(map[string]string)
			for index, column := range readRawdata.GetColumns() {
				data[column] = row.GetValues()[index]
			}
			influxDBRow.Data = append(influxDBRow.Data, data)
		}

		influxDBRows = append(influxDBRows, &influxDBRow)
	}

	return influxDBRows
}

func CompareRawdataWithInfluxResults(readRawdata *Common.ReadRawdata, results []InfluxDBClient.Result) error {
	before := PackMap(results)
	after := ReadRawdataToInfluxDBRow(readRawdata)
	message := ""

	for index, row := range after {
		compRow := before[index]

		// Check Name
		if row.Name != compRow.Name {
			message = message + fmt.Sprintf("Name: %s, %s\n", row.Name, compRow.Name)
			fmt.Printf("[ERROR] Name: %s, %s\n", row.Name, compRow.Name)
		}

		// Check Tags
		for key, value := range row.Tags {
			compValue := compRow.Tags[key]
			if compRow.Tags[key] != value {
				message = message + fmt.Sprintf("Tag[%s]: %s, %s\n", key, value, compValue)
				fmt.Printf("[ERROR] Tag[%s]: %s, %s\n", key, value, compValue)
			}
		}

		// Check Data
		for k, v := range row.Data {
			for key, value := range v {
				compValue := compRow.Data[k][key]
				if compValue != value {
					message = message + fmt.Sprintf("Data[%s]: %s, %s\n", key, value, compValue)
					fmt.Printf("[ERROR] Data[%s]: %s, %s\n", key, value, compValue)
				}
			}
		}
	}

	if message != "" {
		return errors.New(message)
	}

	return nil
}

func changeFormat(value string, dataType Common.DataType) interface{} {
	switch dataType {
	case Common.DataType_DATATYPE_BOOL:
		valueBool, _ := strconv.ParseBool(value)
		return valueBool
	case Common.DataType_DATATYPE_INT32:
		valueInt, _ := strconv.ParseInt(value, 10, 32)
		return valueInt
	case Common.DataType_DATATYPE_INT64:
		valueInt, _ := strconv.ParseInt(value, 10, 64)
		return valueInt
	case Common.DataType_DATATYPE_FLOAT32:
		valueFloat, _ := strconv.ParseFloat(value, 64)
		return valueFloat
	case Common.DataType_DATATYPE_STRING:
		return value
	default:
		fmt.Println("not support")
		return value
	}
}
