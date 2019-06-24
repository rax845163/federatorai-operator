package updateenvvar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func TestUpdateEnvVarsToDeployment(t *testing.T) {

	type testCase struct {
		dep     *appsv1.Deployment
		envVars []corev1.EnvVar
	}

	testCases := []struct {
		have testCase
		want *appsv1.Deployment
	}{
		{
			have: testCase{
				dep: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									corev1.Container{
										Env: []corev1.EnvVar{
											corev1.EnvVar{
												Name:  "e1",
												Value: "v1",
											},
											corev1.EnvVar{
												Name:  "e2",
												Value: "v2",
											},
										},
									},
								},
							},
						},
					},
				},
				envVars: []corev1.EnvVar{
					corev1.EnvVar{
						Name:  "e1",
						Value: "new_v1",
					},
					corev1.EnvVar{
						Name:  "e3",
						Value: "v3",
					},
				},
			},
			want: &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								corev1.Container{
									Env: []corev1.EnvVar{
										corev1.EnvVar{
											Name:  "e1",
											Value: "new_v1",
										},
										corev1.EnvVar{
											Name:  "e2",
											Value: "v2",
										},
										corev1.EnvVar{
											Name:  "e3",
											Value: "v3",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	assert := assert.New(t)
	for _, testCase := range testCases {

		UpdateEnvVarsToDeployment(testCase.have.dep, testCase.have.envVars)
		assert.Equal(testCase.have.dep, testCase.want)
	}
}
