package handlers

import (
	"context"
	integreatly "github.com/integr8ly/managed-services-controller/pkg/apis/integreatly/v1alpha1"
	clients "github.com/integr8ly/managed-services-controller/pkg/clients/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
)

func handleManagedServiceNamespace(
	ctx context.Context,
	event sdk.Event,
	msn *integreatly.ManagedServiceNamespace,
	cfg *rest.Config,
) error {

	ns := msn.Name
	msnsc := clients.NewManagedServiceNamespaces(cfg)

	if event.Deleted {
		logrus.Info("Deleting ManagedServiceNamespace: " + ns)
		if err := msnsc.Delete(msn); err != nil {
			return err
		}
	} else {
		if msnsc.Exists(msn) {
			if err := msnsc.Update(msn); err != nil {
				return err
			}
			return nil
		}

		logrus.Info("New ManagedServiceNamespace event")

		// TODO: Improve this validation.
		// When creating a new managed namespace there must be a ConsumerNamespace set.
		if len(msn.Spec.ConsumerNamespaces) == 0 {
			logrus.Info("ManagedServiceNamespace: " + msn.Name + " has no ConsumerNamespace set")
			return nil
		}

		logrus.Info("Creating ManagedServiceNamespace: " + ns)
		if err := msnsc.Create(msn); err != nil {
			return err
		}

		logrus.Info("ManagedServiceNamespace: " + ns + " setup successfully")
	}

	return nil
}
