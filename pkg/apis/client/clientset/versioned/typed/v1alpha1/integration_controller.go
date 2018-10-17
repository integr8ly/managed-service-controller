package v1alpha1

import (
	"errors"
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

type integrationControllerManager struct {
	k8sClient kubernetes.Interface
}

func NewIntegrationControllerManager(client kubernetes.Interface) ManagedServiceManagerInterface {
	return &integrationControllerManager{
		k8sClient: client,
	}
}

func (icm *integrationControllerManager) Create(msn *integreatly.ManagedServiceNamespace) error {
	if len(msn.Spec.ConsumerNamespaces) == 0 {
		return errors.New("There must be a ConsumerNamespace set")
	}

	ns := msn.Name
	if err := icm.createEnmasseConfigMapRoleBinding(ns);err != nil {
		return err
	}

	cns := msn.Spec.ConsumerNamespaces[0]
	if err := icm.createRoutesAndServicesRoleBinding(cns, ns);err != nil {
		return err
	}

	if err := icm.createIntegrationControllerInstallPlan(ns);err != nil {
		return err
	}

	return nil
}

func (icm *integrationControllerManager) Update(msn *integreatly.ManagedServiceNamespace) error {
	d, err := icm.k8sClient.Apps().Deployments(msn.Name).Get(IntegrationControllerName, metav1.GetOptions{})
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
			_, err := icm.k8sClient.Apps().Deployments(msn.Name).Update(d);if err != nil {
				return err
			}
		}
	}

	return nil
}

func (icm *integrationControllerManager) createEnmasseConfigMapRoleBinding(namespace string) error {
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

	if err := icm.createRoleBinding(EnmasseNamespace, rb);err != nil {
		return err
	}

	return nil
}

func (icm *integrationControllerManager) createRoutesAndServicesRoleBinding(consumerNamespace, managedServiceNamespace string) error {
	rb :=  &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: IntegrationControllerName + "-route-services-",
		},
		RoleRef: clusterRole(RoutesAndServicesClusterRoleName),
		Subjects: []rbacv1.Subject{
			serviceAccountSubject(managedServiceNamespace),
		},
	}

	if err := icm.createRoleBinding(consumerNamespace, rb);err != nil {
		return err
	}

	return nil
}

func (icm *integrationControllerManager) createIntegrationControllerInstallPlan(namespace string) error {
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

func (icm *integrationControllerManager) createRoleBinding(namespace string, rb *rbacv1.RoleBinding) error {
	_, err := icm.k8sClient.Rbac().RoleBindings(namespace).Create(rb);if err != nil {
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