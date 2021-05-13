/*
Copyright 2021.

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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mikutasv1alpha1 "github.com/mikutas/job-deletor/api/v1alpha1"
)

// JobDeletorReconciler reconciles a JobDeletor object
type JobDeletorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=mikutas.example.com,resources=jobdeletors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=mikutas.example.com,resources=jobdeletors/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=mikutas.example.com,resources=jobdeletors/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the JobDeletor object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *JobDeletorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("jobdeletor", req.NamespacedName)

	// your logic here
	var jd mikutasv1alpha1.JobDeletor
	log.Info("fetching DeploymentMaxAge Resource")
	if err := r.Get(ctx, req.NamespacedName, &jd); err != nil {
		log.Error(err, "unable to fetch JobDeletor")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("jd.ClusterName=" + jd.ClusterName)
	log.Info("jd.Namespace=" + jd.Namespace)
	log.Info("jd.Name=" + jd.Name)
	log.Info("jd.Spec.DeletionTargetStatus=" + jd.Spec.DeletionTargetStatus)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *JobDeletorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mikutasv1alpha1.JobDeletor{}).
		Complete(r)
}
