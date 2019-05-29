package resourceapply

import (
	"fmt"
	"time"

	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/alamedaserviceparamter"
	"github.com/containers-ai/federatorai-operator/pkg/util"
	"github.com/pkg/errors"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextclientv1beta1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("controller_alamedaservice")

func ApplyCustomResourceDefinition(client apiextclientv1beta1.CustomResourceDefinitionsGetter, required *apiextv1beta1.CustomResourceDefinition, asp *alamedaserviceparamter.AlamedaServiceParamter) (*apiextv1beta1.CustomResourceDefinition, error) {

	waitInterval := 500 * time.Millisecond

	cluster, err := client.CustomResourceDefinitions().Get(required.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		log.Info("Not Found CRD And Create", "CustomResourceDefinition.Name", required.Name)
		actual, err := client.CustomResourceDefinitions().Create(required)
		if err != nil {
			return nil, errors.Wrapf(err, "apply CustomResourceDefinition %s failed", required.Name)
		}
		if err = wait.Poll(waitInterval, 30*time.Second, func() (bool, error) {
			_, getErr := client.CustomResourceDefinitions().Get(required.Name, metav1.GetOptions{})
			if getErr == nil {
				log.Info("Get CRD Ok", "CRD.Name", required.Name)
				return true, nil
			} else {
				log.V(-1).Info(fmt.Sprintf("get CustomResourceDefinition failed, will retry after %f seconds", waitInterval.Seconds()), "CustomResourceDefinition.Name", required.Name)
				return false, nil
			}
		}); err != nil {
			return nil, errors.Wrapf(err, "apply CustomResourceDefinition %s failed: waiting CustomResourceDefinition timeout", required.Name)
		}
		return actual, err
	} else {
		log.Info("Found CRD", "CustomResourceDefinition.Name", required.Name)
		if asp.CheckCurrentCRDIsChangeVersion() && cluster.Name == util.AlamedaScalerName { //if user change scaler CRD version (component version change)
			cluster.Spec = required.Spec //replace crd spec
			actual, err := client.CustomResourceDefinitions().Update(cluster)
			if err != nil {
				return nil, errors.Wrapf(err, "update CustomResourceDefinition %s failed", required.Name)
			}
			asp.SetCurrentCRDChangeVersionToFalse()
			log.Info("change CRD Version", "CustomResourceDefinition.Name", required.Name)
			return actual, err
		}
	}
	if err != nil {
		return nil, errors.Wrapf(err, "apply CustomResourceDefinition %s failed", required.Name)
	}
	return nil, nil
}
func DeleteCustomResourceDefinition(client apiextclientv1beta1.CustomResourceDefinitionsGetter, required *apiextv1beta1.CustomResourceDefinition) (*apiextv1beta1.CustomResourceDefinition, bool, error) {
	existing, err := client.CustomResourceDefinitions().Get(required.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		log.Info("Not Found CRD", "CustomResourceDefinition.Name", required.Name)
		return nil, false, err
	}
	if err != nil {
		return nil, false, err
	} else {
		log.Info("Found CRD And Delete", "CustomResourceDefinition.Name", required.Name)
		err = client.CustomResourceDefinitions().Delete(existing.Name, nil)
	}
	return nil, false, err
}
