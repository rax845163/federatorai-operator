package enumconv

import (
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
)

const (
	Stable  string = "Stable"
	Compact string = "Compact"
)

var RecommendationPolicyEnum map[string]datahub_v1alpha1.RecommendationPolicy = map[string]datahub_v1alpha1.RecommendationPolicy{
	Stable:  datahub_v1alpha1.RecommendationPolicy_STABLE,
	Compact: datahub_v1alpha1.RecommendationPolicy_COMPACT,
}

var RecommendationPolicyDisp map[datahub_v1alpha1.RecommendationPolicy]string = map[datahub_v1alpha1.RecommendationPolicy]string{
	datahub_v1alpha1.RecommendationPolicy_STABLE:  Stable,
	datahub_v1alpha1.RecommendationPolicy_COMPACT: Compact,
}
