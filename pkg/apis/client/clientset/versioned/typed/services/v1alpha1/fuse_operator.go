package v1alpha1

import (
	olm "github.com/integr8ly/managed-services-controller/pkg/apis/olm/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	resources "github.com/integr8ly/managed-services-controller/pkg/apis/client/clientset/versioned/typed/resources/v1alpha1"
)

const (
	FuseInstallPlanName           = "syndesis.0.0.1-install"
	FuseClusterServiceVersionName = "syndesis-0.0.1"
)

type FuseOperatorInterface interface {
	Create() error
}

// fuseOperator implements FuseOperatorInterface
type fuseOperator struct {
	Namespace string
}

func NewFuseOperator(namespace string) FuseOperatorInterface {
	return &fuseOperator{
		Namespace: namespace,
	}
}

func (f *fuseOperator) Create() error {
	ip := &olm.InstallPlan{
		TypeMeta: metav1.TypeMeta{
			APIVersion: olm.SchemeGroupVersion.String(),// groupName + "/" + version,
			Kind: olm.InstallPlanKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: FuseInstallPlanName,
			Namespace: f.Namespace,
		},
		Spec: olm.InstallPlanSpec{
			Approval: olm.ApprovalsAutomatic,
			ClusterServiceVersionNames: []string{
				FuseClusterServiceVersionName,
			},
		},
	}

	ips := resources.NewInstallPlans(f.Namespace)
	_, err := ips.Create(ip);if err != nil {
		return err
	}

	return nil
}