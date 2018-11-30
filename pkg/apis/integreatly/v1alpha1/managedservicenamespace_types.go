package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ManagedServiceNamespaceSpec defines the desired state of ManagedServiceNamespace
type ManagedServiceNamespaceSpec struct {
	metav1.TypeMeta    `json:",inline"`
	metav1.ObjectMeta  `json:"metadata"`
	ConsumerNamespaces []string `json:"consumerNamespaces"`
	UserID             string   `json:"userId"`
}

// ManagedServiceNamespaceStatus defines the observed state of ManagedServiceNamespace
type ManagedServiceNamespaceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ManagedServiceNamespace is the Schema for the managedservicenamespaces API
// +k8s:openapi-gen=true
type ManagedServiceNamespace struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ManagedServiceNamespaceSpec   `json:"spec,omitempty"`
	Status ManagedServiceNamespaceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ManagedServiceNamespaceList contains a list of ManagedServiceNamespace
type ManagedServiceNamespaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ManagedServiceNamespace `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ManagedServiceNamespace{}, &ManagedServiceNamespaceList{})
}
