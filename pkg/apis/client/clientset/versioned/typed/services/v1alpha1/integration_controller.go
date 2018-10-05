package v1alpha1

import (
	resources "github.com/integr8ly/managed-services-controller/pkg/apis/client/clientset/versioned/typed/resources/v1alpha1"
	olm "github.com/integr8ly/managed-services-controller/pkg/apis/olm/v1alpha1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	IntegrationControllerInstallPlanName           = "integration-controller.0.0.1-install"
	IntegrationControllerClusterServiceVersionName = "integration-controller-0.0.1"
	EnmasseNamespace                               = "enmasse"
	EnmasseClusterRoleName                         = "enmasse-integration-viewer"
	RoutesAndServicesClusterRoleName               = "route-service-viewer"
	IntegrationControllerName                      = "integration-controller"
)

type IntegrationControllerInterface interface {
	Create() error
}

// integrationController implements IntegrationControllerInterface
type integrationController struct {
	Client    kubernetes.Interface
	Namespace string
}

func NewIntegrationController(client kubernetes.Interface, namespace string) FuseOperatorInterface {
	return &integrationController{
		Client: client,
		Namespace: namespace,
	}
}

func (ic *integrationController) Create() error {
	err := ic.createEnmasseConfigMapRoleBinding();if err != nil {
		return err
	}

	err = ic.createRoutesAndServicesRoleBinding();if err != nil {
		return err
	}

	err = ic.createIntegrationControllerInstallPlan();if err != nil {
		return err
	}


	return nil
}

func (ic *integrationController) createEnmasseConfigMapRoleBinding() error {
	rb := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: IntegrationControllerName + "-enmasse-view-",
			Namespace: EnmasseNamespace,
			Labels: map[string]string{
				"for": IntegrationControllerName,
			},
		},
		RoleRef: clusterRole(EnmasseClusterRoleName),
		Subjects: []rbacv1.Subject{
			serviceAccountSubject(ic.Namespace),
		},
	}

	err := ic.createRoleBinding(EnmasseNamespace, rb);if err != nil {
		return err
	}

	return nil
}

func (ic *integrationController) createRoutesAndServicesRoleBinding() error {
	rb :=  &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: IntegrationControllerName + "-route-viewer-",
		},
		RoleRef: clusterRole(RoutesAndServicesClusterRoleName),
		Subjects: []rbacv1.Subject{
			serviceAccountSubject(ic.Namespace),
		},
	}

	err := ic.createRoleBinding(ic.Namespace, rb);if err != nil {
		return err
	}

	return nil
}

func clusterRole(roleName string) rbacv1.RoleRef {
	return rbacv1.RoleRef{
		Kind: "ClusterRole",
		Name: roleName,
	}
}

func serviceAccountSubject(namespace string) rbacv1.Subject {
	return rbacv1.Subject{
		Kind: "ServiceAccount",
		Name: IntegrationControllerName,
		Namespace: namespace,
	}
}

func (ic *integrationController) createIntegrationControllerInstallPlan() error {
	ip := &olm.InstallPlan{
		TypeMeta: metav1.TypeMeta{
			APIVersion: olm.SchemeGroupVersion.String(),
			Kind: olm.InstallPlanKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: IntegrationControllerInstallPlanName,
			Namespace: ic.Namespace,
		},
		Spec: olm.InstallPlanSpec{
			Approval: olm.ApprovalsAutomatic,
			ClusterServiceVersionNames: []string{
				IntegrationControllerClusterServiceVersionName,
			},
		},
	}

	ips := resources.NewInstallPlans(ic.Namespace)
	_, err := ips.Create(ip);if err != nil {
		return err
	}

	return nil
}

func (ic *integrationController) createRoleBinding(namespace string, rb *rbacv1.RoleBinding) error {
	_, err := ic.Client.Rbac().RoleBindings(namespace).Create(rb);if err != nil {
		return err
	}

	return nil
}