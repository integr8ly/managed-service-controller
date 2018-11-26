package v1alpha1

import (
	integreatly "github.com/integr8ly/managed-service-controller/pkg/apis/integreatly/v1alpha1"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
)

type fuseOperatorManager struct {
	k8sClient       kubernetes.Interface
	osClientFactory *ClientFactory
	cfg             map[string]string
}

func NewFuseOperatorManager(client kubernetes.Interface, oscf *ClientFactory, cfg map[string]string) ManagedServiceManagerInterface {
	return &fuseOperatorManager{
		k8sClient:       client,
		osClientFactory: oscf,
		cfg:              cfg,
	}
}

func (fom *fuseOperatorManager) Create(msn *integreatly.ManagedServiceNamespace) error {
	if err := fom.createFuseOperator(msn.Name); err != nil {
		return err
	}

	return nil
}

func (fom *fuseOperatorManager) Update(msn *integreatly.ManagedServiceNamespace) error {
	return nil
}

func (fom *fuseOperatorManager) createRoleBindings(namespace string) error {
	if _, err := fom.k8sClient.RbacV1beta1().RoleBindings(namespace).Create(getFuseServiceRoleBinding(fom.cfg["name"])); err != nil {
		return errors.Wrap(err, "failed to create install role binding for fuse service")
	}

	authClient, err := fom.osClientFactory.AuthClient()
	if err != nil {
		return errors.Wrap(err, "failed to create an openshift authorization client")
	}

	if _, err := authClient.RoleBindings(namespace).Create(getViewRoleBinding(fom.cfg["name"])); err != nil {
		return errors.Wrap(err, "failed to create view role binding for fuse service")
	}

	if _, err := authClient.RoleBindings(namespace).Create(getEditRoleBinding(fom.cfg["name"])); err != nil {
		return errors.Wrap(err, "failed to create edit role binding for fuse service")
	}

	return nil
}

func (fom *fuseOperatorManager) createFuseOperator(namespace string) error {
	if _, err := fom.k8sClient.CoreV1().ServiceAccounts(namespace).Create(getFuseServiceAccount(fom.cfg["name"])); err != nil {
		return errors.Wrap(err, "failed to create service account for fuse service")
	}

	if err := fom.createRoleBindings(namespace); err != nil {
		return err
	}

	dcClient, err := fom.osClientFactory.AppsClient()
	if err != nil {
		return errors.Wrap(err, "failed to create an openshift deployment config client")
	}

	if _, err = dcClient.DeploymentConfigs(namespace).Create(getFuseDeploymentConfig(fom.cfg)); err != nil {
		return errors.Wrap(err, "failed to create deployment config for fuse service")
	}

	return nil
}
