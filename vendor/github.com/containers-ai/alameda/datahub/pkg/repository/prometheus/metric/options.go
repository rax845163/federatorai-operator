package metric

import (
	"time"

	"github.com/containers-ai/alameda/datahub/pkg/dao"
)

type Options struct {
	startTime             *time.Time
	endTime               *time.Time
	stepTime              *time.Duration
	aggregateOverTimeFunc dao.AggregateFunction
}

func buildDefaultOptions() Options {

	copyDefaultStepTime := defaultStepTime

	return Options{
		stepTime:              &copyDefaultStepTime,
		aggregateOverTimeFunc: dao.None,
	}
}

type Option func(*Options)

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
