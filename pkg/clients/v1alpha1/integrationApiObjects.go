package v1alpha1

import (
	integreatlyv1alpha1 "github.com/integr8ly/managed-service-controller/pkg/apis/integreatly/v1alpha1"
	appsv1 "github.com/openshift/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

const (
	EnmasseNamespace                 = "enmasse"
	EnmasseClusterRoleName           = "enmasse-integration-viewer"
	RoutesAndServicesClusterRoleName = "route-service-viewer"
)

func getIntegrationServiceAccount(namespace string) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "integration-controller",
			Namespace: namespace,
		},
	}
}

func getIntegrationServiceRoleBinding(namespace string) *rbacv1beta1.RoleBinding {
	return &rbacv1beta1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ns-integration-controller",
			Namespace: namespace,
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
}

func getIntegrationDeploymentConfig(msn *integreatlyv1alpha1.ManagedServiceNamespace, cfg map[string]string) *appsv1.DeploymentConfig {
	return &appsv1.DeploymentConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: cfg["name"],
		},
		Spec: appsv1.DeploymentConfigSpec{
			Strategy: appsv1.DeploymentStrategy{
				Type: "Recreate",
			},
			Replicas: 1,
			Selector: map[string]string{
				"name": cfg["name"],
			},
			Template: &corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": cfg["name"],
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: cfg["name"],
					Containers: []corev1.Container{
						{
							Name:  "integration-controller",
							Image: cfg["imageOrg"] + "/" + cfg["name"] + ":" + cfg["imageTag"],
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 60000,
									Name:          "metrics",
								},
							},
							Command: []string{
								cfg["name"],
								"--allow-insecure=" + cfg["allowInsecure"],
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
									Value: cfg["name"],
								},
								{
									Name:  "USER_NAMESPACES",
									Value: strings.Join(msn.Spec.ConsumerNamespaces, ","),
								},
							},
						},
					},
				},
			},
		},
	}
}

func getEnmasseConfigMapRoleBinding(namespace string, cfg map[string]string) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: cfg["name"] + "-enmasse-view-",
			Namespace:    EnmasseNamespace,
			Labels: map[string]string{
				"for": cfg["name"],
			},
		},
		RoleRef: clusterRole(EnmasseClusterRoleName),
		Subjects: []rbacv1.Subject{
			serviceAccountSubject(namespace, cfg["name"]),
		},
	}
}

func getRoutesAndServicesRoleBinding(consumerNamespace, managedServiceNamespace string, cfg map[string]string) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    consumerNamespace,
			GenerateName: cfg["name"] + "-route-services-",
			Labels: map[string]string{
				"for": "route-services",
			},
		},
		RoleRef: clusterRole(RoutesAndServicesClusterRoleName),
		Subjects: []rbacv1.Subject{
			serviceAccountSubject(managedServiceNamespace, cfg["name"]),
		},
	}
}

func getUpdateIntegrationsRoleBinding(msn *integreatlyv1alpha1.ManagedServiceNamespace) *rbacv1.RoleBinding {
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

func serviceAccountSubject(namespace, name string) rbacv1.Subject {
	return rbacv1.Subject{
		Kind:      "ServiceAccount",
		Name:      name,
		Namespace: namespace,
	}
}
