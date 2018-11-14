package v1alpha1

import (
	integreatly "github.com/integr8ly/managed-service-controller/pkg/apis/integreatly/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/k8sclient"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
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
	// TODO: Remove
	namespaces             <-chan watch.Event
}

type ManagedServiceManagerInterface interface {
	Create(*integreatly.ManagedServiceNamespace) error
	Update(*integreatly.ManagedServiceNamespace) error
}

func NewManagedServiceNamespaceClient(cfg *rest.Config) ManagedServiceNamespaceInterface {
	k8sClient := k8sclient.GetKubeClient()
	osClient := NewClientFactory(cfg)
	return &managedServiceNamespacesClient{
		k8sClient: k8sClient,
		managedServiceManagers: []ManagedServiceManagerInterface{
			NewFuseOperatorManager(k8sClient, osClient),
			NewIntegrationControllerManager(k8sClient, osClient),
		},
		namespaces: nil,
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
	ns, err := msnsc.k8sClient.Core().Namespaces().Get(msn.Name, metav1.GetOptions{})
	if err == nil && ns != nil {
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

func (msnsc *managedServiceNamespacesClient) Validate(msn *integreatly.ManagedServiceNamespace) error {
	if len(msn.Spec.ConsumerNamespaces) == 0 {
		return errors.New("ManagedServiceNamespace: " + msn.Name + " has no consumerNamespace set")
	}

	nsList, err := msnsc.k8sClient.CoreV1().Namespaces().List(metav1.ListOptions{}); if err != nil {
		return err
	}

	return msn.Validate(nsList)
	// TODO: Use a chan?
	// TODO: Validate the user
	//if msnsc.namespaces == nil {
	//	logrus.Info("namespaces are nil")
	//	var event watch.Event
	//	var test *corev1.NamespaceList
	//	var ok bool
	//	nsWatch, err := msnsc.k8sClient.CoreV1().Namespaces().Watch(metav1.ListOptions{}); if err != nil {
	//		return err
	//	}
	//	switch <-nsWatch.ResultChan() {
	//	case event:
	//		test, ok = event.Object.(*corev1.NamespaceList)
	//	}
	//	//msnsc.namespaces = nsWatch.ResultChan()
	//	//namespaces := <- msnsc.namespaces
	//	//namespaceslist, ok := namespaces.Object.(*corev1.NamespaceList)
	//	logrus.Info("namespaces are nil %v", ok)
	//	if ok {
	//		return msn.Validate(test)
	//	}
	//}



	//return nil
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
