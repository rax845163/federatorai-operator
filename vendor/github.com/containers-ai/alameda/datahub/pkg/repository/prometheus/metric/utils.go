package metric

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/containers-ai/alameda/datahub/pkg/repository/prometheus"
	"github.com/pkg/errors"
)

type operateOptions struct {
	matchType matchType
	labels    []string
}

func newDefaultOperateOptions() operateOptions {
	return operateOptions{
		matchType: matchTypeNone,
		labels:    make([]string, 0),
	}
}

type matchType int

var (
	matchTypeNone   matchType = 0
	matchTypeIgnore matchType = 1
	matchTypeOn     matchType = 2
)

type operateOption func(*operateOptions)

func operateOptionsMatchType(t matchType) operateOption {
	return func(o *operateOptions) {
		o.matchType = t
	}
}

func operateOptionsLabels(labels []string) operateOption {
	return func(o *operateOptions) {
		o.labels = labels
	}
}

func oneToOneMultiply(entity1, entity2 prometheus.Entity, opts ...operateOption) (prometheus.Entity, error) {

	operateOptions := newDefaultOperateOptions()
	for _, opt := range opts {
		opt(&operateOptions)
	}

	entity1.Labels = filterLabelsByOperateOptions(operateOptions, entity1.Labels)
	entity2.Labels = filterLabelsByOperateOptions(operateOptions, entity2.Labels)

	entity := prometheus.Entity{}

	if len(entity1.Values) != len(entity2.Values) {
		return entity, errors.New("prometheus one to one vector multiply failed: length of two vectors' values are not equal")
	}

	if !reflect.DeepEqual(entity1.Labels, entity2.Labels) {
		return entity, errors.Errorf("prometheus one to one vector multiply failed: labels map of two vectors are not equal, entity1's labels: %+v, entity2's labels: %+v", entity1.Labels, entity2.Labels)
	}

	entity = prometheus.Entity{
		Labels: entity1.Labels,
		Values: make([]prometheus.UnixTimeWithSampleValue, len(entity1.Values)),
	}

	for i, vector1Value := range entity1.Values {

		vector2Value := entity2.Values[i]

		if vector1Value.UnixTime != vector2Value.UnixTime {
			return entity, errors.Errorf("prometheus one to one vector multiply failed: value's timestamp not equal")
		}

		value1, err := strconv.ParseFloat(vector1Value.SampleValue, 64)
		if err != nil {
			return entity, errors.Errorf("prometheus one to one vector multiply failed: parse value failed: %s", err.Error())
		}
		value2, err := strconv.ParseFloat(vector2Value.SampleValue, 64)
		if err != nil {
			return entity, errors.Errorf("prometheus one to one vector multiply failed: parse value failed: %s", err.Error())
		}
		result := value1 * value2
		entity.Values[i] = prometheus.UnixTimeWithSampleValue{
			UnixTime:    vector1Value.UnixTime,
			SampleValue: fmt.Sprintf("%f", result),
		}
	}

	return entity, nil
}

func filterLabelsByOperateOptions(opt operateOptions, labels map[string]string) map[string]string {

	switch opt.matchType {
	case matchTypeIgnore:
		for _, keyToIgnore := range opt.labels {
			if _, exist := labels[keyToIgnore]; exist {
				delete(labels, keyToIgnore)
			}
		}
	case matchTypeOn:
		preserveKeysMap := make(map[string]bool)
		for _, keyToPreserve := range opt.labels {
			preserveKeysMap[keyToPreserve] = true
		}
		for key := range labels {
			if _, exist := preserveKeysMap[key]; !exist {
				delete(labels, key)
			}
		}
	}
	return labels
}
