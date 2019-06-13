/*
Copyright 2019 The Alameda Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package alamedarecommendation

import (
	"context"
	"fmt"
	"strings"
	"time"

	autoscalingv1alpha1 "github.com/containers-ai/alameda/operator/pkg/apis/autoscaling/v1alpha1"
	alamedarecommendation_reconciler "github.com/containers-ai/alameda/operator/pkg/reconciler/alamedarecommendation"
	alamedascaler_reconciler "github.com/containers-ai/alameda/operator/pkg/reconciler/alamedascaler"
	datahubutils "github.com/containers-ai/alameda/operator/pkg/utils/datahub"
	utilsresource "github.com/containers-ai/alameda/operator/pkg/utils/resources"
	logUtil "github.com/containers-ai/alameda/pkg/utils/log"

	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	"google.golang.org/grpc"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	alamedarecommendationScope = logUtil.RegisterScope("alamedarecommendation", "alameda recommendation", 0)
)

var cachedFirstSynced = false

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new AlamedaRecommendation Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileAlamedaRecommendation{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("alamedarecommendation-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to AlamedaRecommendation
	err = c.Watch(&source.Kind{Type: &autoscalingv1alpha1.AlamedaRecommendation{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileAlamedaRecommendation{}

// ReconcileAlamedaRecommendation reconciles a AlamedaRecommendation object
type ReconcileAlamedaRecommendation struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a AlamedaRecommendation object and makes changes based on the state read
// and what is in the AlamedaRecommendation.Spec
func (r *ReconcileAlamedaRecommendation) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	if !cachedFirstSynced {
		time.Sleep(5 * time.Second)
	}
	cachedFirstSynced = true

	getResource := utilsresource.NewGetResource(r)
	listResources := utilsresource.NewListResources(r)

	if alamedaRecommendation, err := getResource.GetAlamedaRecommendation(request.Namespace, request.Name); err == nil {
		// Remove AlamedaResource if target is not existed
		for _, or := range alamedaRecommendation.OwnerReferences {
			if or.Controller != nil && *or.Controller && strings.ToLower(or.Kind) == "alamedascaler" {
				if scalers, err := listResources.ListAllAlamedaScaler(); err == nil {
					for _, scaler := range scalers {
						if scaler.GetUID() == or.UID {
							alamedascalerReconciler := alamedascaler_reconciler.NewReconciler(r, &scaler)
							if !alamedascalerReconciler.HasAlamedaPod(alamedaRecommendation.Namespace, alamedaRecommendation.Name) {
								alamedarecommendationScope.Infof(fmt.Sprintf("AlamedaRecommendation (%s/%s) is already removed from AlamedaScaler (%s/%s)", request.Namespace, request.Name, scaler.Namespace, scaler.Name))
								if err = r.Delete(context.TODO(), alamedaRecommendation); err != nil {
									alamedarecommendationScope.Error(err.Error())
								}
								return reconcile.Result{}, nil
							}
						}
					}
				} else {
					alamedarecommendationScope.Error(err.Error())
				}
			}
		}

		// Update Recommendation from datahub
		alamedarecommendationReconciler := alamedarecommendation_reconciler.NewReconciler(r, alamedaRecommendation)
		if conn, err := grpc.Dial(datahubutils.GetDatahubAddress(), grpc.WithInsecure()); err == nil {
			defer conn.Close()
			aiServiceClnt := datahub_v1alpha1.NewDatahubServiceClient(conn)
			req := datahub_v1alpha1.ListPodRecommendationsRequest{
				NamespacedName: &datahub_v1alpha1.NamespacedName{
					Namespace: alamedaRecommendation.GetNamespace(),
					Name:      alamedaRecommendation.GetName(),
				},
			}
			if podRecommendationsRes, err := aiServiceClnt.ListPodRecommendations(context.Background(), &req); err == nil && len(podRecommendationsRes.GetPodRecommendations()) == 1 {
				if alamedaRecommendation, err = alamedarecommendationReconciler.UpdateResourceRecommendation(podRecommendationsRes.GetPodRecommendations()[0]); err == nil {
					if err = r.Update(context.TODO(), alamedaRecommendation); err != nil {
						alamedarecommendationScope.Error(err.Error())
					}
				}
			}
		} else {
			alamedarecommendationScope.Error(err.Error())
		}
	} else if !k8sErrors.IsNotFound(err) {
		alamedarecommendationScope.Errorf("get AlamedaRecommendation %s/%s failed: %s", request.Namespace, request.Name, err.Error())
	}

	return reconcile.Result{}, nil
}
