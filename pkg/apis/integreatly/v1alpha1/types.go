package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"github.com/pkg/errors"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ManagedServiceNamespaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []ManagedServiceNamespace `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ManagedServiceNamespace struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ManagedServiceNamespaceSpec   `json:"spec"`
}

type ManagedServiceNamespaceSpec struct {
	metav1.TypeMeta                 `json:",inline"`
	metav1.ObjectMeta               `json:"metadata"`
	ConsumerNamespaces []string     `json:"consumerNamespaces"`
	UserID             string       `json:"userId"`
}
// TODO: move to the client
// maybe have a hasNamespace function
func (msn *ManagedServiceNamespace) Validate(clusterNamespaces *corev1.NamespaceList) error {
	for _, v := range msn.Spec.ConsumerNamespaces {
		if !contains(clusterNamespaces.Items, v) {
			return errors.New(msn.Name + " is not valid. The namespace " + v + " does not exist.")
		}
	}

	return nil
}

// TODO: move to utils
func contains(s []corev1.Namespace, name string) bool {
	for _, a := range s {
		if a.Name == name {
			return true
		}
	}
	return false
}
