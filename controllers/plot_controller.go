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
	"fmt"
	"reflect"

	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	directorv1 "github.com/kuberik/director/api/v1"
	"github.com/kuberik/director/controllers/frameutil"
	v1 "k8s.io/api/core/v1"
)

// PlotReconciler reconciles a Plot object
type PlotReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

const plotNameLabel = "director.kuberik.io/plot"
const frameNameLabel = "director.kuberik.io/frame"

//+kubebuilder:rbac:groups=director.kuberik.io,resources=plots,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=director.kuberik.io,resources=plots/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=director.kuberik.io,resources=plots/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Plot object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *PlotReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("plot", req.NamespacedName)

	// your logic here
	plot := &directorv1.Plot{}
	err := r.Client.Get(ctx, req.NamespacedName, plot)
	if err != nil {
		return ctrl.Result{}, err
	}

	r.updateFrameStatuses(ctx, plot)
	// TODO: finalizer

frames:
	for _, f := range plot.Spec.Frames {
		if frameStatus := frameutil.FindFrameStatus(plot.Status.FrameStatuses, f.Name); frameStatus == nil {
			// Start frame
			for _, causalFrame := range f.CausedBy {
				// TODO: enable globbing
				status := frameutil.FindFrameStatus(plot.Status.FrameStatuses, causalFrame)
				if status == nil || status.State.Running != nil || !status.State.Finished.Success {
					// Skip because dependencies are either not finished or failed
					continue frames
				}
			}
			job := frameJob(plot, &f)
			err = r.Client.Create(ctx, job)
			if err != nil && !errors.IsAlreadyExists(err) {
				return ctrl.Result{}, err
			}
		} else if frameStatus.State.Running != nil {
			// Do nothing
		} else if frameStatus.State.Finished != nil {
			// Do nothing
		}
	}

	err = r.Client.Status().Update(ctx, plot)
	if err != nil {
		return ctrl.Result{Requeue: true}, nil
	}
	return ctrl.Result{}, nil
}

func (r *PlotReconciler) updateFrameStatuses(ctx context.Context, plot *directorv1.Plot) {
	jobs := &batchv1.JobList{}
	r.Client.List(ctx, jobs, &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(map[string]string{
			plotNameLabel: plot.Name,
		}),
	})

	for _, job := range jobs.Items {
		frameName := job.Labels[frameNameLabel]
		for _, condition := range job.Status.Conditions {
			if condition.Type == batchv1.JobFailed || condition.Type == batchv1.JobComplete {
				frameutil.SetFrameStatus(&plot.Status.FrameStatuses, directorv1.FrameStatus{
					Name: frameName,
					State: directorv1.FrameState{
						Finished: &directorv1.FrameStateFinished{
							// StartedAt:  frameutil.FindFrameStatus(plot.Status.FrameStatuses, frameName).State.Running.StartedAt,
							FinishedAt: condition.LastTransitionTime,
							Success:    condition.Type == batchv1.JobComplete,
						},
					},
				})
				break
			}
		}
		if status := frameutil.FindFrameStatus(plot.Status.FrameStatuses, frameName); status == nil && job.Status.StartTime != nil {
			frameutil.SetFrameStatus(&plot.Status.FrameStatuses, directorv1.FrameStatus{
				Name: frameName,
				State: directorv1.FrameState{
					Running: &directorv1.FrameStateRunning{
						StartedAt: *job.Status.StartTime,
					},
				},
			})
		}
	}
}

var zero int32 = 0

func frameJob(plot *directorv1.Plot, frame *directorv1.Frame) *batchv1.Job {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", plot.Name, frame.Name),
			Namespace: plot.Namespace,
			Labels: map[string]string{
				plotNameLabel:  plot.Name,
				frameNameLabel: frame.Name,
			},
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(
					plot, directorv1.GroupVersion.WithKind(reflect.ValueOf(plot).Elem().Type().Name()),
				),
			},
		},
		Spec: frame.Action,
	}

	// Set defaults
	// TODO: this could probably be done more elegantly
	if job.Spec.BackoffLimit == nil {
		job.Spec.BackoffLimit = &zero
	}
	if job.Spec.Template.Spec.RestartPolicy == "" {
		job.Spec.Template.Spec.RestartPolicy = v1.RestartPolicyNever
	}

	return job
}

// SetupWithManager sets up the controller with the Manager.
func (r *PlotReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&directorv1.Plot{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
