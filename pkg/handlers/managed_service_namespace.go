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
		if err := msnh.Delete(msn); err != nil {
			return err
		}
	} else {
		if err := msnh.Validate(msn); err != nil {
			logrus.Infof("ManagedServiceNamespace %s is invalid: %s" , msn.Name, err.Error())
			return nil
		}

		if msnh.Exists(msn) {
			if err := msnh.Update(msn); err != nil {
				return err
			}
			return nil
		}

		logrus.Info("New ManagedServiceNamespace event")

		logrus.Info("Creating ManagedServiceNamespace: " + ns)
		if err := msnh.Create(msn); err != nil {
			return err
		}

		logrus.Info("ManagedServiceNamespace: " + ns + " setup successfully")
	}

	return nil
}

func (msnh *ManagedServiceNamespaceHandler) Delete(msn *integreatly.ManagedServiceNamespace) error {
	return msnh.client.Delete(msn)
}

func (msnh *ManagedServiceNamespaceHandler) Validate(msn *integreatly.ManagedServiceNamespace) error {
	return msnh.client.Validate(msn)
}

func (msnh *ManagedServiceNamespaceHandler) Exists(msn *integreatly.ManagedServiceNamespace) bool {
	return msnh.client.Exists(msn)
}

func (msnh *ManagedServiceNamespaceHandler) Update(msn *integreatly.ManagedServiceNamespace) error {
	return msnh.client.Update(msn)
}

func (msnh *ManagedServiceNamespaceHandler) Create(msn *integreatly.ManagedServiceNamespace) error {
	return msnh.client.Create(msn)
}