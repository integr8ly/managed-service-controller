package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	Status            ManagedServiceNamespaceStatus `json:"status,omitempty"`
}

type ManagedServiceNamespaceSpec struct {
	metav1.TypeMeta                 `json:",inline"`
	metav1.ObjectMeta               `json:"metadata"`
	ManagedNamespace string         `json:"managedNamespace"`
	ConsumerNamespaces []string     `json:"consumerNamespaces"`
}
type ManagedServiceNamespaceStatus struct {
	// Fill me
}
