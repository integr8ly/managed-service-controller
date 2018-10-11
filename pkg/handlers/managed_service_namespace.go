package handlers

import (
	"context"
	clients "github.com/integr8ly/managed-services-controller/pkg/apis/client/clientset/versioned/typed/v1alpha1"
	integreatly "github.com/integr8ly/managed-services-controller/pkg/apis/integreatly/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

func handleManagedServiceNamespace(ctx context.Context, event sdk.Event, msn *integreatly.ManagedServiceNamespace, k8client kubernetes.Interface) error {
	ns := msn.Spec.ManagedNamespace
	msns := clients.NewManagedServiceNamespaces(k8client)

	if event.Deleted == true {
		logrus.Info("Deleting ManagedServiceNamespace: " + ns)
		err := msns.Delete(msn);if err != nil {
			return err
		}
	} else {
		if msns.Exists(msn) == true {
			err := msns.Update(msn);if err != nil {
				return err
			}
			return nil
		}

		logrus.Info("New ManagedServiceNamespace event")
		logrus.Info("Creating ManagedServiceNamespace: " + ns)

		err := msns.Create(msn);if err != nil {
			return err
		}

		logrus.Info("ManagedServiceNamespace: " + ns +  " setup successfully")
	}

	return nil
}