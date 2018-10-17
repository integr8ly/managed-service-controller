package v1alpha1

import (
	olm "github.com/integr8ly/managed-services-controller/pkg/apis/olm/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	integreatly "github.com/integr8ly/managed-services-controller/pkg/apis/integreatly/v1alpha1"
)

const (
	FuseInstallPlanName           = "syndesis.0.0.1-install"
	FuseClusterServiceVersionName = "syndesis-0.0.1"
)

type fuseOperatorManager struct {}

func NewFuseOperatorManager() ManagedServiceManagerInterface {
	return &fuseOperatorManager{}
}

func (fom *fuseOperatorManager) Create(msn *integreatly.ManagedServiceNamespace) error {
	ns := msn.Name
	ip := &olm.InstallPlan{
		TypeMeta: metav1.TypeMeta{
			APIVersion: olm.SchemeGroupVersion.String(),
			Kind: olm.InstallPlanKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: FuseInstallPlanName,
			Namespace: ns,
		},
		Spec: olm.InstallPlanSpec{
			Approval: olm.ApprovalsAutomatic,
			ClusterServiceVersionNames: []string{
				FuseClusterServiceVersionName,
			},
		},
	}

	ips := NewInstallPlans(ns)
	_, err := ips.Create(ip);if err != nil {
		return err
	}

	return nil
}

func (fom *fuseOperatorManager) Update(msn *integreatly.ManagedServiceNamespace) error {
	return nil
}