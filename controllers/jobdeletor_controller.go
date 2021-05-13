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
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mikutasv1alpha1 "github.com/mikutas/job-deletor/api/v1alpha1"
)

const (
	targetStatusSucceeded = "succeeded"
	targetStatusFailed    = "failed"
	targetStatusAll       = "all"
)

// JobDeletorReconciler reconciles a JobDeletor object
type JobDeletorReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=mikutas.example.com,resources=jobdeletors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=mikutas.example.com,resources=jobdeletors/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=mikutas.example.com,resources=jobdeletors/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=namespaces,verbs=list;watch
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;delete
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch

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
	log.Info("fetching JobDeletor Resource")
	if err := r.Get(ctx, req.NamespacedName, &jd); err != nil {
		log.Error(err, "unable to fetch JobDeletor")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var namespaces v1.NamespaceList
	if jd.Spec.TargetNamespaces == nil {
		if err := r.List(ctx, &namespaces); err != nil {
			log.Error(err, "unable to list Namespaces")
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
	} else {
		for _, n := range jd.Spec.TargetNamespaces {
			var namespace v1.Namespace
			if err := r.Get(ctx, types.NamespacedName{Name: n}, &namespace); err != nil {
				log.Error(err, "unable to get Namespace")
			}
			namespaces.Items = append(namespaces.Items, namespace)
		}
	}

	for _, ns := range namespaces.Items {
		var jobs batchv1.JobList
		if err := r.List(ctx, &jobs, &client.ListOptions{Namespace: ns.Name}); err != nil {
			log.Error(err, "unable to list Jobs")
		}
		for _, job := range jobs.Items {
			if deletable(&job, &jd) {
				background := metav1.DeletePropagationBackground
				if err := r.Delete(ctx, &job, &client.DeleteOptions{PropagationPolicy: &background}); err != nil {
					log.Error(err, "unable to delete Job")
					r.Recorder.Eventf(&jd, v1.EventTypeNormal, "FailedDeleting", "Failed to delete Job %q", job.Name)
				}
				log.Info("deleted Job", "in: "+job.Namespace, job.Name)
				r.Recorder.Eventf(&jd, v1.EventTypeNormal, "Deleted", "Deleted Job %q", job.Name)
				jd.Status.DeletedJobs = append(jd.Status.DeletedJobs, job)
			}
		}
	}

	if len(jd.Status.DeletedJobs) > 10 {
		jd.Status.DeletedJobs = jd.Status.DeletedJobs[len(jd.Status.DeletedJobs)-10:]
	}
	if err := r.Status().Update(ctx, &jd); err != nil {
		log.Error(err, "unable to update JobDeletor status")
		return ctrl.Result{}, err
	}
	r.Recorder.Eventf(&jd, v1.EventTypeNormal, "Updated", "Update jobdeletor.status.DeletedJobs")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *JobDeletorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mikutasv1alpha1.JobDeletor{}).
		Complete(r)
}

func deletable(job *batchv1.Job, jd *mikutasv1alpha1.JobDeletor) bool {
	if jd.Spec.TargetStatus == targetStatusAll || jd.Spec.TargetStatus == targetStatusSucceeded {
		if (job.Spec.Completions != nil && *job.Spec.Completions == job.Status.Succeeded) || job.Status.Succeeded >= 1 {
			return true
		}
	}
	if jd.Spec.TargetStatus == targetStatusAll || jd.Spec.TargetStatus == targetStatusFailed {
		if (job.Spec.BackoffLimit != nil && *job.Spec.BackoffLimit == job.Status.Failed) || job.Status.Failed >= 1 {
			return true
		}
	}
	return false
}
