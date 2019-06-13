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

package node

import (
	"context"

	datahub_node "github.com/containers-ai/alameda/operator/datahub/client/node"
	logUtil "github.com/containers-ai/alameda/pkg/utils/log"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
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
	scope = logUtil.RegisterScope("node_controller", "node controller log", 0)
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Node Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileNode{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("node-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Node
	err = c.Watch(&source.Kind{Type: &corev1.Node{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileNode{}

// ReconcileNode reconciles a Node object
type ReconcileNode struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Node object and makes changes based on the state read
// and what is in the Node.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// +kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=nodes/status,verbs=get;update;patch
func (r *ReconcileNode) Reconcile(request reconcile.Request) (reconcile.Result, error) {

	instance := &corev1.Node{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)

	nodeIsDeleted := false
	if err != nil && k8sErrors.IsNotFound(err) {
		nodeIsDeleted = true
		instance.Namespace = request.Namespace
		instance.Name = request.Name
	} else if err != nil {
		scope.Error(err.Error())
	}

	if err := syncNodeDependentResource(nodeIsDeleted, instance); err != nil {
		scope.Error(err.Error())
	}

	return reconcile.Result{}, nil
}

func syncNodeDependentResource(isDeleted bool, node *corev1.Node) error {

	nodes := make([]*corev1.Node, 0)
	nodes = append(nodes, node)

	if !isDeleted {
		if err := createNodeDependentResource(nodes); err != nil {
			return errors.Wrapf(err, "sync node dependent resource failed: %s", err.Error())
		}
	} else {
		if err := deleteNodeDependentResource(nodes); err != nil {
			return errors.Wrapf(err, "sync node dependent resource failed: %s", err.Error())
		}
	}

	return nil
}

func createNodeDependentResource(nodes []*corev1.Node) error {

	datahubNodeRepo := datahub_node.NewAlamedaNodeRepository()
	if err := datahubNodeRepo.CreateAlamedaNode(nodes); err != nil {
		return errors.Wrapf(err, "create node dependent resource failed: %s", err.Error())
	}

	return nil
}

func deleteNodeDependentResource(nodes []*corev1.Node) error {

	datahubNodeRepo := datahub_node.NewAlamedaNodeRepository()
	if err := datahubNodeRepo.DeleteAlamedaNodes(nodes); err != nil {
		return errors.Wrapf(err, "delete node dependent resource failed: %s", err.Error())
	}

	return nil
}
