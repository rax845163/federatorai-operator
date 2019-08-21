package alamedaserviceparamter

import (
	"testing"

	"github.com/containers-ai/federatorai-operator/pkg/apis/federatorai/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAdd(t *testing.T) {

	type testCaseHave struct {
		origin *Resource
		toAdd  Resource
	}

	type testCase struct {
		have testCaseHave
		want *Resource
	}

	testCases := []testCase{
		testCase{
			have: testCaseHave{
				origin: &Resource{
					ServiceList:   []string{},
					ConfigMapList: []string{"c1"},
				},
				toAdd: Resource{
					ServiceList:   []string{"s1"},
					ConfigMapList: []string{"c2", "c3"},
				},
			},
			want: &Resource{
				ServiceList:   []string{"s1"},
				ConfigMapList: []string{"c1", "c2", "c3"},
			},
		},
		testCase{
			have: testCaseHave{
				origin: &Resource{},
				toAdd: Resource{
					ServiceList:   []string{"s1", "s2", "s3"},
					ConfigMapList: []string{"c1", "c2", "c3"},
				},
			},
			want: &Resource{
				ServiceList:   []string{"s1", "s2", "s3"},
				ConfigMapList: []string{"c1", "c2", "c3"},
			},
		},
	}

	assert := assert.New(t)
	for _, testCase := range testCases {
		actual := testCase.have.origin
		actual.add(testCase.have.toAdd)
		assert.Equal(testCase.want, actual)
	}

}
func TestDelete(t *testing.T) {

	type testCaseHave struct {
		origin   *Resource
		toDelete Resource
	}

	type testCase struct {
		have testCaseHave
		want *Resource
	}

	testCases := []testCase{
		testCase{
			have: testCaseHave{
				origin: &Resource{
					ServiceList:    []string{"s1", "s2", "s3"},
					ConfigMapList:  []string{"c1", "c2", "c3"},
					DeploymentList: []string{"d1", "d2", "d3"},
				},
				toDelete: Resource{
					ServiceList:    []string{"s3"},
					ConfigMapList:  []string{"c1", "c3"},
					DeploymentList: []string{},
				},
			},
			want: &Resource{
				ServiceList:    []string{"s1", "s2"},
				ConfigMapList:  []string{"c2"},
				DeploymentList: []string{"d1", "d2", "d3"},
			},
		},
		testCase{
			have: testCaseHave{
				origin: &Resource{
					ServiceList:    []string{"s1", "s2", "s3"},
					ConfigMapList:  []string{"c1", "c2", "c3"},
					DeploymentList: []string{"d1", "d2", "d3"},
				},
				toDelete: Resource{
					ServiceList:    []string{"s4"},
					ConfigMapList:  []string{"c1", "c3"},
					DeploymentList: []string{},
				},
			},
			want: &Resource{
				ServiceList:    []string{"s1", "s2", "s3"},
				ConfigMapList:  []string{"c2"},
				DeploymentList: []string{"d1", "d2", "d3"},
			},
		},
	}

	assert := assert.New(t)
	for i, testCase := range testCases {
		actual := testCase.have.origin
		actual.delete(testCase.have.toDelete)
		assert.EqualValuesf(testCase.want, actual, "test case: #%d", i)
	}

}

func TestGetInstallResource(t *testing.T) {

	type testCase struct {
		have v1alpha1.AlamedaService
		want Resource
	}

	var (
		ns   = "test"
		name = "test"
	)

	var defaultResource Resource
	for _, defaultInstallList := range defaultInstallLists {
		resource, _ := getResourceFromList(defaultInstallList)
		defaultResource.add(resource)
	}
	t0 := testCase{
		have: v1alpha1.AlamedaService{
			ObjectMeta: metav1.ObjectMeta{Namespace: ns, Name: name},
			Spec:       v1alpha1.AlamedaServiceSpec{},
		},
		want: defaultResource,
	}

	defaultAlamedaScalerCRD, _ := getResourceFromList(alamedaScalerCRD)
	v2AlamedaScalerCRD, _ := getResourceFromList(alamedaScalerCRDV2)
	t1Want := defaultResource
	t1Want.add(v2AlamedaScalerCRD)
	t1Want.delete(defaultAlamedaScalerCRD)
	t1 := testCase{
		have: v1alpha1.AlamedaService{
			ObjectMeta: metav1.ObjectMeta{Namespace: ns, Name: name},
			Spec:       v1alpha1.AlamedaServiceSpec{Version: "latest"},
		},
		want: t1Want,
	}

	dispatcherResource := GetDispatcherResource()
	weavescopeResource, _ := getResourceFromList(weavescopeList)
	t2Want := defaultResource
	t2Want.add(*dispatcherResource)
	t2Want.add(weavescopeResource)
	t2 := testCase{
		have: v1alpha1.AlamedaService{
			ObjectMeta: metav1.ObjectMeta{Namespace: ns, Name: name},
			Spec:       v1alpha1.AlamedaServiceSpec{EnableDispatcher: true, EnableWeavescope: true},
		},
		want: t2Want,
	}

	testCases := []testCase{t0, t1, t2}

	assert := assert.New(t)
	for i, testCase := range testCases {
		asp := NewAlamedaServiceParamter(&testCase.have)
		actual := asp.GetInstallResource()

		assert.EqualValuesf(testCase.want, *actual, "test_case: index #%d ", i)
	}

}
