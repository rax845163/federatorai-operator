package prometheus

import (
	"github.com/containers-ai/alameda/datahub/pkg/dao"
	"time"
)

type Options struct {
	startTime             *time.Time
	endTime               *time.Time
	timeout               *time.Time
	stepTime              *time.Duration
	aggregateOverTimeFunc dao.AggregateFunction
}

type Option func(*Options)

func BuildDefaultOptions() Options {

	copyDefaultStepTime, _ := time.ParseDuration("30s")

	return Options{
		stepTime:              &copyDefaultStepTime,
		aggregateOverTimeFunc: dao.None,
	}
}

func StartTime(t *time.Time) Option {
	return func(o *Options) {
		o.startTime = t
	}
}

func EndTime(t *time.Time) Option {
	return func(o *Options) {
		o.endTime = t
	}
}

func Timeout(t *time.Time) Option {
	return func(o *Options) {
		o.endTime = t
	}
}

func StepDuration(d *time.Duration) Option {
	return func(o *Options) {
		o.stepTime = d
	}
}

func AggregateOverTimeFunction(f dao.AggregateFunction) Option {
	return func(o *Options) {
		o.aggregateOverTimeFunc = f
	}
}
