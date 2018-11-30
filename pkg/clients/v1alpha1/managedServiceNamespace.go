package v1alpha1

import (
	"context"
	integreatlyv1alpha1 "github.com/integr8ly/managed-service-controller/pkg/apis/integreatly/v1alpha1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

const (
	ViewClusterRole = "view"
	ClusterRoleType = "ClusterRole"
)

var log = logf.Log.WithName("ManagedServiceNamespace client")

type managedServiceNamespacesClient struct {
	k8sClient                 client.Client
	osClient                  *ClientFactory
	scheme                    *runtime.Scheme
	managedServiceReconcilers []MsnReconcilerInterface
}

func NewManagedServiceNamespaceClient(mgr manager.Manager, sCfg map[string]map[string]string) MsnReconcilerInterface {
	osClient := NewClientFactory(mgr.GetConfig())
	k8sClient := mgr.GetClient()
	return &managedServiceNamespacesClient{
		k8sClient: k8sClient,
		osClient:  osClient,
		scheme:    mgr.GetScheme(),
		managedServiceReconcilers: []MsnReconcilerInterface{
			NewFuseOperatorManager(k8sClient, osClient, sCfg["fuse"]),
			NewIntegrationControllerManager(k8sClient, osClient, sCfg["integrationController"], mgr),
		},
	}
}

func (msnsc *managedServiceNamespacesClient) Reconcile(msn *integreatlyv1alpha1.ManagedServiceNamespace) error {
	logger := log.WithValues("Namespace", msn.Namespace, "Name", msn.Name)

	if err := msnsc.validate(msn); err != nil {
		return err
	}

	if !msnsc.exists(msn) {
		logger.Info("Creating a new namespace", "Namespace", msn.Name)

		if err := msnsc.createNamespace(msn); err != nil {
			return err
		}

		if err := msnsc.createViewRoleBindingForUser(msn); err != nil {
			return err
		}
	}

	logger.Info("Reconciling Managed Services")
	for _, msr := range msnsc.managedServiceReconcilers {
		err := msr.Reconcile(msn)
		if err != nil {
			return err
		}
	}

	return nil
}

func (msnsc *managedServiceNamespacesClient) exists(msn *integreatlyv1alpha1.ManagedServiceNamespace) bool {
	err := msnsc.k8sClient.Get(context.TODO(), types.NamespacedName{Name: msn.Name}, &corev1.Namespace{})
	if err == nil {
		return true
	}

	return false
}

func (msnsc *managedServiceNamespacesClient) validate(msn *integreatlyv1alpha1.ManagedServiceNamespace) error {
	if len(msn.Spec.ConsumerNamespaces) == 0 {
		return errors.New(" No consumerNamespace set")
	}

	if err := msnsc.validateConsumerNamespaces(msn); err != nil {
		return err
	}

	if err := msnsc.validateUser(msn); err != nil {
		return err
	}

	return nil
}

func (msnsc *managedServiceNamespacesClient) validateConsumerNamespaces(msn *integreatlyv1alpha1.ManagedServiceNamespace) error {
	nsList := &corev1.NamespaceList{}
	err := msnsc.k8sClient.List(context.TODO(), &client.ListOptions{}, nsList)
	if err != nil {
		return err
	}

	for _, a := range msn.Spec.ConsumerNamespaces {
		valid := false
		for _, b := range nsList.Items {
			if a == b.Name {
				valid = true
				break
			}
		}
		if !valid {
			return errors.New("The consumer namespace " + a + " does not exist.")
		}
	}

	return nil
}

func (msnsc *managedServiceNamespacesClient) validateUser(msn *integreatlyv1alpha1.ManagedServiceNamespace) error {
	userClient, err := msnsc.osClient.UserClient()
	if err != nil {
		return err
	}
	if _, err = userClient.Users().Get(msn.Spec.UserID, metav1.GetOptions{}); err != nil {
		if apiErrors.IsNotFound(err) {
			return errors.New("User " + msn.Spec.UserID + " does not exist")
		}
		return err
	}

	return nil
}

func (msnsc *managedServiceNamespacesClient) createNamespace(msn *integreatlyv1alpha1.ManagedServiceNamespace) error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: msn.Name,
		},
	}

	if err := controllerutil.SetControllerReference(msn, ns, msnsc.scheme); err != nil {
		return err
	}

	return msnsc.k8sClient.Create(context.TODO(), ns)
}

func (msnsc *managedServiceNamespacesClient) createViewRoleBindingForUser(msn *integreatlyv1alpha1.ManagedServiceNamespace) error {
	rb := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    msn.Name,
			GenerateName: msn.Spec.UserID + "-view-" + msn.Name + "-",
		},
		RoleRef: rbacv1.RoleRef{
			Kind: ClusterRoleType,
			Name: ViewClusterRole,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: "User",
				Name: msn.Spec.UserID,
			},
		},
	}

	if err := controllerutil.SetControllerReference(msn, rb, msnsc.scheme); err != nil {
		return err
	}

	return msnsc.k8sClient.Create(context.TODO(), rb)
}
