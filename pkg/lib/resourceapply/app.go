package resourceapply

import (
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsclientv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

func ApplyDeployment(client appsclientv1.DeploymentsGetter, required *appsv1.Deployment) (*appsv1.Deployment, bool, error) {
	_, err := client.Deployments(required.Namespace).Get(required.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		logrus.Info("Not Found Deployments And Create %v Deployments", required.Name)
		actual, err := client.Deployments(required.Namespace).Create(required)
		logrus.Info("Create Deployments ", required.Name)
		return actual, true, err
	}
	if err != nil {
		return nil, false, err
	}
	logrus.Info("Found Deployment ", required.Name)
	return nil, false, err
}
