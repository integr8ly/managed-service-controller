package handlers

import (
	"context"
	integreatly "github.com/integr8ly/managed-service-controller/pkg/apis/integreatly/v1alpha1"
	clients "github.com/integr8ly/managed-service-controller/pkg/clients/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
)

type ManagedServiceNamespaceHandler struct {
	client clients.ManagedServiceNamespaceInterface
}

func NewManagedServiceNamespaceHandler(cfg *rest.Config) *ManagedServiceNamespaceHandler {
	return &ManagedServiceNamespaceHandler{
		client: clients.NewManagedServiceNamespaceClient(cfg),
	}
}

func (msnh *ManagedServiceNamespaceHandler) Handle(
	ctx context.Context,
	event sdk.Event,
	msn *integreatly.ManagedServiceNamespace,
) error {

	ns := msn.Name
	if event.Deleted {
		logrus.Info("Deleting ManagedServiceNamespace: " + ns)
		// TODO: Change these to msnh.Delete(msn)
		if err := msnh.client.Delete(msn); err != nil {
			return err
		}
	} else {
		if err := msnh.client.Validate(msn); err != nil {
			logrus.Info("ManagedServiceNamespace resource is invalid: " + err.Error())
			return nil
		}

		if msnh.client.Exists(msn) {
			if err := msnh.client.Update(msn); err != nil {
				return err
			}
			return nil
		}

		logrus.Info("New ManagedServiceNamespace event")

		logrus.Info("Creating ManagedServiceNamespace: " + ns)
		if err := msnh.client.Create(msn); err != nil {
			return err
		}

		logrus.Info("ManagedServiceNamespace: " + ns + " setup successfully")
	}

	return nil
}
