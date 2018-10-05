package v1alpha1

import (
	olm "github.com/integr8ly/managed-services-controller/pkg/apis/olm/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
)

type InstallPlanInterface interface {
	Create(*olm.InstallPlan) (*olm.InstallPlan, error)
}

// installPlans implements InstallPlanInterface
type installPlans struct {
	Namespace string
}

// TODO: Why can't this be *InstallPlanInterface returned
func NewInstallPlans(namespace string) InstallPlanInterface {
	return &installPlans{
		Namespace: namespace,
	}
}

func (ips *installPlans) Create(ip *olm.InstallPlan) (*olm.InstallPlan, error) {
	err := sdk.Create(ip);
	return ip, err
}