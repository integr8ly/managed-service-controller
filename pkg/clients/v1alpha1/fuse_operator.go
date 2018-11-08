package v1alpha1

import (
	integreatly "github.com/integr8ly/managed-services-controller/pkg/apis/integreatly/v1alpha1"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
)

const FUSE_IMAGE_STREAMS_NAMESPACE string = "openshift"

type fuseOperatorManager struct {
	k8sClient       kubernetes.Interface
	osClientFactory *ClientFactory
}

func NewFuseOperatorManager(client kubernetes.Interface, oscf *ClientFactory) ManagedServiceManagerInterface {
	return &fuseOperatorManager{
		k8sClient:       client,
		osClientFactory: oscf,
	}
}

func (fom *fuseOperatorManager) Create(msn *integreatly.ManagedServiceNamespace) error {
	namespace := msn.Name

	if _, err := fom.k8sClient.CoreV1().ServiceAccounts(namespace).Create(fuseServiceAccount); err != nil {
		return errors.Wrap(err, "failed to create service account for fuse service")
	}

	if _, err := fom.k8sClient.RbacV1beta1().Roles(namespace).Create(fuseServiceRole); err != nil {
		return errors.Wrap(err, "failed to create role for fuse service")
	}

	if err := fom.createRoleBindings(namespace); err != nil {
		return err
	}

	if err := fom.createFuseOperator(namespace); err != nil {
		return err
	}

	return nil
}

func (fom *fuseOperatorManager) Update(msn *integreatly.ManagedServiceNamespace) error {
	return nil
}

func (fom *fuseOperatorManager) createRoleBindings(namespace string) error {
	if _, err := fom.k8sClient.RbacV1beta1().RoleBindings(namespace).Create(serviceRoleBinding); err != nil {
		return errors.Wrap(err, "failed to create install role binding for fuse service")
	}

	authClient, err := fom.osClientFactory.AuthClient()
	if err != nil {
		return errors.Wrap(err, "failed to create an openshift authorization client")
	}

	if _, err := authClient.RoleBindings(namespace).Create(viewRoleBinding); err != nil {
		return errors.Wrap(err, "failed to create view role binding for fuse service")
	}

	if _, err := authClient.RoleBindings(namespace).Create(editRoleBinding); err != nil {
		return errors.Wrap(err, "failed to create edit role binding for fuse service")
	}

	return nil
}

func (fom *fuseOperatorManager) createFuseOperator(namespace string) error {
	dcClient, err := fom.osClientFactory.AppsClient()
	if err != nil {
		return errors.Wrap(err, "failed to create an openshift deployment config client")
	}

	if _, err = dcClient.DeploymentConfigs(namespace).Create(fuseDeploymentConfig); err != nil {
		return errors.Wrap(err, "failed to create deployment config for fuse service")
	}

	return nil
}
