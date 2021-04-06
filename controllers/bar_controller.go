/*


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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	foo "github.com/thetirefire/example/api/v1"
)

// BarReconciler reconciles a Bar object.
type BarReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=foo.example.thetirefire,resources=bars,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=foo.example.thetirefire,resources=bars/status,verbs=get;update;patch

func (r *BarReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	log.Info("Reconciling Bar")

	var bar foo.Bar
	if err := r.Get(ctx, req.NamespacedName, &bar); err != nil {
		log.Error(err, "unable to fetch Bar")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// TODO: handle deletion
	// TODO: test if we need to create bar and create if needed

	// Set Status
	bar.Status.Path = fmt.Sprintf("/%s/%s", bar.Namespace, bar.Name)

	if err := r.Status().Update(ctx, &bar); err != nil {
		log.Error(err, "unable to update Bar status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *BarReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager) error {
	if r.Scheme == nil {
		r.Scheme = mgr.GetScheme()
	}

	if r.Client == nil {
		r.Client = mgr.GetClient()
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&foo.Bar{}).
		Complete(r)
}
