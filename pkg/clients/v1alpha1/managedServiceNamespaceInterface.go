package v1alpha1

import (
	integreatly "github.com/integr8ly/managed-service-controller/pkg/apis/integreatly/v1alpha1"
)

type ManagedServiceNamespaceInterface interface {
	Create(*integreatly.ManagedServiceNamespace) error
	Exists(msn *integreatly.ManagedServiceNamespace) bool
	Delete(msn *integreatly.ManagedServiceNamespace) error
	Update(msn *integreatly.ManagedServiceNamespace) error
	Validate(*integreatly.ManagedServiceNamespace) error
}
