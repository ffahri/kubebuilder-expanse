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
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"simplekubebuilder/pkg/database"
	"time"

	expansev1beta1 "simplekubebuilder/api/v1beta1"
)

// SpaceShipsReconciler reconciles a SpaceShips object
type SpaceShipsReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	DBConfig *database.DBConfig
}

// +kubebuilder:rbac:groups=expanse.blog.webischia.com,resources=spaceships,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=expanse.blog.webischia.com,resources=spaceships/status,verbs=get;update;patch

func (r *SpaceShipsReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("spaceships", req.NamespacedName)
	spaceship := expansev1beta1.SpaceShips{}
	if err := r.Get(context.Background(), req.NamespacedName, &spaceship); err != nil {
		return ctrl.Result{}, nil
	}
	switch spaceship.Status.Phase {
	case "":
		// with new created ones are doesn't have any status so its empty
		return r.createSpaceship(spaceship)
	case expansev1beta1.Created:
		return r.startTheEngines(spaceship)
	}
	return ctrl.Result{}, nil
}

func (r *SpaceShipsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&expansev1beta1.SpaceShips{}).
		Complete(r)
}

func (r *SpaceShipsReconciler) createSpaceship(ship expansev1beta1.SpaceShips) (ctrl.Result, error) {
	ship.Status.Phase = expansev1beta1.Created
	if err := r.DBConfig.Write(ship); err != nil {
		r.Log.Error(err, "Could not write entity into DDB")
		return ctrl.Result{}, err
	}
	return r.setStatus(ship)
}

func (r *SpaceShipsReconciler) startTheEngines(ship expansev1beta1.SpaceShips) (ctrl.Result, error) {
	ship.Status.Phase = expansev1beta1.Active
	if err := r.DBConfig.Update(ship); err != nil {
		r.Log.Error(err, "Could not write entity into DDB")
		return ctrl.Result{}, err
	}
	return r.setStatus(ship)

}

func (r *SpaceShipsReconciler) setStatus(ship expansev1beta1.SpaceShips) (ctrl.Result, error) {
	if err := r.Client.Status().Update(context.Background(), &ship); err != nil {
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: 300 * time.Millisecond,
		}, err
	}
	return ctrl.Result{}, nil
}
