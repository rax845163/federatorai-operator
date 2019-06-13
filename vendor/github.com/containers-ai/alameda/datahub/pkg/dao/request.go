package dao

import (
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

type AggregateFunction = int

const (
	None AggregateFunction = 0
	Max  AggregateFunction = 1
)

// QueryCondition Others query condition
type QueryCondition struct {
	StartTime                 *time.Time
	EndTime                   *time.Time
	StepTime                  *time.Duration
	TimestampOrder            Order
	Limit                     int
	AggregateOverTimeFunction AggregateFunction
}
