package v1alpha1

import (
	integreatly "github.com/integr8ly/managed-service-controller/pkg/apis/integreatly/v1alpha1"
	appsv1 "github.com/openshift/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	EnmasseNamespace                   = "enmasse"
	EnmasseClusterRoleName             = "enmasse-integration-viewer"
	RoutesAndServicesClusterRoleName   = "route-service-viewer"
	IntegrationControllerName          = "integration-controller"
	IntegrationUserNamespacesEnvVarKey = "USER_NAMESPACES"
)

var integrationServiceAccount = &corev1.ServiceAccount{
	ObjectMeta: metav1.ObjectMeta{
		Name: "integration-controller",
	},
}

var integrationServiceRoleBinding = &rbacv1beta1.RoleBinding{
	ObjectMeta: metav1.ObjectMeta{
		Name: "ns-integration-controller",
	},
	Subjects: []rbacv1beta1.Subject{
		{
			Kind: "ServiceAccount",
			Name: "integration-controller",
		},
	},
	RoleRef: rbacv1beta1.RoleRef{
		Kind:     "ClusterRole",
		Name:     "integration-controller",
		APIGroup: "rbac.authorization.k8s.io",
	},
}

var integrationDeploymentConfig = &appsv1.DeploymentConfig{
	ObjectMeta: metav1.ObjectMeta{
		Name: IntegrationControllerName,
	},
	Spec: appsv1.DeploymentConfigSpec{
		Strategy: appsv1.DeploymentStrategy{
			Type: "Recreate",
		},
		Replicas: 1,
		Selector: map[string]string{
			"name": IntegrationControllerName,
		},
		Template: &corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"name": IntegrationControllerName,
				},
			},
			Spec: corev1.PodSpec{
				ServiceAccountName: IntegrationControllerName,
				Containers: []corev1.Container{
					{
						Name:  "integration-controller",
						Image: "quay.io/integreatly/" + IntegrationControllerName + ":dev",
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 60000,
								Name:          "metrics",
							},
						},
						Command: []string{
							IntegrationControllerName,
							"--allow-insecure=true",
							"--log-level=debug",
						},
						ImagePullPolicy: "Always",
						Env: []corev1.EnvVar{
							{
								Name: "WATCH_NAMESPACE",
								ValueFrom: &corev1.EnvVarSource{
									FieldRef: &corev1.ObjectFieldSelector{
										FieldPath: "metadata.namespace",
									},
								},
							},
							{
								Name:  "OPERATOR_NAME",
								Value: IntegrationControllerName,
							},
							{
								Name:  "USER_NAMESPACES",
								Value: "",
							},
						},
					},
				},
			},
		},
	},
}

func getEnmasseConfigMapRoleBinding(namespace string) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: IntegrationControllerName + "-enmasse-view-",
			Namespace:    EnmasseNamespace,
			Labels: map[string]string{
				"for": IntegrationControllerName,
			},
		},
		RoleRef: clusterRole(EnmasseClusterRoleName),
		Subjects: []rbacv1.Subject{
			serviceAccountSubject(namespace),
		},
	}
}

func getRoutesAndServicesRoleBinding(managedServiceNamespace string) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: IntegrationControllerName + "-route-services-",
			Labels: map[string]string{
				"for": "route-services",
			},
		},
		RoleRef: clusterRole(RoutesAndServicesClusterRoleName),
		Subjects: []rbacv1.Subject{
			serviceAccountSubject(managedServiceNamespace),
		},
	}
}

func getUpdateIntegrationsRoleBinding(msn *integreatly.ManagedServiceNamespace) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    msn.Name,
			GenerateName: msn.Spec.UserID + "-update-integrations-" + msn.Name + "-",
		},
		RoleRef: clusterRole("integration-update"),
		Subjects: []rbacv1.Subject{
			{
				Kind: "User",
				Name: msn.Spec.UserID,
			},
		},
	}
}

func clusterRole(roleName string) rbacv1.RoleRef {
	return rbacv1.RoleRef{
		Kind: "ClusterRole",
		Name: roleName,
	}
}

func serviceAccountSubject(namespace string) rbacv1.Subject {
	return rbacv1.Subject{
		Kind:      "ServiceAccount",
		Name:      IntegrationControllerName,
		Namespace: namespace,
	}
}
