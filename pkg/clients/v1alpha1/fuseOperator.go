package v1alpha1

import (
	"context"
	integreatlyv1alpha1 "github.com/integr8ly/managed-service-controller/pkg/apis/integreatly/v1alpha1"
	osClientAppv1 "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	"github.com/pkg/errors"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type fuseOperatorManager struct {
	k8sClient       client.Client
	osClientFactory *ClientFactory
	cfg             map[string]string
	dcClient        *osClientAppv1.AppsV1Client
}

func NewFuseOperatorManager(client client.Client, oscf *ClientFactory, cfg map[string]string) MsnReconcilerInterface {
	return &fuseOperatorManager{
		k8sClient:       client,
		osClientFactory: oscf,
		cfg:             cfg,
		dcClient:        nil,
	}
}

func (fom *fuseOperatorManager) Reconcile(msn *integreatlyv1alpha1.ManagedServiceNamespace) error {
	if fom.dcClient == nil {
		dcClient, err := fom.osClientFactory.AppsClient()
		if err != nil {
			return errors.Wrap(err, "failed to create an openshift deployment config client")
		}
		fom.dcClient = dcClient
	}

	if !fom.exists(msn) {
		if err := fom.createFuseOperator(msn.Name); err != nil {
			return err
		}
	}

	return nil
}

func (fom *fuseOperatorManager) exists(msn *integreatlyv1alpha1.ManagedServiceNamespace) bool {
	_, err := fom.dcClient.DeploymentConfigs(msn.Name).Get(getFuseDeploymentConfig(fom.cfg).Name, metav1.GetOptions{})
	if err != nil && apiErrors.IsNotFound(err) {
		return false
	}

	return true
}

func (fom *fuseOperatorManager) createFuseOperator(namespace string) error {
	fsa := getFuseServiceAccount(fom.cfg["name"], namespace)
	if err := fom.k8sClient.Create(context.TODO(), fsa); err != nil {
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

func (fom *fuseOperatorManager) createRoleBindings(namespace string) error {
	fsrb := getFuseServiceRoleBinding(fom.cfg["name"], namespace)
	if err := fom.k8sClient.Create(context.TODO(), fsrb); err != nil {
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
