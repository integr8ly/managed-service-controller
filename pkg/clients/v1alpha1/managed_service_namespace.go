package v1alpha1

import (
	integreatly "github.com/integr8ly/managed-services-controller/pkg/apis/integreatly/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/k8sclient"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	ViewClusterRole = "view"
	ClusterRoleType = "ClusterRole"
)

type managedServiceNamespacesClient struct {
	k8sClient              kubernetes.Interface
	managedServiceManagers []ManagedServiceManagerInterface
}

type ManagedServiceManagerInterface interface {
	Create(*integreatly.ManagedServiceNamespace) error
	Update(*integreatly.ManagedServiceNamespace) error
}

func NewManagedServiceNamespaces(cfg *rest.Config) ManagedServiceNamespaceInterface {
	k8sClient := k8sclient.GetKubeClient()
	osClient := NewClientFactory(cfg)
	return &managedServiceNamespacesClient{
		k8sClient: k8sClient,
		managedServiceManagers: []ManagedServiceManagerInterface{
			NewFuseOperatorManager(k8sClient, osClient),
			NewIntegrationControllerManager(k8sClient),
		},
	}
}

func (msnsc *managedServiceNamespacesClient) Create(msn *integreatly.ManagedServiceNamespace) error {
	if err := createNamespace(msnsc.k8sClient, msn.Name); err != nil {
		return err
	}

	if err := createViewRoleBindingForUser(msnsc.k8sClient, msn); err != nil {
		return err
	}

	for _, ms := range msnsc.managedServiceManagers {
		if err := ms.Create(msn); err != nil {
			return err
		}
	}

	return nil
}

func (msnsc *managedServiceNamespacesClient) Exists(msn *integreatly.ManagedServiceNamespace) bool {
	_, err := msnsc.k8sClient.Core().Namespaces().Get(msn.Name, metav1.GetOptions{})
	if err != nil && errors.IsAlreadyExists(err) {
		return true
	}

	return false
}

func (msnsc *managedServiceNamespacesClient) Delete(msn *integreatly.ManagedServiceNamespace) error {
	return msnsc.k8sClient.Core().Namespaces().Delete(msn.Name, &metav1.DeleteOptions{})
}

func (msnsc *managedServiceNamespacesClient) Update(msn *integreatly.ManagedServiceNamespace) error {
	for _, ms := range msnsc.managedServiceManagers {
		if err := ms.Update(msn); err != nil {
			return err
		}
	}

	return nil
}

func createNamespace(c kubernetes.Interface, namespace string) error {
	_, err := c.Core().Namespaces().Create(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	})
	return err
}

func createViewRoleBindingForUser(c kubernetes.Interface, msn *integreatly.ManagedServiceNamespace) error {
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
	if _, err := c.Rbac().RoleBindings(msn.Name).Create(rb); err != nil {
		return err
	}

	return nil
}
