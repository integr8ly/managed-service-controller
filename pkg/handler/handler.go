package handlers

import (
	"context"
	services "github.com/integr8ly/managed-services-controller/pkg/apis/client/clientset/versioned/typed/services/v1alpha1"
	apis "github.com/integr8ly/managed-services-controller/pkg/apis/integreatly/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/k8sclient"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *apis.ManagedServiceNamespace:
		ns := o.Spec.ManagedNamespace
		c := k8sclient.GetKubeClient()
		if event.Deleted == true {
			err := deleteNamespace(c, ns);if err != nil {
				return err
			}

			logrus.Info("Deleted ManagedServiceNamespace: " + o.Spec.ManagedNamespace)
		} else {
			err := createNamespace(c, ns)
			if err != nil {
				if errors.IsAlreadyExists(err) == true {
					// Namespace already exists. Return silently.
					return nil
				}
				return err
			}
			logrus.Info("New ManagedServiceNamespace event")
			logrus.Info("Created ManagedServiceNamespace: " + o.Spec.ManagedNamespace)

			logrus.Info("Creating fuse operator")
			fuseOperator := services.NewFuseOperator(ns)
			err = fuseOperator.Create();if err != nil {
				return err
			}

			logrus.Info("Creating Integration Controller")
			integrationController := services.NewIntegrationController(c, ns)
			err = integrationController.Create();if err != nil {
				return err
			}

			logrus.Info("ManagedServiceNamespace: " + o.Spec.ManagedNamespace +  " setup successfully")
		}
	}
	return nil
}

func createNamespace(c kubernetes.Interface, namespace string) error{
	_, err := c.Core().Namespaces().Create(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	})
	return err
}

func deleteNamespace(c kubernetes.Interface, namespace string) error{
	return  c.Core().Namespaces().Delete(namespace, &metav1.DeleteOptions{})
}