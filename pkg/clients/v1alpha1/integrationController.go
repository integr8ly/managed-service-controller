package v1alpha1

import (
	"context"
	integreatlyv1alpha1 "github.com/integr8ly/managed-service-controller/pkg/apis/integreatly/v1alpha1"
	appsv1 "github.com/openshift/api/apps/v1"
	osClientAppv1 "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"strings"
)

type integrationControllerManager struct {
	k8sCacheClient  client.Client
	osClientFactory *ClientFactory
	cfg             map[string]string
	mgr             manager.Manager
	dcClient        *osClientAppv1.AppsV1Client
	// k8sCacheClient only seems to have cached the namespace that the operator is deployed to.
	// A client that queries the apiserver directly is needed to get the rolebindings in consumer namespaces
	apiServerClient client.Client
}

func NewIntegrationControllerManager(client client.Client, oscf *ClientFactory, cfg map[string]string, mgr manager.Manager) MsnReconcilerInterface {
	return &integrationControllerManager{
		k8sCacheClient:  client,
		osClientFactory: oscf,
		cfg:             cfg,
		mgr:             mgr,
		dcClient:        nil,
		apiServerClient: nil,
	}
}

func (icm *integrationControllerManager) Reconcile(msn *integreatlyv1alpha1.ManagedServiceNamespace) error {
	ns := msn.Name
	if icm.dcClient == nil {
		dcClient, err := icm.osClientFactory.AppsClient()
		if err != nil {
			return errors.Wrap(err, "failed to create an openshift deployment config client")
		}
		icm.dcClient = dcClient
	}

	if !icm.exists(msn) {
		erb := getEnmasseConfigMapRoleBinding(ns, icm.cfg)
		if err := controllerutil.SetControllerReference(msn, erb, icm.mgr.GetScheme()); err != nil {
			return err
		}
		if err := icm.createRoleBinding(EnmasseNamespace, erb); err != nil {
			return err
		}

		if err := icm.createRoleBinding(ns, getUpdateIntegrationsRoleBinding(msn)); err != nil {
			return err
		}

		if err := icm.createRoutesAndServicesRoleBindings(msn.Spec.ConsumerNamespaces, msn); err != nil {
			return err
		}

		if err := icm.createIntegrationController(msn); err != nil {
			return err
		}

		return nil
	}

	if err := icm.setConsumerNamespaces(msn); err != nil {
		return err
	}

	return nil
}

func (icm *integrationControllerManager) exists(msn *integreatlyv1alpha1.ManagedServiceNamespace) bool {
	_, err := icm.dcClient.DeploymentConfigs(msn.Name).Get(getIntegrationDeploymentConfig(msn, icm.cfg).Name, metav1.GetOptions{})
	if err != nil && apiErrors.IsNotFound(err) {
		return false
	}

	return true
}

func (icm *integrationControllerManager) setConsumerNamespaces(msn *integreatlyv1alpha1.ManagedServiceNamespace) error {
	ns := msn.Name
	d, err := icm.dcClient.DeploymentConfigs(ns).Get(getIntegrationDeploymentConfig(msn, icm.cfg).Name, metav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "Failed to get deployment config")
	}

	c := getContainer(icm.cfg["name"], d)
	if c != nil {
		e := getEnvVar(icm.cfg["namespacesEnvVarKey"], c)
		oldValues := filter(strings.Split(e.Value, ","), func(s string) bool {
			// strings.Split can return an empty string if e.Value is empty.
			return len(s) != 0
		})
		added, deleted := consumerNamespaceDiff(oldValues, msn.Spec.ConsumerNamespaces)
		if len(deleted) > 0 {
			if err := icm.deleteRoutesAndServicesRoleBindings(deleted); err != nil {
				return err
			}
		}

		if len(added) > 0 {
			if err := icm.createRoutesAndServicesRoleBindings(added, msn); err != nil {
				return err
			}
		}

		if len(added) > 0 || len(deleted) > 0 {
			setEnvVarValue(e, strings.Join(msn.Spec.ConsumerNamespaces, ","), c)
			if _, err = icm.dcClient.DeploymentConfigs(ns).Update(d); err != nil {
				return err
			}
		}
	}

	return nil
}

func consumerNamespaceDiff(old, new []string) ([]string, []string) {
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

func (icm *integrationControllerManager) createRoutesAndServicesRoleBindings(namespaces []string, msn *integreatlyv1alpha1.ManagedServiceNamespace) error {
	for _, ns := range namespaces {
		rb := getRoutesAndServicesRoleBinding(ns, msn.Name, icm.cfg)
		if err := controllerutil.SetControllerReference(msn, rb, icm.mgr.GetScheme()); err != nil {
			return err
		}
		if err := icm.createRoleBinding(ns, rb); err != nil {
			return err
		}
	}

	return nil
}

func (icm *integrationControllerManager) deleteRoutesAndServicesRoleBindings(namespaces []string) error {
	if icm.apiServerClient == nil {
		c, err := client.New(icm.mgr.GetConfig(), client.Options{})
		if err != nil {
			return err
		}
		icm.apiServerClient = c
	}

	for _, ns := range namespaces {
		rbL := &rbacv1.RoleBindingList{}
		listOpts := &client.ListOptions{}
		listOpts.SetLabelSelector("for=route-services")
		listOpts.InNamespace(ns)

		if err := icm.apiServerClient.List(context.TODO(), listOpts, rbL); err != nil {
			return err
		}
		for _, rb := range rbL.Items {
			if err := icm.k8sCacheClient.Delete(context.TODO(), &rb); err != nil {
				return err
			}
		}
	}

	return nil
}

func (icm *integrationControllerManager) createIntegrationController(msn *integreatlyv1alpha1.ManagedServiceNamespace) error {
	namespace := msn.Name
	if err := icm.k8sCacheClient.Create(context.TODO(), getIntegrationServiceAccount(namespace)); err != nil {
		return errors.Wrap(err, "failed to create service account for Integration Controller")
	}

	if err := icm.k8sCacheClient.Create(context.TODO(), getIntegrationServiceRoleBinding(namespace)); err != nil {
		return errors.Wrap(err, "failed to create role binding for integration controller service")
	}

	dcClient, err := icm.osClientFactory.AppsClient()
	if err != nil {
		return errors.Wrap(err, "failed to create an openshift deployment config client")
	}

	dc := getIntegrationDeploymentConfig(msn, icm.cfg)
	if _, err = dcClient.DeploymentConfigs(namespace).Create(dc); err != nil {
		return errors.Wrap(err, "failed to create deployment config for integration controller service")
	}

	return nil
}

func (icm *integrationControllerManager) createRoleBinding(namespace string, rb *rbacv1.RoleBinding) error {
	err := icm.k8sCacheClient.Create(context.TODO(), rb)
	if err != nil {
		return err
	}

	return nil
}
