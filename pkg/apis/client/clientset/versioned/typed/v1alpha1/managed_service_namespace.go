package v1alpha1

import (
	integreatly "github.com/integr8ly/managed-services-controller/pkg/apis/integreatly/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ManagedServiceNamespaceInterface interface {
	Create(*integreatly.ManagedServiceNamespace) error
	Exists(msn *integreatly.ManagedServiceNamespace) bool
	Delete(msn *integreatly.ManagedServiceNamespace) error
	Update(msn *integreatly.ManagedServiceNamespace) error
}

// managedServiceNamespaces implements ManagedServiceNamespaceInterface
type managedServiceNamespaces struct {
	client                 kubernetes.Interface
	managedServiceManagers []ManagedServiceManagerInterface
}

type ManagedServiceManagerInterface interface {
	Create(*integreatly.ManagedServiceNamespace) error
	Update(*integreatly.ManagedServiceNamespace) error
}

func NewManagedServiceNamespaces(c kubernetes.Interface) ManagedServiceNamespaceInterface {
	return &managedServiceNamespaces{
		client: c,
		managedServiceManagers: []ManagedServiceManagerInterface{
			NewFuseManager(),
			NewIntegrationControllerManager(c),
        },
	}
}

func (msns *managedServiceNamespaces) Create(msn *integreatly.ManagedServiceNamespace) error {
	err := createNamespace(msns.client, msn.Spec.ManagedNamespace);if err != nil {
		return err
	}

	for _, ms := range msns.managedServiceManagers {
		err = ms.Create(msn);if err != nil {
			return err
		}
	}

	return nil
}

func (msns *managedServiceNamespaces) Exists(msn *integreatly.ManagedServiceNamespace) bool {
	_, err := msns.client.Core().Namespaces().Get(msn.Spec.ManagedNamespace, metav1.GetOptions{});if err != nil {
		return false
	}

	return true
}

func (msns *managedServiceNamespaces) Delete(msn *integreatly.ManagedServiceNamespace) error {
	return msns.client.Core().Namespaces().Delete(msn.Spec.ManagedNamespace, &metav1.DeleteOptions{})
}

func (msns *managedServiceNamespaces) Update(msn *integreatly.ManagedServiceNamespace) error {
	for _, ms := range msns.managedServiceManagers {
		err := ms.Update(msn);if err != nil {
			return err
		}
	}

	return nil
}

func createNamespace(c kubernetes.Interface, namespace string) error{
	_, err := c.Core().Namespaces().Create(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	})
	return err
}
