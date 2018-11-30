package managedservicenamespace

import (
	"context"
	"encoding/json"
	"github.com/gobuffalo/packr"
	integreatlyv1alpha1 "github.com/integr8ly/managed-service-controller/pkg/apis/integreatly/v1alpha1"
	msnsClients "github.com/integr8ly/managed-service-controller/pkg/clients/v1alpha1"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_managedservicenamespace")

// Add creates a new ManagedServiceNamespace Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileManagedServiceNamespace{
		client:     mgr.GetClient(),
		scheme:     mgr.GetScheme(),
		msnsClient: getMsnsClient(mgr),
	}
}

// Get the managed service namespace client
func getMsnsClient(mgr manager.Manager) msnsClients.MsnReconcilerInterface {
	config := packr.NewBox("../../../config")
	sCfgBytes, err := config.Find("service-config.json")
	if err != nil {
		logrus.Fatalf("failed to get managed service config: %v", err)
	}

	var sCfg map[string]map[string]string
	json.Unmarshal(sCfgBytes, &sCfg)

	return msnsClients.NewManagedServiceNamespaceClient(mgr, sCfg)
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("managedservicenamespace-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ManagedServiceNamespace
	err = c.Watch(&source.Kind{Type: &integreatlyv1alpha1.ManagedServiceNamespace{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// ReconcileManagedServiceNamespace reconciles a ManagedServiceNamespace object
type ReconcileManagedServiceNamespace struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client     client.Client
	scheme     *runtime.Scheme
	msnsClient msnsClients.MsnReconcilerInterface
}

// Reconcile reads that state of the cluster for a ManagedServiceNamespace object and makes changes based on the state read
// and what is in the ManagedServiceNamespace.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileManagedServiceNamespace) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ManagedServiceNamespace")

	// Fetch the ManagedServiceNamespace instance
	msn := &integreatlyv1alpha1.ManagedServiceNamespace{}
	err := r.client.Get(context.TODO(), request.NamespacedName, msn)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if err := r.msnsClient.Reconcile(msn); err != nil {
		return reconcile.Result{}, err
	}

	reqLogger.Info("ManagedServiceNamespace " + msn.Name + " reconciled")
	return reconcile.Result{}, nil
}
