package v1alpha1

import (
	integreatly "github.com/integr8ly/managed-services-controller/pkg/apis/integreatly/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ManagedServiceNamespaceInterface interface {
	Create(*integreatly.ManagedServiceNamespace) error
	Exists(msn *integreatly.ManagedServiceNamespace) bool
	Delete(msn *integreatly.ManagedServiceNamespace) error
	Update(msn *integreatly.ManagedServiceNamespace) error
}

type managedServiceNamespacesClient struct {
	k8sClient              kubernetes.Interface
	managedServiceManagers []ManagedServiceManagerInterface
}

type ManagedServiceManagerInterface interface {
	Create(*integreatly.ManagedServiceNamespace) error
	Update(*integreatly.ManagedServiceNamespace) error
}

func NewManagedServiceNamespaces(c kubernetes.Interface) ManagedServiceNamespaceInterface {
	return &managedServiceNamespacesClient{
		k8sClient: c,
		managedServiceManagers: []ManagedServiceManagerInterface{
			NewFuseOperatorManager(),
			NewIntegrationControllerManager(c),
        },
	}
}

func (msnsc *managedServiceNamespacesClient) Create(msn *integreatly.ManagedServiceNamespace) error {
	if err := createNamespace(msnsc.k8sClient, msn.Spec.ManagedNamespace);err != nil {
		return err
	}

	for _, ms := range msnsc.managedServiceManagers {
		if err := ms.Create(msn);err != nil {
			return err
		}
	}

	return nil
}

func (msnsc *managedServiceNamespacesClient) Exists(msn *integreatly.ManagedServiceNamespace) bool {
	_, err := msnsc.k8sClient.Core().Namespaces().Get(msn.Spec.ManagedNamespace, metav1.GetOptions{})
	if err != nil && errors.IsAlreadyExists(err) {
		return true
	}

	return false
}

func (msnsc *managedServiceNamespacesClient) Delete(msn *integreatly.ManagedServiceNamespace) error {
	return msnsc.k8sClient.Core().Namespaces().Delete(msn.Spec.ManagedNamespace, &metav1.DeleteOptions{})
}

func (msnsc *managedServiceNamespacesClient) Update(msn *integreatly.ManagedServiceNamespace) error {
	for _, ms := range msnsc.managedServiceManagers {
		if err := ms.Update(msn);err != nil {
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
