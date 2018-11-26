package v1alpha1

import (
	integreatly "github.com/integr8ly/managed-service-controller/pkg/apis/integreatly/v1alpha1"
	appsv1 "github.com/openshift/api/apps/v1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

type integrationControllerManager struct {
	k8sClient       kubernetes.Interface
	osClientFactory *ClientFactory
	cfg             map[string]string
}

func NewIntegrationControllerManager(client kubernetes.Interface, oscf *ClientFactory, cfg map[string]string) ManagedServiceManagerInterface {
	return &integrationControllerManager{
		k8sClient:       client,
		osClientFactory: oscf,
		cfg:             cfg,
	}
}

func (icm *integrationControllerManager) Create(msn *integreatly.ManagedServiceNamespace) error {
	ns := msn.Name
	if err := icm.createRoleBinding(EnmasseNamespace, getEnmasseConfigMapRoleBinding(ns, icm.cfg)); err != nil {
		return err
	}

	if err := icm.createRoleBinding(ns, getUpdateIntegrationsRoleBinding(msn)); err != nil {
		return err
	}

	if err := icm.createIntegrationController(ns); err != nil {
		return err
	}

	return nil
}

func (icm *integrationControllerManager) Update(msn *integreatly.ManagedServiceNamespace) error {
	ns := msn.Name
	dcClient, err := icm.osClientFactory.AppsClient()
	if err != nil {
		return err
	}

	d, err := dcClient.DeploymentConfigs(ns).Get(icm.cfg["name"], metav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to create an openshift deployment config client")
	}

	c := getContainer(icm.cfg["name"], d)
	if c != nil {
		e := getEnvVar(icm.cfg["namespacesEnvVarKey"], c)
		oldValues := filter(strings.Split(e.Value, ","), func(s string) bool {
			// strings.Split can return an empty string if e.Value is empty.
			return len(s) != 0
		})
		added, deleted := arrayDiff(oldValues, msn.Spec.ConsumerNamespaces)
		if len(deleted) > 0 {
			if err := icm.deleteRoutesAndServicesRoleBindings(deleted); err != nil {
				return err
			}

		}

		if len(added) > 0 {
			if err := icm.createRoutesAndServicesRoleBindings(added, ns); err != nil {
				return err
			}
		}

		if len(added) > 0 || len(deleted) > 0 {
			setEnvVarValue(e, strings.Join(msn.Spec.ConsumerNamespaces, ","), c)
			if _, err = dcClient.DeploymentConfigs(ns).Update(d); err != nil {
				return err
			}
		}
	}

	return nil
}

func arrayDiff(old, new []string) ([]string, []string) {
	added := filter(new, func(s string) bool {
		return !contains(old, s)
	})
	deleted := filter(old, func(s string) bool {
		return !contains(new, s)
	})

	return added, deleted
}

type predicateFunc func(string) bool

func filter(a []string, predicate predicateFunc) []string {
	var result []string
	for _, s := range a {
		if predicate(s) {
			result = append(result, s)
		}
	}

	return result
}

func getContainer(name string, d *appsv1.DeploymentConfig) *corev1.Container {
	for _, c := range d.Spec.Template.Spec.Containers {
		if c.Name == name {
			return &c
		}
	}

	return nil
}

func contains(s []string, b string) bool {
	for _, a := range s {
		if a == b {
			return true
		}
	}
	return false
}

func setEnvVarValue(envVar *corev1.EnvVar, value string, c *corev1.Container) {
	for i, e := range c.Env {
		if e.Name == envVar.Name {
			c.Env[i].Value = value
		}
	}
}

func getEnvVar(name string, c *corev1.Container) *corev1.EnvVar {
	for _, e := range c.Env {
		if e.Name == name {
			return &e
		}
	}

	return nil
}

func (icm *integrationControllerManager) createRoutesAndServicesRoleBindings(namespaces []string, managedServiceNamespace string) error {
	for _, cns := range namespaces {
		rb := getRoutesAndServicesRoleBinding(managedServiceNamespace, icm.cfg)
		if err := icm.createRoleBinding(cns, rb); err != nil {
			return err
		}
	}

	return nil
}
func (icm *integrationControllerManager) deleteRoutesAndServicesRoleBindings(namespaces []string) error {
	for _, ns := range namespaces {
		err := icm.k8sClient.Rbac().RoleBindings(ns).DeleteCollection(&metav1.DeleteOptions{}, metav1.ListOptions{
			LabelSelector: "for=route-services",
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (icm *integrationControllerManager) createIntegrationController(namespace string) error {
	if _, err := icm.k8sClient.CoreV1().ServiceAccounts(namespace).Create(integrationServiceAccount); err != nil {
		return errors.Wrap(err, "failed to create service account for Integration Controller")
	}

	if _, err := icm.k8sClient.RbacV1beta1().RoleBindings(namespace).Create(integrationServiceRoleBinding); err != nil {
		return errors.Wrap(err, "failed to create role binding for integration controller service")
	}

	dcClient, err := icm.osClientFactory.AppsClient()
	if err != nil {
		return errors.Wrap(err, "failed to create an openshift deployment config client")
	}

	if _, err = dcClient.DeploymentConfigs(namespace).Create(getIntegrationDeploymentConfig(icm.cfg)); err != nil {
		return errors.Wrap(err, "failed to create deployment config for integration controller service")
	}

	return nil
}

func (icm *integrationControllerManager) createRoleBinding(namespace string, rb *rbacv1.RoleBinding) error {
	_, err := icm.k8sClient.Rbac().RoleBindings(namespace).Create(rb)
	if err != nil {
		return err
	}

	return nil
}
