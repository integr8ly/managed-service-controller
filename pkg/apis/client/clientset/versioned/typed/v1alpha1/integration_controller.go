package v1alpha1

import (
	integreatly "github.com/integr8ly/managed-services-controller/pkg/apis/integreatly/v1alpha1"
	olm "github.com/integr8ly/managed-services-controller/pkg/apis/olm/v1alpha1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

const (
	IntegrationControllerInstallPlanName           = "integration-controller.0.0.1-install"
	IntegrationControllerClusterServiceVersionName = "integration-controller-0.0.1"
	EnmasseNamespace                               = "enmasse"
	EnmasseClusterRoleName                         = "enmasse-integration-viewer"
	RoutesAndServicesClusterRoleName               = "route-service-viewer"
	IntegrationControllerName                      = "integration-controller"
	IntegrationUserNamespacesEnvVarKey             = "USER_NAMESPACES"
)

type integrationController struct {
	Client    kubernetes.Interface
}

func NewIntegrationControllerManager(client kubernetes.Interface) ManagedServiceManagerInterface {
	return &integrationController{
		Client: client,
	}
}

func (ic *integrationController) Create(msn *integreatly.ManagedServiceNamespace) error {
	ns := msn.Spec.ManagedNamespace
	err := ic.createEnmasseConfigMapRoleBinding(ns);if err != nil {
		return err
	}

	err = ic.createRoutesAndServicesRoleBinding(ns);if err != nil {
		return err
	}

	// When creating a new Integration Controller there should be only one ConsumerNamespace.
	// Still not sure about this.
	cns := msn.Spec.ConsumerNamespaces[0]
	err = ic.createIntegrationControllerInstallPlan(cns);if err != nil {
		return err
	}

	return nil
}

func (ic *integrationController) Update(msn *integreatly.ManagedServiceNamespace) error {
	d, err := ic.Client.Apps().Deployments(msn.Spec.ManagedNamespace).Get(IntegrationControllerName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	for _, c := range d.Spec.Template.Spec.Containers {
		if c.Name == IntegrationControllerName {
			for _, e := range c.Env {
				if e.Name == IntegrationUserNamespacesEnvVarKey {
					e.Value = strings.Join(msn.Spec.ConsumerNamespaces, ",")
				}
			}
			_, err := ic.Client.Apps().Deployments(msn.Spec.ManagedNamespace).Update(d);if err != nil {
				return err
			}
		}
	}

	return nil
}

func (ic *integrationController) createEnmasseConfigMapRoleBinding(namespace string) error {
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
			serviceAccountSubject(namespace),
		},
	}

	err := ic.createRoleBinding(EnmasseNamespace, rb);if err != nil {
		return err
	}

	return nil
}

func (ic *integrationController) createRoutesAndServicesRoleBinding(namespace string) error {
	rb :=  &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: IntegrationControllerName + "-route-services-",
		},
		RoleRef: clusterRole(RoutesAndServicesClusterRoleName),
		Subjects: []rbacv1.Subject{
			serviceAccountSubject(namespace),
		},
	}

	err := ic.createRoleBinding(namespace, rb);if err != nil {
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

func (ic *integrationController) createIntegrationControllerInstallPlan(namespace string) error {
	ip := &olm.InstallPlan{
		TypeMeta: metav1.TypeMeta{
			APIVersion: olm.SchemeGroupVersion.String(),
			Kind: olm.InstallPlanKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: IntegrationControllerInstallPlanName,
			Namespace: namespace,
		},
		Spec: olm.InstallPlanSpec{
			Approval: olm.ApprovalsAutomatic,
			ClusterServiceVersionNames: []string{
				IntegrationControllerClusterServiceVersionName,
			},
		},
	}

	ips := NewInstallPlans(namespace)
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