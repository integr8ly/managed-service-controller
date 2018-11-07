package v1alpha1

import (
	integreatly "github.com/integr8ly/managed-services-controller/pkg/apis/integreatly/v1alpha1"
)

type fuseOperatorManager struct {}

func NewFuseOperatorManager() ManagedServiceManagerInterface {
	return &fuseOperatorManager{}
}

func (fom *fuseOperatorManager) Create(msn *integreatly.ManagedServiceNamespace) error {
	return nil
}

func (fom *fuseOperatorManager) Update(msn *integreatly.ManagedServiceNamespace) error {
	return nil
}