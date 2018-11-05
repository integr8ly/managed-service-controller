package v1alpha1

import (
	"encoding/json"
	"github.com/golang/glog"
	utils "github.com/integr8ly/managed-services-controller/admission-webhook/pkg/utils/v1alpha1"
	"github.com/integr8ly/managed-services-controller/pkg/apis/integreatly/v1alpha1"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ManagedServiceNamespaceClient struct {
	K8sClient kubernetes.Interface
}

func contains(s []corev1.Namespace, name string) bool {
	for _, a := range s {
		if a.Name == name {
			return true
		}
	}
	return false
}

func (msnc *ManagedServiceNamespaceClient) Validate(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	raw := ar.Request.Object.Raw
	msn := &v1alpha1.ManagedServiceNamespace{}
	err := json.Unmarshal(raw, msn)
	if err != nil {
		glog.Error(err)
		return utils.ToAdmissionResponse(err)
	}

	namespaces, err := msnc.K8sClient.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		glog.Error(err)
		return utils.ToAdmissionResponse(err)
	}

	reviewResponse := v1beta1.AdmissionResponse{
		Result: &metav1.Status{},
	}
	for _, v := range msn.Spec.ConsumerNamespaces {
		if !contains(namespaces.Items, v) {
			reviewResponse.Allowed = false
			reviewResponse.Result.Message = msn.Name + " is not valid. The namespace " + v + " does not exist."
			return &reviewResponse
		}
	}

	glog.V(2).Info("ManagedServiceNamespace " + msn.Name + " valid.")
	reviewResponse.Allowed = true
	return &reviewResponse
}
