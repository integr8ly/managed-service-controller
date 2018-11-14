package v1alpha1

import (
	integreatly "github.com/integr8ly/managed-service-controller/pkg/apis/integreatly/v1alpha1"
	appsv1 "github.com/openshift/api/apps/v1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

// TODO: remove and use object names where possible. move to api objects as well
const (
	EnmasseNamespace                   = "enmasse"
	EnmasseClusterRoleName             = "enmasse-integration-viewer"
	RoutesAndServicesClusterRoleName   = "route-service-viewer"
	IntegrationControllerName          = "integration-controller"
	IntegrationUserNamespacesEnvVarKey = "USER_NAMESPACES"
)

type integrationControllerManager struct {
	k8sClient kubernetes.Interface
	osClientFactory *ClientFactory
}

func NewIntegrationControllerManager(client kubernetes.Interface, oscf *ClientFactory) ManagedServiceManagerInterface {
	return &integrationControllerManager{
		k8sClient: client,
		osClientFactory: oscf,
	}
}

func (icm *integrationControllerManager) Create(msn *integreatly.ManagedServiceNamespace) error {
	// TODO: Should be part of validate
	if len(msn.Spec.ConsumerNamespaces) == 0 {
		return errors.New("There must be a ConsumerNamespace set")
	}

	ns := msn.Name
	// TODO: We my need some validation of enmasse but maybe adding the person to a group may solve this issue of needing an enmasse namespace?
	if err := icm.createEnmasseConfigMapRoleBinding(ns); err != nil {
		return err
	}

	if err := icm.createRoutesAndServicesRoleBindings(msn.Spec.ConsumerNamespaces, ns); err != nil {
		return err
	}

	if err := icm.createUpdateIntegrationsRoleBinding(msn); err != nil {
		return err
	}

	if err := icm.createIntegrationController(ns); err != nil {
		return err
	}

	return nil
}

//func (icm *integrationControllerManager) Update(msn *integreatly.ManagedServiceNamespace) error {
//	ns := msn.Name
//	dcClient, err := icm.osClientFactory.AppsClient(); if err != nil {
//		return err
//	}
//	// TODO: Use the object names instead of these constants where possible
//	d, err := dcClient.DeploymentConfigs(ns).Get(IntegrationControllerName, metav1.GetOptions{}); if err != nil {
//		return err
//	}
//
//	c := getContainer(IntegrationControllerName, d)
//	if c != nil {
//		updated, env := updateEnvVar(IntegrationUserNamespacesEnvVarKey, strings.Join(msn.Spec.ConsumerNamespaces, ","), c)
//		if updated {
//			if err := icm.deleteRoutesAndServicesRoleBindings(strings.Split(env.Value, ",")); err != nil {
//				return err
//			}
//
//			if err := icm.createRoutesAndServicesRoleBindings(msn.Spec.ConsumerNamespaces, ns); err != nil {
//				return err
//			}
//
//			if _, err = dcClient.DeploymentConfigs(ns).Update(d); err != nil {
//				return err
//			}
//		}
//	}
//
//	return nil
//}

func (icm *integrationControllerManager) Update(msn *integreatly.ManagedServiceNamespace) error {
	logrus.Info("update called")
	ns := msn.Name
	// TODO: getter?
	dcClient, err := icm.osClientFactory.AppsClient(); if err != nil {
		return err
	}
	// TODO: Use the object names instead of these constants where possible
	d, err := dcClient.DeploymentConfigs(ns).Get(IntegrationControllerName, metav1.GetOptions{}); if err != nil {
		// TODO: wrap all these as the errors are vague
		logrus.Info("error getting the deployment")
		return err
	}

	c := getContainer(IntegrationControllerName, d)
	if c != nil {
		e := getEnvVar(IntegrationUserNamespacesEnvVarKey, c)
		oldValues := filter(strings.Split(e.Value, ","), func(s string) bool {
			// filter empty values
			return len(s) != 0
		})
		added, deleted := arrayDiff(oldValues, msn.Spec.ConsumerNamespaces)
		logrus.Printf("Added: %v", added)
		logrus.Printf("Deleted: %v", deleted)
		if len(deleted) > 0 {
			if err := icm.deleteRoutesAndServicesRoleBindings(deleted); err != nil {
				logrus.Info("error deleteRoutesAndServicesRoleBindings")
				return err
			}

		}

		logrus.Printf("Length Added: %b", len(added))
		logrus.Printf("Length Deleted: %b", len(deleted))
		if len(added) > 0 {
			logrus.Print("Calling delete:")
			if err := icm.createRoutesAndServicesRoleBindings(added, ns); err != nil {
				logrus.Info("error createRoutesAndServicesRoleBindings")
				return err
			}
		}

		if len(added) > 0 || len(deleted) > 0 {
			logrus.Printf("Calling update: %t", len(added) > 0 || len(deleted) > 0)
			setEnvVarValue(e, strings.Join(msn.Spec.ConsumerNamespaces, ","), c)
			if _, err = dcClient.DeploymentConfigs(ns).Update(d); err != nil {
				logrus.Info("error updating")
				return err
			}
		}
	}

	return nil
}

// TODO: Performance comment
// TODO: Does interface{} work?
// TODO: Can It be genericised more? remove string type from filter?
func arrayDiff(old, new []string) ([]string, []string){
	added := filter(new, func(s string) bool {
		return !contains(old, s)
	})
	deleted := filter(old, func(s string) bool {
		return !contains(new, s)
	})

	return added, deleted
}

// TODO: Does interface{} work?
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
// TODO: Does interface{} work?
// use simple equality
func contains(s []string, b string) bool {
	for _, a := range s {
		if a == b {
			return true
		}
	}
	return false
}


//func updateEnvVar(envVar, newValue string, c *corev1.Container) (bool, *corev1.EnvVar) {
//	for i, e := range c.Env {
//		if e.Name == envVar {
//			if c.Env[i].Value != newValue {
//				c.Env[i].Value = newValue
//				return true, &c.Env[1]
//			}
//		}
//	}
//
//	return false, nil
//}

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

func (icm *integrationControllerManager) createEnmasseConfigMapRoleBinding(namespace string) error {
	rb := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: IntegrationControllerName + "-enmasse-view-",
			Namespace:    EnmasseNamespace,
			Labels: map[string]string{
				"for": IntegrationControllerName,
			},
		},
		RoleRef: clusterRole(EnmasseClusterRoleName),
		Subjects: []rbacv1.Subject{
			serviceAccountSubject(namespace),
		},
	}

	if err := icm.createRoleBinding(EnmasseNamespace, rb); err != nil {
		return err
	}

	return nil
}

func (icm *integrationControllerManager) createRoutesAndServicesRoleBindings(namespaces []string, managedServiceNamespace string) error {
	for _, cns := range namespaces {
		rb := &rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: IntegrationControllerName + "-route-services-",
				Labels: map[string]string{
					"for": "route-services",
				},
			},
			RoleRef: clusterRole(RoutesAndServicesClusterRoleName),
			Subjects: []rbacv1.Subject{
				serviceAccountSubject(managedServiceNamespace),
			},
		}

		if err := icm.createRoleBinding(cns, rb); err != nil {
			return err
		}
	}

	return nil
}
func (icm *integrationControllerManager) deleteRoutesAndServicesRoleBindings(namespaces []string) error {
	logrus.Printf("The namespaces arrray: %v", namespaces)
	for _, ns := range namespaces{
		logrus.Printf("The namespaces arrray has a value: %s", ns)
		err := icm.k8sClient.Rbac().RoleBindings(ns).DeleteCollection(&metav1.DeleteOptions{}, metav1.ListOptions{
			LabelSelector: "for=route-services",
		})
		if err != nil {
			return err
		}
	}

	return nil
}

//func (icm *integrationControllerManager) deleteRoutesAndServicesRoleBindings(namespaces []string) error {
//	for _, ns := range namespaces{
//		roleBindings, err := icm.k8sClient.Rbac().RoleBindings(ns).List(metav1.ListOptions{
//			LabelSelector: "for=route-services",
//		})
//		if err != nil {
//			return err
//		}
//
//		for _, rb := range roleBindings.Items {
//			logrus.Info("The namespace is: " + ns)
//			logrus.Info("The rolebinding name is: " + rb.Name)
//			return icm.k8sClient.Rbac().RoleBindings(ns).Delete(rb.Name, &metav1.DeleteOptions{})
//		}
//
//	}
//
//	return nil
//}

func (icm *integrationControllerManager) createUpdateIntegrationsRoleBinding(msn *integreatly.ManagedServiceNamespace) error {
	rb := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    msn.Name,
			GenerateName: msn.Spec.UserID + "-update-integrations-" + msn.Name + "-",
		},
		RoleRef: clusterRole("integration-update"),
		Subjects: []rbacv1.Subject{
			{
				Kind: "User",
				Name: msn.Spec.UserID,
			},
		},
	}

	if err := icm.createRoleBinding(msn.Name, rb); err != nil {
		return err
	}

	return nil
}

func (icm *integrationControllerManager) createIntegrationController(namespace string) error {
	if _, err := icm.k8sClient.CoreV1().ServiceAccounts(namespace).Create(integrationServiceAccount); err != nil {
		return errors.Wrap(err, "failed to create service account for Integration Controller")
	}

	//if _, err := icm.k8sClient.RbacV1beta1().Roles(namespace).Create(integrationServiceRole); err != nil {
	//	return errors.Wrap(err, "failed to create role for integration controller service")
	//}

	if _, err := icm.k8sClient.RbacV1beta1().RoleBindings(namespace).Create(integrationServiceRoleBinding); err != nil {
		return errors.Wrap(err, "failed to create role binding for integration controller service")
	}

	dcClient, err := icm.osClientFactory.AppsClient()
	if err != nil {
		return errors.Wrap(err, "failed to create an openshift deployment config client")
	}

	if _, err = dcClient.DeploymentConfigs(namespace).Create(integrationDeploymentConfig); err != nil {
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

// TODO: Move to api objects
func clusterRole(roleName string) rbacv1.RoleRef {
	return rbacv1.RoleRef{
		Kind: "ClusterRole",
		Name: roleName,
	}
}

func serviceAccountSubject(namespace string) rbacv1.Subject {
	return rbacv1.Subject{
		Kind:      "ServiceAccount",
		Name:      IntegrationControllerName,
		Namespace: namespace,
	}
}
