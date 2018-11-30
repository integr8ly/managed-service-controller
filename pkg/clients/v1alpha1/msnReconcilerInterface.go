package v1alpha1

import (
	integreatlyv1alpha1 "github.com/integr8ly/managed-service-controller/pkg/apis/integreatly/v1alpha1"
)

type MsnReconcilerInterface interface {
	Reconcile(*integreatlyv1alpha1.ManagedServiceNamespace) error
}
